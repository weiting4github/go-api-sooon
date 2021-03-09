package member

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/wtg42/go-api-sooon/app"
	"github.com/wtg42/go-api-sooon/models"

	"github.com/gin-gonic/gin"
)

const detailCodePrefix = "MEM02"

// Do loginHistory: 登入紀錄,
// @Summary  Member detail API
// @Tags Member
// @param action path string true "profile|loginHistory"
// @param mid path string true "用戶編號(1000000001)"
// @param Authorization header string false "JWT"
// @version 1.0
// @produce json
// @Success 200 {object} historySuccessResponse "登入紀錄"
// @Success 201 {object} historySuccessResponse "個人檔案"
// @Failure 400 {object} apiFailResponse
// @host localhost:3000
// @Router /member/{action}/{mid} [get]
func Do(c *gin.Context) {
	switch c.Param("action") {
	// 個人檔案
	case "profile":
		memberID := c.Param("mid")
		v, _ := c.Get("memberID")
		c.JSON(http.StatusOK, gin.H{
			"s":         1,
			"c":         v,
			"profileID": memberID,
		})
		return
	case "loginHistory":
		// 使用者登入紀錄
		memberID := c.Param("mid")

		// 自己才能看自己紀錄
		chkMemberID, _ := strconv.ParseInt(memberID, 10, 64)
		if memberIDSelf, ok := c.Get("memberID"); ok == false || memberIDSelf.(int64) != chkMemberID {
			c.JSON(http.StatusBadRequest, apiFailResponse{
				S:       -9,
				ErrCode: app.SFunc.DumpErrorCode(detailCodePrefix),
				ErrMsg:  errors.New("Session Invalid").Error(),
			})
			return
		}
		models.DBM.SetQuery("SELECT * FROM `sooon_db`.`member_login_log` WHERE `member_id` = ?")
		stmt, err := models.DBM.DB.Prepare(models.DBM.GetQuery())
		defer stmt.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, apiFailResponse{
				S:       -9,
				ErrCode: app.SFunc.DumpErrorCode(detailCodePrefix),
				ErrMsg:  err.Error(),
			})
			return
		}

		rows, err := stmt.Query(memberID)
		defer rows.Close()
		if err != nil && err != sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, apiFailResponse{
				S:       -9,
				ErrCode: app.SFunc.DumpErrorCode(detailCodePrefix),
				ErrMsg:  err.Error(),
			})
			return
		}
		// log to stdout
		models.DBM.SQLDebug(memberID)

		var logs = []historyLog{}
		for rows.Next() {
			var log historyLog
			if err := rows.Scan(&log.MemberID, &log.Device, &log.LoginTs, &log.IP, &log.CreateDt); err != nil {
				c.JSON(http.StatusBadRequest, apiFailResponse{
					S:       -9,
					ErrCode: app.SFunc.DumpErrorCode(detailCodePrefix),
					ErrMsg:  err.Error(),
				})
				return
			}

			logs = append(logs, log)
		}

		c.JSON(http.StatusOK, historySuccessResponse{
			S:    1,
			Data: logs,
		})
		return
	default:
		c.JSON(http.StatusOK, historySuccessResponse{
			S:    1,
			Data: []historyLog{},
		})
		return
	}

}
