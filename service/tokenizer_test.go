package service

import (
	"context"
	"github.com/eran-levy/tokenizer-gophercon/cache/mock_cache"
	"github.com/eran-levy/tokenizer-gophercon/logger"
	"github.com/eran-levy/tokenizer-gophercon/repository/mock_repository"
	"github.com/eran-levy/tokenizer-gophercon/telemetry"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTokenizer_TokenizeText(t *testing.T) {
	ai:="tokenizer-gophercon"
	logger.New(logger.Config{LogLevel: "debug", ApplicationId: ai})
	defer logger.Close()
	telem, flush, err := telemetry.New(telemetry.Config{ApplicationID: ai, ServiceName:ai, AgentEndpoint: "http://localhost:14268/api/traces"})
	if err != nil {
		logger.Log.Fatal(err)
	}
	defer flush()
	//for demonstration purposes to test the slow api call
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		rw.Write([]byte(`OK`))
	}))
	// Close the server when test finishes
	defer server.Close()

	tests := []struct {
		name string
		req  TokenizeTextRequest
		resp TokenizeTextResponse
		ctx  context.Context
		e    error
	}{
		{
			name: "fail not provided cache size",
			req:  TokenizeTextRequest{RequestId: "req-1", GlobalTxId: "uuid-1", Txt: "my test"},
			resp: TokenizeTextResponse{RequestId: "req-1", TokenizedTxt: []string{"my", "test"}, NumOfWords: 2},
			ctx:  context.Background(),
			e:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			c := mock_cache.NewMockCache(mockCtrl)
			p := mock_repository.NewMockPersistence(mockCtrl)
			c.EXPECT().Get(tt.ctx,tt.req.GlobalTxId).Return([]byte{},false).MaxTimes(1)
			//b, _ := json.Marshal(tt.resp)
			c.EXPECT().Set(tt.ctx,tt.req.GlobalTxId,gomock.Any()).Return(nil).MaxTimes(1)
			p.EXPECT().StoreMetadata(tt.ctx,gomock.Any()).Return(nil)
			s := New(c,p,telem,server.Client())
			rs,err := s.TokenizeText(tt.ctx,tt.req)
			assert.EqualValues(t, tt.resp,rs)
			assert.Nil(t, err)
		})
	}
}
