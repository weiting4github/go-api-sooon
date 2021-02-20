package member

import (
	"database/sql"
	"fmt"
	"go-api-sooon/app"
	"go-api-sooon/config"
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

		stmt, err := db.Prepare("SELECT `member_id` FROM `sooon_db`.`member_login_log` WHERE `member_id` = ?")
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
			var _id int
			if err := rows.Scan(&_id); err != nil {
				c.JSON(http.StatusOK, gin.H{
					"s":       -1, // -9系統層級 APP不顯示錯誤訊息
					"errCode": app.DumpErrorCode(loginCodePrefix),
					"errMsg":  err.Error(),
				})
			}

			r["memberID"] = _id
			log = append(log, r)
		}

		c.JSON(http.StatusOK, gin.H{
			"s":    1,
			"data": log,
		})
		return
	}
}
