package app

type apiFailResponse struct {
	S       int    `json:"s" example:"-9"`
	ErrCode string `json:"errCode" example:"APP00_143"`
	ErrMsg  string `json:"errMsg" example:"unauthorized"`
}

type initSuccessResponse struct {
	S int `json:"s" example:"1"`
}
