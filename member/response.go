// Package member reponse.go定義api回傳結構
package member

// historySuccessResponse 定義api success回傳struct
type historySuccessResponse struct {
	S    int          `json:"s" example:"1"`
	Data []historyLog `json:"data"`
}

// apiFailResponse 定義api fail回傳struct
type apiFailResponse struct {
	S       int    `json:"s" example:"-9"`
	ErrCode string `json:"errCode" example:"APP00_143"`
	ErrMsg  string `json:"errMsg" example:"No such account exists"`
}

// historyLog ...
type historyLog struct {
	MemberID int64  `json:"memberID" example:"1000000001"`
	Device   string `json:"device" example:"A"`
	LoginTs  int    `json:"loginTs" example:"1614650853"`
	IP       string `json:"ip" example:"::1"`
	CreateDt string `json:"createDt" example:"2021-03-02 10:07:34"`
}

// 定義signup api success回傳struct
type signupSuccessResponse struct {
	S    int       `json:"s" example:"1"`
	Data signupMsg `json:"data"`
}

// 註冊成功訊息
type signupMsg struct {
	MemberID int64  `json:"memberID" example:"1000000001"`
	Msg      string `json:"msg" example:"註冊成功"`
}

// 定義login api success回傳struct
type loginSuccessResponse struct {
	S    int      `json:"s" example:"1"`
	Data loginMsg `json:"data"`
}

// 登入發送jwt token
type loginMsg struct {
	MemberID int64  `json:"memberID" example:"1000000001"`
	Msg      string `json:"msg" example:"登入成功"`
	Token    string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImQyd3UxODdAZ21haWwuY29tIiwiUm9sZSI6Im1lbWJlciIsIk1lbWJlcklEIjoxMDAwMDAwMDM0LCJMYW5nIjoiIiwiYXVkIjoiZDJ3dTE4N0BnbWFpbC5jb20iLCJleHAiOjE2MTQ2NjA0MzIsImp0aSI6ImQyd3UxODdAZ21haWwuY29tMTYxNDY1NjgzMiIsImlhdCI6MTYxNDY1NjgzMiwiaXNzIjoibG9naW4iLCJuYmYiOjE2MTQ2NTY4MzMsInN1YiI6ImQyd3UxODdAZ21haWwuY29tIn0.9jUJBPCZgMr4AIWv_JfwVSN9gMVeLbxI8Ck5HrGcknk"`
}
