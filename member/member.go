// Package member ...
package member

// Member defines member contents
type Member struct {
	id       int64
	email    string
	ip       string
	nickName string
	birthTs  uint32 /* 生日 */
	country  int    /* 註冊國家 預設1台灣 */
	city     int    /* 註冊城市 */
	gender   int    /* 性別 */
	imei     string /* imei */
	avatar   string /* 大頭貼 */
}

// NewMember returns a new NewMember to read from id, email.
func NewMember(id int64, email string) *Member {
	return &Member{
		id:    id,
		email: email,
	}
}
