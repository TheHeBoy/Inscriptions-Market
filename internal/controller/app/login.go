package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gohub/internal/redisI"
	"gohub/internal/request/app"
	"gohub/internal/request/validators"
	"gohub/pkg/config"
	"gohub/pkg/eth"
	"gohub/pkg/jwt"
	"gohub/pkg/logger"
	"gohub/pkg/redisP"
	"gohub/pkg/response"
)

type LoginController struct {
}

// GetMessageAuth 得到签名所需要的 Message
func (lc *LoginController) GetMessageAuth(c *gin.Context) {
	// 1. 验证表单
	req := app.GetMessageReq{}
	if ok := validators.Validate(c, &req); !ok {
		return
	}
	address := req.Address

	// 2. 创建 nonce
	nonce := uuid.New().String()

	// 3. 组成 message
	message := fmt.Sprintf("Welcome to %s! Nonce:%s", config.Get("app.name"), nonce)

	// 4. 保存到redis
	redisKey := redisI.AuthNonce
	redisP.Redis.Set(redisKey.SKey(address), message, redisKey.Expired)

	response.SuccessData(c, gin.H{
		"message": message,
		"address": address,
	})
}

// LoginBySignatureAuth 签名登录
func (lc *LoginController) LoginBySignatureAuth(c *gin.Context) {
	// 1. 验证表单
	req := app.LoginBySignatureReq{}
	if ok := validators.Validate(c, &req); !ok {
		return
	}

	// 2. 验证 nonce
	redisKey := redisI.AuthNonce
	val := redisP.Redis.Get(redisKey.SKey(req.Address))
	if val == "" {
		response.ErrorStr(c, "message已失效，请重新获取")
		return
	}

	// 3. 验证签名
	err := eth.VerifySignature(req.Address, val, req.Signature)
	if err != nil {
		logger.Errorv(err)
		response.ErrorStr(c, "签名验证失败")
		return
	}

	// 4. 生成 jwt token
	token, err := jwt.NewJWT().IssueToken(req.Address)
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
