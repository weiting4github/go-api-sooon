package approuter

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload" // 環境變數套件os.Getenv
)

// GinEngine default engine
var GinEngine *gin.Engine

// GinRouterGroup ...
var GinRouterGroup *gin.RouterGroup

func init() {
	gin.SetMode(os.Getenv("GIN_MODE")) // dev release 切換
	GinEngine, GinRouterGroup = switchRouter(gin.Mode())
}

func switchRouter(ginMode string) (*gin.Engine, *gin.RouterGroup) {
	var engine *gin.Engine
	var group *gin.RouterGroup
	switch ginMode {
	case "debug", "test":
		//RouterD DEV
		var engine = gin.New()
		engine.Use(gin.Logger())
		engine.Use(gin.Recovery())
		// 設定Redis Sessions
		store, err := redis.NewStore(10, "tcp", "localhost:30001", "", []byte(os.Getenv("SESSION_KEY"))) // idle connections 10, Close connections after remaining idle for this duration(240sec).
		if err != nil {
			panic(err.Error()) // 嚴重錯誤
		}
		engine.Use(sessions.Sessions("dev_sessions", store))

		{ // log
			f, _ := os.Create("gin.log")
			gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
			// defer f.Close()
			engine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
				// your custom format
				return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
					param.ClientIP,
					param.TimeStamp.Format(time.RFC1123),
					param.Method,
					param.Path,
					param.Request.Proto,
					param.StatusCode,
					param.Latency,
					param.Request.UserAgent(),
					param.ErrorMessage,
				)
			}))
		}
		//Dev ...
		var group = engine.Group("/dev")
		return engine, group
	case "release":
		// RouterR RELEASE
		var engine = gin.New()
		engine.Use(gin.Logger())
		engine.Use(gin.Recovery())
		// 設定Redis Sessions
		store, err := redis.NewStore(10, "tcp", "localhost:30001", "", []byte(os.Getenv("SESSION_KEY"))) // idle connections 10, Close connections after remaining idle for this duration(240sec).
		if err != nil {
			panic(err.Error()) // 嚴重錯誤
		}
		engine.Use(sessions.Sessions("release_sessions", store))

		{ // log
			f, _ := os.Create("gin.log")
			gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
			// defer f.Close()
			engine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
				// your custom format
				return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
					param.ClientIP,
					param.TimeStamp.Format(time.RFC1123),
					param.Method,
					param.Path,
					param.Request.Proto,
					param.StatusCode,
					param.Latency,
					param.Request.UserAgent(),
					param.ErrorMessage,
				)
			}))
		}
		//Release ...
		var group = engine.Group("/release")
		return engine, group
	}

	return engine, group
}
