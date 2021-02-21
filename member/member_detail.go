package member

import (
	"database/sql"
	"fmt"
	"go-api-sooon/app"
	"go-api-sooon/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Do 用switch case 對應用戶功能
func Do(c *gin.Context) {
	fmt.Println(c.Param("action"))
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
		memberID, _ := c.Get("memberID")

		// stmt, err := models.DBM.DB.Prepare("SELECT * FROM `sooon_db`.`member_login_log` WHERE `member_id` = ?")
		stmt, err := models.DBM.SelMemberLoginLog()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"s":       -9, // -9系統層級 APP不顯示錯誤訊息
				"errCode": app.DumpErrorCode(loginCodePrefix),
				"errMsg":  err.Error(),
			})
			return
		}
		defer stmt.Close()
		rows, err := stmt.Query(memberID)
		if err != nil && err != sql.ErrNoRows {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"s":       -9, // -9系統層級 APP不顯示錯誤訊息
				"errCode": app.DumpErrorCode(loginCodePrefix),
				"errMsg":  err.Error(),
			})
			return
		}
		defer rows.Close()

		var log []map[string]interface{}
		for rows.Next() {
			r := make(map[string]interface{})
			var _id, _loginTS int
			var _device, _createDt string
			if err := rows.Scan(&_id, &_device, &_loginTS, &_createDt); err != nil {
				c.JSON(http.StatusOK, gin.H{
					"s":       -1, // -9系統層級 APP不顯示錯誤訊息
					"errCode": app.DumpErrorCode(loginCodePrefix),
					"errMsg":  err.Error(),
				})
			}

			r["memberID"] = _id
			r["device"] = _device
			r["loginTs"] = _loginTS
			r["createDt"] = _createDt
			log = append(log, r)
		}

		c.JSON(http.StatusOK, gin.H{
			"s":    1,
			"data": log,
		})
		return
	}
}
