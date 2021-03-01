// Error code prefix
// 用戶註冊

package member

import (
	"crypto/sha256"
	"fmt"
	"go-api-sooon/app"
	"go-api-sooon/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// signupCodePrefix 錯誤代碼追蹤
const signupCodePrefix = "MEM00"

// RegInfo POST參數
type RegInfo struct {
	RegEmail string `form:"email" binding:"required"`
	Pwd      string `form:"p" binding:"required"` /* 密碼 */
	// RegNickName string `form:"nickName" binding:"required"`
	// RegBirthTs  uint32 `form:"birthday" binding:"required"`    /* 生日 */
	// CPW         string `form:"cpw" binding:"required"`         /* 再次確認密碼 */
	// Country     int    `form:"country" binding:"required"`     /* 註冊國家 預設1台灣 */
	// City        int    `form:"city" binding:"required"`        /* 註冊城市 */
	// Gender      int    `form:"gender" binding:"required"`      /* 性別 */
	// ExpectType  int    `form:"expectType" binding:"required"`  /* 尋找的旅遊類型 */
	// RegImei     string `form:"imei" binding:"required"`        /* imei */
	// ProfileShot string `form:"profileShot" binding:"required"` /* 大頭貼 */
}

// NewMemberReg POST參數
// @Summary  Signup API
// @Tags Sign Up
// @param email formData string true "註冊信箱 binding:RegEmail"
// @param p formData string true "註冊密碼 binding:Pwd"
// @version 1.0
// @produce application/json
// @Success 200 {object} signupSuccessResponse "註冊成功"
// @Failure 400 {object} apiFailResponse
// @host localhost:3000
// @Router /signup [post]
func NewMemberReg(c *gin.Context) {
	var reginfo RegInfo
	err := c.ShouldBind(&reginfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, apiFailResponse{
			S:       -9,
			ErrCode: app.SFunc.DumpErrorCode(detailCodePrefix),
			ErrMsg:  err.Error(),
		})
		return
	}

	// db salt column
	salt := app.SFunc.NewMd5String(3)
	// sha256雜湊
	h := sha256.New()
	h.Write([]byte(reginfo.Pwd + salt))
	hashPWD := fmt.Sprintf("%x", h.Sum(nil))

	models.DBM.SetQuery("INSERT IGNORE INTO `sooon_db`.`member`(`email`, `pwd`, `salt`, `ip_field`, `ipv4v6`, `create_ts`) VALUES(?, ?, ?, ?, INET6_ATON(?), ?)")
	stmtIns, err := models.DBM.DB.Prepare(models.DBM.GetQuery())
	defer stmtIns.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, apiFailResponse{
			S:       -9,
			ErrCode: app.SFunc.DumpErrorCode(detailCodePrefix),
			ErrMsg:  err.Error(),
		})
		return
	}

	result, err := stmtIns.Exec(reginfo.RegEmail, hashPWD, salt, c.ClientIP(), c.ClientIP(), time.Now().Unix())
	if err != nil {
		c.JSON(http.StatusBadRequest, apiFailResponse{
			S:       -9,
			ErrCode: app.SFunc.DumpErrorCode(detailCodePrefix),
			ErrMsg:  err.Error(),
		})
		return
	}
	// log to stdout
	models.DBM.SQLDebug(reginfo.RegEmail, hashPWD, salt, c.ClientIP(), c.ClientIP(), time.Now().Unix())

	newMember, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusBadRequest, apiFailResponse{
			S:       -9,
			ErrCode: app.SFunc.DumpErrorCode(detailCodePrefix),
			ErrMsg:  err.Error(),
		})
		return
	}

	var outputMsg string
	var s int
	if newMember == 0 {
		s = -1
		outputMsg = "註冊失敗"
	} else {
		s = 1
		outputMsg = "註冊成功"
	}
	// 語系
	translation := app.SFunc.Localizer(c, outputMsg)

	data := signupMsg{
		MemberID: newMember,
		Msg:      translation,
	}
	c.JSON(http.StatusOK, signupSuccessResponse{
		S:    s,
		Data: data,
	})
}
