// Package jwt 处理 JWT 认证
package jwt

import (
	"github.com/pkg/errors"
	"gohub/pkg/app"
	"gohub/pkg/config"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/golang-jwt/jwt"
)

var (
	ErrTokenExpired    = errors.New("令牌已过期")
	ErrTokenMalformed  = errors.New("请求令牌格式有误")
	ErrTokenInvalid    = errors.New("请求令牌无效")
	ErrHeaderEmpty     = errors.New("需要认证才能访问！")
	ErrHeaderMalformed = errors.New("请求头中 Authorization 格式有误")
	ErrTokenCreateFail = errors.New("创建 token 失败")
	ErrTokenTypeFail   = errors.New("token 类型错误")
)

// JWT 定义一个jwt对象
type JWT struct {

	// 秘钥，用以加密 JWT，读取配置信息 app.key
	SignKey []byte

	// 刷新 Token 的最大过期时间
	MaxRefresh int64
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// CustomClaims 自定义载荷
type CustomClaims struct {
	UserID    string `json:"user_id"`
	IsRefresh bool   `json:"is_refresh"` // 是否是刷新 Token
	// StandardClaims 结构体实现了 Claims 接口继承了  Valid() 方法
	// JWT 规定了7个官方字段，提供使用:
	// - iss (issuer)：发布者
	// - sub (subject)：主题
	// - iat (Issued At)：生成签名的时间
	// - exp (expiration time)：签名过期时间
	// - aud (audience)：观众，相当于接受者
	// - nbf (Not Before)：生效时间
	// - jti (JWT ID)：编号
	jwtpkg.StandardClaims
}

func NewJWT() *JWT {
	return &JWT{
		SignKey:    []byte(config.GetString("app.key")),
		MaxRefresh: config.GetInt64("jwt.max_refresh_time"),
	}
}

// ParserToken 解析 Token，中间件中调用
func (jwt *JWT) ParserToken(c *gin.Context) (*CustomClaims, error) {

	tokenString, parseErr := jwt.getTokenFromHeader(c)
	if parseErr != nil {
		return nil, parseErr
	}

	// 1. 调用 jwt 库解析用户传参的 Token
	token, err := jwt.parseTokenString(tokenString)

	// 2. 解析出错
	if err != nil {
		var validationErr *jwtpkg.ValidationError
		ok := errors.As(err, &validationErr)
		if ok {
			if validationErr.Errors == jwtpkg.ValidationErrorMalformed {
				return nil, ErrTokenMalformed
			} else if validationErr.Errors == jwtpkg.ValidationErrorExpired {
				return nil, ErrTokenExpired
			}
		}
		return nil, ErrTokenInvalid
	}

	// 3. 将 token 中的 claims 信息解析出来和 CustomClaims 数据结构进行校验
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// RefreshToken 更新 Token，用以提供 refresh token 接口
func (jwt *JWT) RefreshToken(c *gin.Context) (Tokens, error) {

	// 1. 获取 token
	token, err := jwt.ParserToken(c)
	if err != nil {
		return Tokens{}, err
	}

	if !token.IsRefresh {
		return Tokens{}, ErrTokenTypeFail
	}
	return jwt.IssueToken(token.UserID)
}

// IssueToken 生成  Token，在登录成功时调用
func (jwt *JWT) IssueToken(userID string) (Tokens, error) {

	var expireTime int64
	if config.GetBool("app.debug") {
		expireTime = config.GetInt64("jwt.debug_expire_time")
	} else {
		expireTime = config.GetInt64("jwt.expire_time")
	}

	expireAtTime := jwt.expireAtTime(expireTime)
	accessToken, err := jwt.createToken(userID, expireAtTime, false)
	if err != nil {
		return Tokens{}, ErrTokenCreateFail
	}

	refreshExpireAtTime := jwt.expireAtTime(jwt.MaxRefresh)
	refreshAccessToken, err := jwt.createToken(userID, refreshExpireAtTime, true)
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{accessToken, refreshAccessToken}, nil
}

// createToken 创建 Token，内部使用，外部请调用 IssueToken
func (jwt *JWT) createToken(userID string, expiresAt int64, isRefresh bool) (string, error) {
	// 使用HS256算法进行token生成
	token := jwtpkg.NewWithClaims(jwtpkg.SigningMethodHS256, CustomClaims{
		userID,
		isRefresh,
		jwtpkg.StandardClaims{
			NotBefore: app.TimenowInTimezone().Unix(), // 签名生效时间
			ExpiresAt: expiresAt,                      // 签名过期时间
			Issuer:    config.GetString("app.name"),   // 签名颁发者
		},
	})
	return token.SignedString(jwt.SignKey)
}

// expireAtTime 过期时间
func (jwt *JWT) expireAtTime(expireTime int64) int64 {
	timenow := app.TimenowInTimezone()
	expire := time.Duration(expireTime) * time.Minute
	return timenow.Add(expire).Unix()
}

// parseTokenString 使用 jwtpkg.ParseWithClaims 解析 Token
func (jwt *JWT) parseTokenString(tokenString string) (*jwtpkg.Token, error) {
	return jwtpkg.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwtpkg.Token) (interface{}, error) {
		return jwt.SignKey, nil
	})
}

// getTokenFromHeader 使用 jwtpkg.ParseWithClaims 解析 Token
// Authorization:Bearer xxxxx
func (jwt *JWT) getTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrHeaderEmpty
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", ErrHeaderMalformed
	}
	return parts[1], nil
}
