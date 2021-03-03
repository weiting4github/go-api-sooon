// Package member ...
package member

import (
	"net"
)

// Member defines member contents
type Member struct {
	id       uint64
	email    string
	ip       net.IP
	nickName string
	birthTs  uint32 /* 生日 */
	country  uint16 /* 註冊國家 預設1台灣 */
	city     uint16 /* 註冊城市 */
	gender   uint8  /* 性別 */
	imei     string /* imei */
	avatar   string /* 大頭貼 */
}

// NewMember returns a new NewMember to read from id, email.
func NewMember(id uint64, email string) *Member {
	return &Member{
		id:    id,
		email: email,
	}
}
