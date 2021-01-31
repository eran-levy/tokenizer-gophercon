package service

import (
	"context"
	"encoding/json"
	"github.com/eran-levy/tokenizer-gophercon/cache"
	"github.com/eran-levy/tokenizer-gophercon/logger"
	"github.com/pkg/errors"
	"strings"
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
}
type tokenizer struct {
	c cache.Cache
}

func New(c cache.Cache) TokenizerService {
	return &tokenizer{c: c}
}

func (t *tokenizer) TokenizeText(ctx context.Context, request TokenizeTextRequest) (TokenizeTextResponse, error) {
	//TODO: process concurrently sentences by newline
	//TODO: call another http api to show an example of retries, etc - the http will predict text lang 200 first chars
	//TODO: telemetry re get/set cache
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
		}
		logger.Log.With("request_id", request.RequestId).With("global_tx_id", request.GlobalTxId).
			Debug("response retrieved from cache")
		return resp, err
	}
	//processed if this global tx id hasnt found in cache
	spt := strings.Split(request.Txt, textSplitSepChar)
	resp = TokenizeTextResponse{RequestId: request.RequestId, TokenizedTxt: spt, NumOfWords: len(spt)}
	//persist in cache
	b, err := json.Marshal(resp)
	if err != nil {
		return TokenizeTextResponse{}, errors.Wrap(err, "could not marshal resp to persist in cache")
	}
	err = t.c.Set(ctx, request.GlobalTxId, b)
	if err != nil {
		logger.Log.With("request_id", request.RequestId).With("global_tx_id", request.GlobalTxId).
			Error("could not persist response in cache")
	}
	return resp, nil
}
