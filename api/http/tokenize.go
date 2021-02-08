package http

import (
	"context"
	"github.com/eran-levy/tokenizer-gophercon/api/http/internal"
	"github.com/eran-levy/tokenizer-gophercon/service"
	"github.com/eran-levy/tokenizer-gophercon/telemetry"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"net/http"
	"time"
)

func (s *restApiAdapter) tokenizeTextHandler(c *gin.Context) {
	//can be plugged in a middleware - for demonstration purposes
	const handlerTimeout = 15 * time.Second
	ctx, cancel := context.WithTimeout(c.Request.Context(), handlerTimeout)
	defer cancel()

	ctx, span := s.telemetry.Tracer.Start(ctx, "tokenizeTextHandler")
	defer span.End()
	var r internal.TokenizeTextAPIRequest
	if c.ShouldBindJSON(&r) != nil {
		span.SetStatus(codes.Error, requestNotValid.Error())
		telemetry.IncAPIRequestCounter(c.Request.Context(), 1, telemetry.FailStatusValue)
		c.JSON(http.StatusBadRequest, ErrorResponse{Status: "Invalid request", Message: requestNotValid.Error()})
		return
	}
	dr := service.TokenizeTextRequest{GlobalTxId: r.GlobalTxId, RequestId: uuid.New().String(), Txt: r.Txt}
	tr, err := s.ts.TokenizeText(ctx, dr)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, ErrorResponse{Status: "Tokenize text in request failed", Message: err.Error()})
		return
	}
	span.SetAttributes(telemetry.GlobalTxIdKey.String(r.GlobalTxId), telemetry.ReuqestIdTagKey.String(dr.RequestId))
	telemetry.IncAPIRequestCounter(c.Request.Context(), 1, telemetry.SuccessStatusValue)
	c.JSON(http.StatusOK, internal.TokenizeTextAPIResponse{GlobalTxId: r.GlobalTxId, RequestId: dr.RequestId, TokenizedTxt: tr.TokenizedTxt, NumOfWords: tr.NumOfWords})
}
