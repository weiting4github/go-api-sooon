// Package app api程式常用小工具
package app

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
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

// IDump ...
type IDump interface {
	DumpAnything(i interface{})
}

// ErrCodePrefix 錯誤代碼追蹤
var ErrCodePrefix = "APP00"

// SooonFunc ...
type SooonFunc struct {
	APIAuthorizedKey string
}

// SFunc ...
var SFunc *SooonFunc

func init() {
	SFunc = &SooonFunc{APIAuthorizedKey: ""}
}

// APIAuthorized ...
func (s *SooonFunc) APIAuthorized() string {
	if s.APIAuthorizedKey == "" {
		h := sha256.New()
		h.Write([]byte(os.Getenv("APP_NAME")))
		return hex.EncodeToString(h.Sum(nil))
	}
	return s.APIAuthorizedKey

}

/*DumpErrorCode debug印出function name*/
func (s *SooonFunc) DumpErrorCode(codePrefix string) string {
	_, _, fileLine, _ := runtime.Caller(1)
	// return runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("%s_%d", codePrefix, fileLine)
}

// DumpAnything dumping your stuff like a boss
func (s *SooonFunc) DumpAnything(i interface{}) {
	fmt.Println("-------------------------")
	fmt.Printf("%#v\n", i)
	fmt.Println("-------------------------")
}

// NewMd5String BD密碼salt產生器 會產生len * 2長度的字串 最多16*2
func (s *SooonFunc) NewMd5String(len int) string {
	t := strconv.FormatInt(time.Now().Unix(), 10) // int64 to int to string
	b := []byte(t)
	m := md5.Sum(b)
	final := hex.EncodeToString(m[0:len])

	return final
}

// Dump implement IDump interface
func (s *SooonFunc) Dump(d IDump, i interface{}) {
	d.DumpAnything(i)
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

// Loadi18n 多語設定
// 抓不到語言預設使用en
func (s *SooonFunc) loadi18n(c *gin.Context) *i18n.Localizer {
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

// Localizer 返回翻譯後的多語
func (s *SooonFunc) Localizer(c *gin.Context, outputMsg string) string {
	// 語系
	localizer := s.loadi18n(c)
	translation, _ := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: outputMsg,
	})
	return translation
}

// Init APP基本認證 過了APP才能去要JWT token
// @Summary  APP基本認證 過了APP才能去要JWT token
// @Tags App
// @param hash formData string true "sha256"
// @version 1.0
// @produce json
// @Success 200 {object} initSuccessResponse "登入紀錄"
// @Failure 400 {object} apiFailResponse
// @host localhost:3000
// @Router /init [post]
func (s *SooonFunc) Init(c *gin.Context) {
	h := c.PostForm("hash")
	if h != s.APIAuthorized() {
		c.JSON(http.StatusBadRequest, apiFailResponse{
			S:       -9,
			ErrCode: SFunc.DumpErrorCode(ErrCodePrefix),
			ErrMsg:  errors.New("unauthorized").Error(),
		})

		return
	}
	c.JSON(http.StatusOK, initSuccessResponse{
		S: 1,
	})
}
