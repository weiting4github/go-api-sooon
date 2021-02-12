package member

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Do ...
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
		c.JSON(http.StatusOK, gin.H{
			"s": 1,
			"c": v,
		})
		return
	}
}
