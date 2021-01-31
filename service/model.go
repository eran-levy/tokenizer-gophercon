package service

type TokenizeTextRequest struct {
	RequestId string
	Txt       string
}

type TokenizeTextResponse struct {
	RequestId    string
	TokenizedTxt []string
	NumOfWords   int
}
