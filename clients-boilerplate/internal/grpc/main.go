package main

import (
	"context"
	"github.com/eran-levy/tokenizer-gophercon/pkg/proto/tokenizer"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := tokenizer.NewTokenizerClient(conn)
	req := &tokenizer.TokenizePayloadRequest{GlobalTxId: "111", Text: "MY TEXT", OrganizationId: "1", UserId: "2"}
	res, err := c.GetTokens(context.Background(), req)
	if err != nil {
		log.Fatalf("could not call gettokens %s", err)
	}
	log.Printf("my resp %+v", res)
}
