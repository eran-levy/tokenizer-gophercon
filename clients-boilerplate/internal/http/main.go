package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	const numOfRetries = 3
	//use timeout
	//DefaultTransport caching and re-using connections see http.Client{} for more info
	client := &http.Client{Timeout: time.Second * 30}
	reqBody := strings.NewReader(`
{
    "global_tx_id":"22-11",
    "text":"mytest hello"
}
`)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	b, c, err := doReq(ctx, client, reqBody)
	if err == nil {
		log.Printf("%+v", string(b))
		if c == http.StatusOK {
			log.Printf("%+v", string(b))
			return
		}
		if !isRetryable(c) {
			log.Print("couldnt retry")
			return
		}
	}

	//retryNum := 0
	//for retryNum < numOfRetries {
	//	select {
	//	case <-ctx.Done():
	//		log.Fatalf("context done: %s", ctx.Err())
	//	default:
	//
	//	}
	//	ctx, cancel := context.WithTimeout(ctx, time.Second)
	//	b, c, err := doReq(ctx, client, reqBody)
	//	cancel()
	//	if err == nil {
	//		if c == http.StatusOK {
	//			log.Printf("%+v", string(b))
	//			return
	//		}
	//		if !isRetryable(c) {
	//			log.Print("couldnt retry")
	//			return
	//		}
	//	}
	//	time.Sleep(1 * time.Second)
	//	retryNum++
	//}
}

func doReq(ctx context.Context, client *http.Client, r *strings.Reader) ([]byte, int, error) {
	//go func() {
	//	select {
	//	case <- ctx.Done():
	//		log.Printf("routine context info: %s", ctx.Err())
	//	}
	//}()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8080/v1/tokenize", r)
	if err != nil {
		return []byte{}, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []byte{}, resp.StatusCode, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, resp.StatusCode, err
	}
	return body, resp.StatusCode, err
}

func isRetryable(code int) bool {
	if code <= 399 {
		return true
	}
	return false
}
