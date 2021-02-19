// Package app api程式常用小工具
package app

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Loadi18n 抓不到語言預設使用en
func Loadi18n(c *gin.Context) *i18n.Localizer {
	lang := c.Request.FormValue("lang")
	if len(lang) <= 0 {
		if v, ok := c.Get("lang"); ok {
			lang = v.(string)
		} else {
			lang = "zh"
		}
	}
	bundle := i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFile("config/locales/" + lang + ".toml")

	localizer := i18n.NewLocalizer(bundle, lang)

	return localizer
}

// APIAuthorized ...
func APIAuthorized() string {
	h := sha256.New()
	h.Write([]byte(os.Getenv("APP_NAME")))
	return hex.EncodeToString(h.Sum(nil))
}

// GetTokenVia ...
func GetTokenVia(c *gin.Context) {
	// validate request body
	var body struct {
		Account  string
		Password string
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
}

/*DumpErrorCode debug印出function name*/
func DumpErrorCode(codePrefix string) string {
	_, _, fileLine, _ := runtime.Caller(1)
	// return runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("%s_%d", codePrefix, fileLine)
}

// DumpAnyLikeABoss dumping your stuff like a boss
func DumpAnyLikeABoss(i interface{}) {
	fmt.Println("-------------------------")
	fmt.Printf("%#v\n", i)
	fmt.Println("-------------------------")
}

// NewMd5String BD密碼salt產生器 會產生len * 2長度的字串 最多16*2
func NewMd5String(len int) string {
	t := strconv.FormatInt(time.Now().Unix(), 10) // int64 to int to string
	b := []byte(t)
	m := md5.Sum(b)
	final := hex.EncodeToString(m[0:len])

	return final
}