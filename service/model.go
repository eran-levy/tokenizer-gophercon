package service

type TokenizeTextRequest struct {
	RequestId  string `json:"request_id"`
	GlobalTxId string `json:"global_tx_id"`
	Txt        string `json:"txt"`
}

type TokenizeTextResponse struct {
	RequestId    string   `json:"request_id"`
	TokenizedTxt []string `json:"tokenized_text,omitempty"`
	NumOfWords   int      `json:"num_of_words"`
}
