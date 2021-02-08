package main

import (
	"context"
	"github.com/eran-levy/tokenizer-gophercon/pkg/proto/tokenizer"
	"google.golang.org/grpc"
	"log"
	"time"
)

//*** FOR DEMONSTRATION PURPOSES ***
func main() {
	//TODO: retry policy and service health configs - https://github.com/grpc/grpc-go/tree/master/examples/features
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	conn, err := grpc.Dial("localhost:7070", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := tokenizer.NewTokenizerClient(conn)
	req := &tokenizer.TokenizePayloadRequest{GlobalTxId: "111", Text: "MY TEXT", OrganizationId: "1", UserId: "2"}
	res, err := c.GetTokens(ctx, req)
	if err != nil {
		log.Fatalf("could not call gettokens %s", err)
	}
	log.Printf("my resp %+v", res)
}
