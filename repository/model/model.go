package model

import "time"

type TokenizeTextMetadata struct {
	RequestId   string
	GlobalTxId  string
	CreatedDate time.Time
	Language    string
}
