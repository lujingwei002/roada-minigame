package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shark/minigame-common/db"
)

type AccountController struct {
}

func NewAccountController() *AccountController {
	return new(AccountController)
}

type AccountLoginRequest struct {
	Openid   string `form:"openid" binding:"required"`
	Avatar   string `form:"avatar" binding:""`
	Nickname string `form:"nickname" binding:""`
	Token    string `form:"token" binding:"required"`
	Time     int64  `form:"time" binding:"required"`
}

type AccountLoginResponse struct {
	Code int    `json:"code"`
	Host string `json:"host"`
	Port int32  `json:"port"`
}

type AccountTestRequest struct {
}

type AccountTestResponse struct {
	Code int `json:"code"`
}

// @Summary 测试登陆
// @Description
// @Tags account
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Param userid formData int64 true "用户id" default(2)
// @Param code formData string true "密钥"
// @Success 0 {object} AccountDebugLoginResponse
// @Router /account/debug_login [post]
func (ac *AccountController) Login(c *gin.Context) {
	var req = AccountLoginRequest{}
	if err := c.ShouldBind(&req); err != nil {
		c.PureJSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "msg": err.Error()})
		return
	}
	var resp = AccountLoginResponse{
		Code: 0,
	}
	gate, err := db.GateStat_Get()
	if err != nil {
		log.Printf("[account] GateStat_Get failed, error=%s\n", err.Error())
		resp.Code = 1
		c.PureJSON(http.StatusOK, gin.H{"code": 0, "data": resp})
		return
	}
	if gate == nil {
		resp.Code = 1
		c.PureJSON(http.StatusOK, gin.H{"code": 0, "data": resp})
		return
	}
	resp.Host = gate.Host
	resp.Port = gate.Port
	c.PureJSON(http.StatusOK, gin.H{"code": 0, "data": resp})
}

func (ac *AccountController) Test(c *gin.Context) {
	var req = AccountTestRequest{}
	if err := c.ShouldBind(&req); err != nil {
		c.PureJSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "msg": err.Error()})
		return
	}
	var resp = AccountTestResponse{
		Code: 0,
	}
	c.PureJSON(http.StatusOK, gin.H{"code": 0, "data": resp})
}
