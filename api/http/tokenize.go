package http

import (
	"github.com/eran-levy/tokenizer-gophercon/api/http/internal"
	"github.com/eran-levy/tokenizer-gophercon/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (s *RestApiAdapter) tokenizeTextHandler(c *gin.Context) {
	var r internal.TokenizeTextAPIRequest
	if c.ShouldBindJSON(&r) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Status: "Invalid request", Message: requestNotValid.Error()})
		return
	}

	dr := service.TokenizeTextRequest{RequestId: uuid.New().String(), Txt: r.Txt}
	tr, err := s.ts.TokenizeText(c.Request.Context(), dr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Status: "Tokenize text in request failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, internal.TokenizeTextAPIResponse{GlobalTxId: r.GlobalTxId, RequestId: dr.RequestId, TokenizedTxt: tr.TokenizedTxt, NumOfWords: tr.NumOfWords})

}
