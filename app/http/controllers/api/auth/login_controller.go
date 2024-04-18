package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gohub/app/requests"
	"gohub/pkg/eth"
	"gohub/pkg/jwt"
	"gohub/pkg/logger"
	"gohub/pkg/redis"
	"gohub/pkg/response"
)

type LoginController struct {
}

// GetMessage 得到签名所需要的 Message
func (lc *LoginController) GetMessage(c *gin.Context) {
	// 1. 验证表单
	request := GetMessageReq{}
	if ok := requests.Validate(c, &request, GetMessageVal); !ok {
		return
	}
	address := request.Address

	// 2. 创建 nonce
	nonce := uuid.New().String()

	// 3. 组成 message

	//message := fmt.Sprintf("Welcome to %s!\n\n"+
	//	"Wallet address:\n%s\n\n"+
	//	"Nonce:\n%s", config.Get("app.name"), address, nonce)

	message := fmt.Sprintf("Nonce:%s", nonce)

	// 4. 保存到redis
	redisKey := redis.AUTH_NONCE
	redis.Redis.Set(redisKey.SKey(address), message, redisKey.Expired)

	response.SuccessData(c, gin.H{
		"message": message,
		"address": address,
	})
}

// LoginBySignature 签名登录
func (lc *LoginController) LoginBySignature(c *gin.Context) {
	// 1. 验证表单
	request := LoginBySignatureReq{}
	if ok := requests.Validate(c, &request, LoginBySignatureVal); !ok {
		return
	}
	// 2. 验证 nonce
	redisKey := redis.AUTH_NONCE
	val := redis.Redis.Get(redisKey.SKey(request.Address))
	if val == "" {
		response.ErrorStr(c, "message已失效，请重新获取")
		return
	}

	// 3. 验证签名
	err := eth.VerifySignature(request.Address, val, request.Signature)
	if err != nil {
		logger.Errorf("%+v", err)
		response.ErrorStr(c, "签名验证失败")
		return
	}

	// 4. 生成 jwt token
	token, err := jwt.NewJWT().IssueToken(request.Address)
	if err != nil {
		response.Error(c, errors.Wrap(err, "生成 token 失败"))
		return
	}

	response.SuccessData(c, token)
}

// RefreshToken 刷新 Access Token
func (lc *LoginController) RefreshToken(c *gin.Context) {
	tokens, err := jwt.NewJWT().RefreshToken(c)
	if err != nil {
		logger.Errorf("刷新 token 失败:%+v", err)
		response.Error(c, err)
	} else {
		response.SuccessData(c, tokens)
	}
}
