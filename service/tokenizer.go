package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eran-levy/tokenizer-gophercon/cache"
	"github.com/eran-levy/tokenizer-gophercon/logger"
	"github.com/eran-levy/tokenizer-gophercon/repository"
	"github.com/eran-levy/tokenizer-gophercon/repository/model"
	"github.com/eran-levy/tokenizer-gophercon/telemetry"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	txtSizeLimitInBytes = 1000
	textSplitSepChar    = " "
)

var (
	TextSizeExceedsMaxLimitBytesError = errors.New("text size it greater than the limit in Bytes")
)

type TokenizerService interface {
	TokenizeText(ctx context.Context, request TokenizeTextRequest) (TokenizeTextResponse, error)
	IsServiceHealthy(ctx context.Context) (bool, error)
}
type tokenizer struct {
	c        cache.Cache
	p        repository.Persistence
	t        telemetry.Telemetry
	htClient *http.Client
}

func New(c cache.Cache, p repository.Persistence, t telemetry.Telemetry, htClient *http.Client) TokenizerService {
	return &tokenizer{c: c, p: p, t: t, htClient: htClient}
}

func (t *tokenizer) TokenizeText(ctx context.Context, request TokenizeTextRequest) (TokenizeTextResponse, error) {
	//TODO: process concurrently sentences by newline
	if len(request.Txt) > txtSizeLimitInBytes {
		return TokenizeTextResponse{}, TextSizeExceedsMaxLimitBytesError
	}
	var resp TokenizeTextResponse
	//try to pick up from cache
	v, found := t.c.Get(ctx, request.GlobalTxId)
	if found {
		err := json.Unmarshal(v, &resp)
		if err != nil {
			logger.Log.With("request_id", request.RequestId).With("global_tx_id", request.GlobalTxId).
				Error("response found in cache but could not unmarshal cached request, reprocessing request")
			telemetry.IncTokenizeRequestCounter(ctx, 1, found, telemetry.FailStatusValue)
		}
		logger.Log.With("request_id", request.RequestId).With("global_tx_id", request.GlobalTxId).
			Debug("response retrieved from cache")
		telemetry.IncTokenizeRequestCounter(ctx, 1, found, telemetry.SuccessStatusValue)
		return resp, err
	}
	//dummy slow http request to demonstrate
	err := doSlowCallWithRetry(ctx, t.htClient)
	if err != nil {
		return TokenizeTextResponse{}, errors.Wrap(err, "called slow http and failed")
	}

	//processed if this global tx id hasnt found in cache
	spt := strings.Split(request.Txt, textSplitSepChar)
	resp = TokenizeTextResponse{RequestId: request.RequestId, TokenizedTxt: spt, NumOfWords: len(spt)}
	//persist in cache
	b, err := json.Marshal(resp)
	if err != nil {
		telemetry.IncTokenizeRequestCounter(ctx, 1, found, telemetry.FailStatusValue)
		return TokenizeTextResponse{}, errors.Wrap(err, "could not marshal resp to persist in cache")
	}
	err = t.c.Set(ctx, request.GlobalTxId, b)
	if err != nil {
		telemetry.IncTokenizeRequestCounter(ctx, 1, found, telemetry.FailStatusValue)
		logger.Log.With("request_id", request.RequestId).With("global_tx_id", request.GlobalTxId).With("error", err).
			Error("could not persist response in cache")
	}
	//persist in metastore
	err = t.p.StoreMetadata(ctx, model.TokenizeTextMetadata{RequestId: request.RequestId, GlobalTxId: request.GlobalTxId, CreatedDate: time.Now().UTC(), Language: "English"})
	if err != nil {
		telemetry.IncTokenizeRequestCounter(ctx, 1, found, telemetry.FailStatusValue)
		logger.Log.With("request_id", request.RequestId).With("global_tx_id", request.GlobalTxId).With("error", err).
			Error("could not persist metadata in db")
	}
	telemetry.IncTokenizeRequestCounter(ctx, 1, found, telemetry.SuccessStatusValue)
	return resp, nil
}

func doSlowCallWithRetry(ctx context.Context, client *http.Client) error {
	const (
		numOfRetries        = 3
		callHttpCtxTimeout  = 5 * time.Second
		waitBetweenRequests = time.Second
	)

	retryNum := 0
	for retryNum < numOfRetries {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		ctx, cancel := context.WithTimeout(ctx, callHttpCtxTimeout)
		c, err := doReq(ctx, client)
		cancel()
		if err == nil {
			if c == http.StatusOK {
				return nil
			}
			if !isRetryable(c) {
				return fmt.Errorf("couldnt retry http code %d", c)
			}
		}
		//backoff between retries
		time.Sleep(waitBetweenRequests)
		retryNum++
	}
	return fmt.Errorf("could not call http service after retries")
}

func doReq(ctx context.Context, client *http.Client) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:3333/language", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, err
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}
	return resp.StatusCode, err
}
func (t *tokenizer) IsServiceHealthy(ctx context.Context) (bool, error) {
	h, err := t.p.IsServiceHealthy(ctx)
	if !h {
		return h, err
	}
	h, err = t.c.IsServiceHealthy(ctx)
	if !h {
		return h, err
	}
	return h, nil
}

func isRetryable(code int) bool {
	if code <= 399 {
		return true
	}
	return false
}
