package service

import (
	"context"
	"errors"
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
}

func New() TokenizerService {
	return &tokenizer{}
}

func (t *tokenizer) TokenizeText(ctx context.Context, request TokenizeTextRequest) (TokenizeTextResponse, error) {
	if len(request.Txt) > txtSizeLimitInBytes {
		return TokenizeTextResponse{}, TextSizeExceedsMaxLimitBytesError
	}
	spt := strings.Split(request.Txt, textSplitSepChar)
	return TokenizeTextResponse{RequestId: request.RequestId, TokenizedTxt: spt, NumOfWords: len(spt)}, nil
}
