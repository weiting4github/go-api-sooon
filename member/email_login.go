// Package member 會員
// JWT 會在此驗證發行token
// Sessions 會初始化並塞入用戶資料
package member

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"

	"github.com/wtg42/go-api-sooon/app"

	"github.com/wtg42/go-api-sooon/config"

	"net/http"
	"time"

	"github.com/wtg42/go-api-sooon/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
// @Summary  Login API
// @Tags Login
// @param email formData string true "登入信箱 binding:Email"
// @param p formData string true "登入密碼 binding:RawPWD"
// @param lang formData string false "用戶語系"
// @param client formData string false "用戶裝置"
// @version 1.0
// @produce application/json
// @Failure 200 {object} loginSuccessResponse "登入成功"
// @Failure 400 {object} apiFailResponse
// @host localhost:3000
// @Router /login/email [post]
func Login(c *gin.Context) {
	// form body parameter
	var loginBody LoginBody
	err := c.ShouldBind(&loginBody)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"s":       -9, // -9系統層級 APP不顯示錯誤訊息
			"errCode": app.SFunc.DumpErrorCode(loginCodePrefix),
			"errMsg":  err.Error(),
		})
		return
	}

	// 查email
	ch := make(chan int64) // member_id
	errch := make(chan error)
	tokenString := make(chan string) // jwt
	var email, pwd, salt string
	var memberID int64
	go func() {
		{
			models.DBM.SetQuery("SELECT `member_id`, `email`, `pwd`, `salt` FROM `sooon_db`.`member` WHERE `email` = ?")
			stmt, err := models.DBM.DB.Prepare(models.DBM.GetQuery())
			defer stmt.Close()
			if err != nil {
				errch <- err // 錯誤跳出
				return
			}

			err = stmt.QueryRow(loginBody.Email).Scan(&memberID, &email, &pwd, &salt)
			if err != nil {
				errch <- err // 錯誤跳出
				close(errch)
				return
			}
			// log to stdout
			models.DBM.SQLDebug(loginBody.Email)

			// 驗證密碼
			h := sha256.New()
			h.Write([]byte(loginBody.RawPWD + salt))
			if pwd == fmt.Sprintf("%x", h.Sum(nil)) {
				ch <- memberID // 傳給jwt payload
				close(ch)
			} else {
				fmt.Println(pwd)
				fmt.Printf("%x", h.Sum(nil))
				errch <- sql.ErrNoRows // 密碼錯故意使用sql.ErrNoRows顯示帳密錯誤
				close(errch)
				return
			}

		}
		// 更新使用者登入Log
		{
			models.DBM.SetQuery("INSERT INTO `sooon_db`.`member_login_log`(`member_id`, `client_device`, `login_ts`, `ip`) VALUES (?, ?, ?, ?)")
			stmt, err := models.DBM.DB.Prepare(models.DBM.GetQuery())
			if err != nil {
				// 非致命錯誤 可以寄信通知或是寫入redis做定期排查
				fmt.Println(app.SFunc.DumpErrorCode(loginCodePrefix) + err.Error())
				return
			}
			defer stmt.Close()
			_, err = stmt.Exec(memberID, loginBody.Device, time.Now().Unix(), c.ClientIP())
			if err != nil {
				// 非致命錯誤 可以寄信通知或是寫入redis做定期排查
				fmt.Println(app.SFunc.DumpErrorCode(loginCodePrefix) + err.Error())
				return
			}
			// log to stdout
			models.DBM.SQLDebug(memberID, loginBody.Device, time.Now().Unix(), c.ClientIP())
		}
	}()

	go func() {
		_memberID := <-ch
		// SESSIONS
		{
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
			c.JSON(http.StatusBadRequest, apiFailResponse{
				S:       -9,
				ErrCode: app.SFunc.DumpErrorCode(loginCodePrefix),
				ErrMsg:  err.Error(),
			})
			return
		}
		tokenString <- token
	}()

	// 不設定default讓select強制等待goroutine
	select {
	case _token := <-tokenString: // DB驗證登入
		// 多語
		translation := app.SFunc.Localizer(c, "登入成功")

		c.JSON(http.StatusOK, loginSuccessResponse{
			S: 1,
			Data: loginMsg{
				MemberID: memberID,
				Msg:      translation,
				Token:    _token,
			},
		})
		return
	case err := <-errch:
		s := -1
		switch {
		case err == sql.ErrNoRows:
			s = -1
			err = errors.New(app.SFunc.Localizer(c, "帳號不存在"))
		case err != nil:
			s = -9
		}

		// DB initialization failed
		c.JSON(http.StatusBadRequest, apiFailResponse{
			S:       s,
			ErrCode: app.SFunc.DumpErrorCode(loginCodePrefix),
			ErrMsg:  err.Error(),
		})

		return
	case <-time.After(time.Second * 3):
		c.JSON(http.StatusBadRequest, apiFailResponse{
			S:       -9,
			ErrCode: app.SFunc.DumpErrorCode(loginCodePrefix),
			ErrMsg:  errors.New("Timeout").Error(),
		})
		return
		// default:
		// DB initialization failed
		// c.JSON(http.StatusOK, apiFailResponse{
		// 	S:       -9,
		// 	ErrCode: app.SFunc.DumpErrorCode(loginCodePrefix),
		// 	ErrMsg:  errors.New("DB initialization failed").Error(),
		// })
	}
}
