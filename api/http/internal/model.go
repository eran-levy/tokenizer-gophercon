package internal

type TokenizeTextAPIRequest struct {
	GlobalTxId string `json:"global_tx_id" binding:"required"`
	Txt        string `json:"text" binding:"required"`
}

type TokenizeTextAPIResponse struct {
	GlobalTxId   string   `json:"global_tx_id"`
	RequestId    string   `json:"request_id"`
	TokenizedTxt []string `json:"tokenized_text,omitempty"`
	NumOfWords   int      `json:"num_of_words"`
}
