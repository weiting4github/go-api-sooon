/*main
  cyclesooon
  主程式進入點
*/
package main

import (
	"context"
	"fmt"
	"go-api-sooon/app"
	"go-api-sooon/approuter"
	"go-api-sooon/config"
	"go-api-sooon/member"
	m "go-api-sooon/myplayground"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	_ "github.com/joho/godotenv/autoload" // 環境變數套件os.Getenv
)

// mainCodePrefix 錯誤代碼追蹤
const mainCodePrefix = "MA00"

func main() {
	// PLAYGROUND TEST YOUR CODE
	approuter.GinRouterGroup.GET("/playground", playground)
	// PLAYGROUND TEST YOUR CODE

	// MARK: dev
	// APP 打開APP先打這支認證 過了APP才能去要JWT token
	approuter.GinRouterGroup.POST("/init", func(c *gin.Context) {
		s := c.PostForm("hash")
		if s != app.APIAuthorized() {
			c.JSON(http.StatusUnauthorized, gin.H{
				"s":       -9,
				"errCode": app.DumpErrorCode(mainCodePrefix),
				"errMsg":  "unauthorized",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"s": 1,
		})
	})

	// 用戶註冊
	approuter.GinRouterGroup.POST("/signup", member.NewMemberReg)

	// email登入
	approuter.GinRouterGroup.POST("/login/email", member.Login)
	approuter.GinRouterGroup.Use(config.JWTAuth)
	{
		approuter.GinRouterGroup.GET("/member/:action/:mid", member.Do)
	}

	// approuter.GinEngine.Run(":3000")

	// 原本是用router.Run()，要使用net/http套件的shutdown的話，需要使用原生的ListenAndServe
	srv := &http.Server{
		Addr:    ":3000",
		Handler: approuter.GinEngine,
	}
	//新增一個channel，type是os.Signal
	ch := make(chan os.Signal, 1)
	//call goroutine啟動http server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("SERVER GG惹:", err)
		}
	}()
	//Notify：將系統訊號轉發至channel
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	//阻塞channel
	<-ch

	//收到關機訊號時做底下的流程
	fmt.Println("Graceful Shutdown start signal")
	//透過context.WithTimeout產生一個新的子context，它的特性是有生命週期，這邊是設定10秒
	//只要超過10秒就會自動發出Done()的訊息
	ctx := context.Background()
	c, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	fmt.Println("Graceful Shutdown start context.Background()")
	//使用net/http的shutdown進行關閉http server，參數是上面產生的子context，會有生命週期10秒，
	//所以10秒內要把request全都消化掉，如果超時一樣會強制關閉，所以如果http server要處理的是
	//需要花n秒才能處理的request就要把timeout時間拉長一點
	if err := srv.Shutdown(c); err != nil {
		log.Println("srv.Shutdown:", err)
	}
	//使用select去阻塞主線程，當子context發出Done()的訊號才繼續向下走
	select {
	case <-c.Done():
		fmt.Println("Graceful Shutdown start c.Done()")
		close(ch)
	}
}

// PLAYGROUND TEST YOUR CODE
func playground(c *gin.Context) {
	// ip := net.ParseIP("2001:db8::68")
	// bint := ipv6ToInt(ip)
	// fmt.Println(net.IPv4(byte("192"), byte("192"), byte("192"), byte("192")))

	// 陣列值2個數字相加等於4 且是唯一解
	r := m.TwoSum([]int{3, 4, 1, 2}, 4)
	fmt.Println(r)

	// 陣列反轉
	// myplay.ArrReverse([]int{6, 4, 3, 1})

	// 費式數列
	// for i := 0; i < 20; i++ {
	// 	app.DumpAnyLikeABoss(myplay.Fibonacci1()(i))
	// }

	// store := cookie.NewStore([]byte(os.Getenv("SESSION_KEY")))
	// approuter.GinRouterGroup.Use(sessions.Sessions("testSessions", store))
	// session := sessions.Default(c)
	// session.Set("hello", "world")
	// session.Set("mycookie", "yes done!")
	// // session.Clear()
	// session.Save()

	c.JSON(http.StatusOK, gin.H{
		"s": 1,
	})
	return
}

func ipv6ToInt(IPv6Addr net.IP) *big.Int {
	IPv6Int := big.NewInt(0)
	IPv6Int.SetBytes(IPv6Addr)
	return IPv6Int
}
