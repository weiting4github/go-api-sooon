/*Package config Server App協議常用設定 middleware 機密資料不要放這邊*/
package config

import (
	"fmt"
	"go-api-sooon/app"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"    // for "database/sql"
	_ "github.com/joho/godotenv/autoload" // 環境變數套件os.Getenv
)

const configCodePrefix = "CNF00"

// JWTClaims JWT帶入的使用者資訊
// Payload server端定義要帶入哪些常用的變數
type JWTClaims struct {
	Email    string `json:"Email"`
	Role     string `json:"Role"`
	MemberID int64  `json:"MemberID"`
	Lang     string `json:"Lang"`
	jwt.StandardClaims
}

// JwtSecret SecretKey
var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

// CreateJWTClaims 簽發JWT token
func CreateJWTClaims(memberID int64, email string, role string, issuer string) (JWTClaims, string, error) {
	now := time.Now()
	jwtID := email + strconv.FormatInt(now.Unix(), 10)

	// set claims and sign
	claims := JWTClaims{
		Email:    email,
		Role:     role,
		MemberID: memberID,
		StandardClaims: jwt.StandardClaims{
			Audience:  email,
			ExpiresAt: now.Add(3600 * time.Second).Unix(), // 過期時間
			Id:        jwtID,
			IssuedAt:  now.Unix(), // 發行時間
			Issuer:    issuer,
			NotBefore: now.Add(1 * time.Second).Unix(), // 幾秒後可以開始使用
			Subject:   email,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(JwtSecret)
	return claims, token, err
}

// JWTAuth JWT middleware payload內容設定環境變數
func JWTAuth(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	bearerString := strings.Split(auth, "Bearer ")
	if len(bearerString) < 2 {
		c.JSON(http.StatusOK, gin.H{
			"s":       -9,
			"errMsg":  "no Bearer",
			"errCode": app.SFunc.DumpErrorCode(configCodePrefix),
		})
		return
	}
	// fmt.Println(len(bearerString))
	// fmt.Println(bearerString)
	// fmt.Println(os.Getenv("JWT_SECRET"))
	token := strings.Split(auth, "Bearer ")[1]
	// parse and validate token for six things:
	// validationErrorMalformed => token is malformed
	// validationErrorUnverifiable => token could not be verified because of signing problems
	// validationErrorSignatureInvalid => signature validation failed
	// validationErrorExpired => exp validation failed
	// validationErrorNotValidYet => nbf validation failed
	// validationErrorIssuedAt => iat validation failed
	tokenClaims, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		// JwtSecret 取得驗證
		app.SFunc.DumpAnything(token.Claims)
		// sample token is expired.  override time so it parses as valid?
		return JwtSecret, nil
	})

	if err != nil {
		var message string
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				message = "token is malformed"
			} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
				message = "token could not be verified because of signing problems"
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				message = "signature validation failed"
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				message = "token is expired"
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				message = "token is not yet valid before sometime"
			} else {
				message = "can not handle this token"
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"s":     -1,
			"error": message,
		})
		c.Abort()
		return
	}

	// 設定環境變數
	if claims, ok := tokenClaims.Claims.(*JWTClaims); ok && tokenClaims.Valid {
		fmt.Println("email:", claims.Email)
		fmt.Println("role:", claims.Role)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("memberID", claims.MemberID)
		c.Set("lang", claims.Lang)
		c.Next()
	} else {
		c.Abort()
		return
	}
}

// MemberSessions 存放redis sessions的用戶結構
type MemberSessions struct {
	LoginTs int64
	Lang    string
	Email   string
}

// CORSMiddleware 允許CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
