package app

type apiFailResponse struct {
	S       int    `json:"s"`
	ErrCode string `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

type initSuccessResponse struct {
	S int `json:"s"`
}
