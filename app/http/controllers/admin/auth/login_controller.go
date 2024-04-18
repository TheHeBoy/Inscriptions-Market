package auth

import (
	"gohub/app/requests"
	"gohub/pkg/auth"
	"gohub/pkg/jwt"
	"gohub/pkg/logger"
	"gohub/pkg/response"

	"github.com/gin-gonic/gin"
)

// LoginController 用户控制器
type LoginController struct{}

// LoginByPhone 手机登录
func (lc *LoginController) LoginByPhone(c *gin.Context) {

	// 1. 验证表单
	request := requests.LoginByPhoneRequest{}
	if ok := requests.Validate(c, &request, requests.LoginByPhone); !ok {
		return
	}

	// 2. 尝试登录
	user, err := auth.LoginByPhone(request.Phone)
	if err != nil {
		// 失败，显示错误提示
		response.Error(c, err)
	} else {
		// 登录成功
		tokens, err := jwt.NewJWT().IssueToken(user.GetStringID())
		if err != nil {
			logger.Error(err)
			response.Error(c, err)
		} else {
			response.SuccessData(c, tokens)
		}
	}
}

// LoginByPassword 多种方法登录，支持手机号、email 和用户名
func (lc *LoginController) LoginByPassword(c *gin.Context) {
	// 1. 验证表单
	request := requests.LoginByPasswordRequest{}
	if ok := requests.Validate(c, &request, requests.LoginByPassword); !ok {
		return
	}

	// 2. 尝试登录
	user, err := auth.Attempt(request.LoginID, request.Password)
	if err != nil {
		// 失败，显示错误提示
		response.Error(c, err)
	} else {
		tokens, err := jwt.NewJWT().IssueToken(user.GetStringID())
		if err != nil {
			logger.Error(err)
			response.Error(c, err)
		} else {
			response.SuccessData(c, tokens)
		}
	}
}

// RefreshToken 刷新 Access Token
func (lc *LoginController) RefreshToken(c *gin.Context) {
	tokens, err := jwt.NewJWT().RefreshToken(c)
	if err != nil {
		response.Error(c, err)
	} else {
		response.SuccessData(c, tokens)
	}
}
