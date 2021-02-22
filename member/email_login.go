// Package member 會員
// JWT 會在此驗證發行token
// Sessions 會初始化並塞入用戶資料
package member

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"go-api-sooon/app"
	"go-api-sooon/config"

	"go-api-sooon/models"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// loginCodePrefix ...
const loginCodePrefix = "MEM01"

// LoginBody POST參數
type LoginBody struct {
	Email  string `form:"email" binding:"required"`
	RawPWD string `form:"p" binding:"required"`
	Lang   string `form:"lang"`
	Device string `form:"client" binding:"required"`
}

// Login ...
func Login(c *gin.Context) {
	// form body parameter
	var loginBody LoginBody
	err := c.ShouldBind(&loginBody)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"s":       -9, // -9系統層級 APP不顯示錯誤訊息
			"errCode": app.DumpErrorCode(loginCodePrefix),
			"errMsg":  err.Error(),
		})
		return
	}

	// 查email
	ch := make(chan uint64) // member_id
	errch := make(chan error)
	tokenString := make(chan string) // jwt
	var email, pwd, salt string
	var memberID uint64
	go func() {
		stmt, err := models.DBM.DB.Prepare("SELECT `member_id`, `email`, `pwd`, `salt` FROM `sooon_db`.`member` WHERE `email` = ?")
		defer stmt.Close()
		if err != nil {
			errch <- err // 錯誤跳出
			return
		}

		err = stmt.QueryRow(loginBody.Email).Scan(&memberID, &email, &pwd, &salt)
		if err != nil && err != sql.ErrNoRows {
			errch <- err // 錯誤跳出
			close(errch)
			return
		}

		// 驗證密碼
		h := sha256.New()
		h.Write([]byte(loginBody.RawPWD + salt))
		if pwd == fmt.Sprintf("%x", h.Sum(nil)) {
			ch <- memberID // 傳給jwt payload
			close(ch)
		}

		{ // 更新使用者登入Log
			stmt, err := models.DBM.DB.Prepare("INSERT INTO `sooon_db`.`member_login_log`(`member_id`, `client_device`, `login_ts`) VALUES (?, ?, ?)")
			if err != nil {
				// 非致命錯誤 可以寄信通知或是寫入redis做定期排查
				fmt.Println(app.DumpErrorCode(loginCodePrefix) + err.Error())
				return
			}
			defer stmt.Close()
			_, err = stmt.Exec(memberID, loginBody.Device, time.Now().Unix(), c.ClientIP())
			if err != nil {
				// 非致命錯誤 可以寄信通知或是寫入redis做定期排查
				fmt.Println(app.DumpErrorCode(loginCodePrefix) + err.Error())
				return
			}
		}
	}()

	go func() {
		_memberID := <-ch
		{ // 更新sessions
			session := sessions.Default(c)
			lang := c.Request.FormValue("lang")
			if len(lang) <= 0 {
				_membersessions := session.Get(_memberID)
				if _membersessions != nil {
					lang = _membersessions.(config.MemberSessions).Lang
				} else {
					lang = "zh"
				}
			}
			session.Set(_memberID, config.MemberSessions{
				LoginTs: time.Now().Unix(),
				Lang:    lang,
				Email:   email,
			})
			session.Save()
		}

		// 先做JWT 上面用戶編號拿到
		_, token, err := config.CreateJWTClaims(_memberID, loginBody.Email, "member", "login")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"s":       -9,
				"errCode": app.DumpErrorCode(loginCodePrefix),
				"errMsg":  err.Error(),
			})
			return
		}
		tokenString <- token
	}()

	// 不設定default讓select強制等待goroutine
	select {
	case _token := <-tokenString: // DB驗證登入
		// 多語
		localizer := app.Loadi18n(c)
		translation, _ := localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "登入成功",
		})
		c.JSON(http.StatusOK, gin.H{
			"s":      1,
			"member": memberID,
			"token":  _token,
			"msg":    translation,
		})
		return
	case err = <-errch:
		// DB initialization failed
		c.JSON(http.StatusOK, gin.H{
			"s":       -9,
			"errCode": app.DumpErrorCode(loginCodePrefix),
			"errMsg":  err.Error(),
		})
		return
	case <-time.After(time.Second * 3):
		c.JSON(http.StatusOK, gin.H{
			"s":       -9,
			"errCode": app.DumpErrorCode(loginCodePrefix),
			"errMsg":  errors.New("Timeout").Error(),
		})
		return
		// default:
		// 	// DB initialization failed
		// 	c.JSON(http.StatusOK, gin.H{
		// 		"s":       -9,
		// 		"errCode": app.DumpErrorCode(loginCodePrefix),
		// 		"errMsg":  "DB initialization failed",
		// 	})
	}
}
