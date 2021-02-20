package member

import (
	"go-api-sooon/app"
	"go-api-sooon/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Do 用switch case 對應用戶功能
func Do(c *gin.Context) {
	switch c.Param("action") {
	case "test":
		v, _ := c.Get("memberID")
		c.JSON(http.StatusOK, gin.H{
			"s": 1,
			"c": v,
		})
		return
	case "loginHistory":
		// 使用者登入紀錄
		v, _ := c.Get("memberID")

		// sql
		db, err := config.NewDBConnect()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"s":       -9, // -9系統層級 APP不顯示錯誤訊息
				"errCode": app.DumpErrorCode(loginCodePrefix),
				"errMsg":  err.Error(),
			})
			return
		}
		defer db.Close()

		stmt, err := db.Prepare("SELECT * FROM `sooon_db`.`member_login_log` WHERE `member_id`")

		c.JSON(http.StatusOK, gin.H{
			"s": 1,
			"c": v,
		})
		return
	}
}
