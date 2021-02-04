package grpc

import (
	"context"
	"github.com/eran-levy/tokenizer-gophercon/logger"
	"github.com/eran-levy/tokenizer-gophercon/pkg/proto/tokenizer"
	"github.com/eran-levy/tokenizer-gophercon/service"
	"github.com/eran-levy/tokenizer-gophercon/telemetry"
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
)

type grpcApiAdapter struct {
	cfg Config
	srv *grpc.Server
	ts  service.TokenizerService
	t   telemetry.Telemetry
	//https://github.com/grpc/grpc-go/issues/3669
	tokenizer.UnimplementedTokenizerServer
}

type Config struct {
	GrpcAddress           string
	MaxConnectionAge      time.Duration
	MaxConnectionAgeGrace time.Duration
}

func New(cfg Config, ts service.TokenizerService, t telemetry.Telemetry) *grpcApiAdapter {
	return &grpcApiAdapter{cfg: cfg, ts: ts, t: t}
}

func (g *grpcApiAdapter) GetTokens(ctx context.Context, r *tokenizer.TokenizePayloadRequest) (*tokenizer.TokenizePayloadReresponse, error) {
	logger.Log.With("global_tx_id", r.GlobalTxId).With("organization_id", r.OrganizationId).Info("processing grpc request")
	ctx, span := g.t.Tracer.Start(ctx, "get tokens grpc")
	defer span.End()
	dr := service.TokenizeTextRequest{GlobalTxId: r.GlobalTxId, RequestId: uuid.New().String(), Txt: r.Text}
	tr, err := g.ts.TokenizeText(ctx, dr)
	if err != nil {
		span.SetStatus(codes.Error, "getTokens error grpc")
		return &tokenizer.TokenizePayloadReresponse{}, err
	}
	return &tokenizer.TokenizePayloadReresponse{GlobalTxId: r.GlobalTxId, TokenizedText: tr.TokenizedTxt, Language: "ENGLISH"}, nil
}

func (g *grpcApiAdapter) Start(fatalErrors chan<- error) {
	const network = "tcp"
	l, err := net.Listen(network, g.cfg.GrpcAddress)
	if err != nil {
		fatalErrors <- err
	}
	g.srv = grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionAge: g.cfg.MaxConnectionAge,
		MaxConnectionAgeGrace: g.cfg.MaxConnectionAgeGrace}), grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
	tokenizer.RegisterTokenizerServer(g.srv, g)
	err = g.srv.Serve(l)
	if err != nil {
		logger.Log.Errorf("grpc server closed: %s", err)
		fatalErrors <- err
	}
}

func (g *grpcApiAdapter) Close() {
	//this call pending server close, you can set timeout by creating some select timer to define time for graceful shutdown
	g.srv.GracefulStop()
}
