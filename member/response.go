// Package member reponse.go定義api回傳結構
package member

// historySuccessResponse 定義api success回傳struct
type historySuccessResponse struct {
	S    int          `json:"s"`
	Data []historyLog `json:"data"`
}

// apiFailResponse 定義api fail回傳struct
type apiFailResponse struct {
	S       int    `json:"s"`
	ErrMsg  string `json:"errMsg"`
	ErrCode string `json:"errCode"`
}

// historyLog ...
type historyLog struct {
	MemberID int64  `json:"memberID"`
	Device   string `json:"device"`
	LoginTs  int    `json:"loginTs"`
	IP       string `json:"ip"`
	CreateDt string `json:"createDt"`
}

// 定義signup api success回傳struct
type signupSuccessResponse struct {
	S    int       `json:"s"`
	Data signupMsg `json:"data"`
}

// 註冊成功訊息
type signupMsg struct {
	MemberID int64  `json:"memberID"`
	Msg      string `json:"msg"`
}

// 定義login api success回傳struct
type loginSuccessResponse struct {
	S    int      `json:"s"`
	Data loginMsg `json:"data"`
}

// 登入發送jwt token
type loginMsg struct {
	MemberID int64  `json:"memberID"`
	Msg      string `json:"msg"`
	Token    string `json:"token"`
}
