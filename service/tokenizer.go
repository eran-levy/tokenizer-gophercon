package service

import (
	"context"
	"encoding/json"
	"github.com/eran-levy/tokenizer-gophercon/cache"
	"github.com/eran-levy/tokenizer-gophercon/logger"
	"github.com/eran-levy/tokenizer-gophercon/repository"
	"github.com/eran-levy/tokenizer-gophercon/repository/model"
	"github.com/eran-levy/tokenizer-gophercon/telemetry"
	"github.com/pkg/errors"
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
	c cache.Cache
	p repository.Persistence
	t telemetry.Telemetry
}

func New(c cache.Cache, p repository.Persistence, t telemetry.Telemetry) TokenizerService {
	return &tokenizer{c: c, p: p, t: t}
}

func (t *tokenizer) TokenizeText(ctx context.Context, request TokenizeTextRequest) (TokenizeTextResponse, error) {
	//TODO: process concurrently sentences by newline
	//TODO: call another http api to show an example of retries, etc - the http will predict text lang 200 first chars
	//TODO: context cancelation
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
