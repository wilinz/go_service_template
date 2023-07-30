package account

import (
	"github.com/gin-gonic/gin"
	"server_template/db"
	"server_template/model"
	"server_template/service"
)

func SetUserInfoHandler(c *gin.Context) {
	var p model.UserInfo
	err := c.Bind(&p)
	if err != nil {
		service.HttpParameterError(c)
		return
	}

	logged, username := IsLoggedWithResponse(c)
	if !logged {
		return
	}

	err = db.Mysql.Where(map[string]any{"username": username}).Updates(&model.User{
		UserInfoReadOnly: model.UserInfoReadOnly{UserInfo: p},
	}).Error
	if err != nil {
		service.HttpServerInternalError(c)
		return
	}
	c.JSON(200, model.JsonResponse{
		Code: 200,
		Msg:  "ok",
		Data: nil,
	})
}

func GetUserInfoHandler(c *gin.Context) {
	logged, username := IsLoggedWithResponse(c)
	if !logged {
		return
	}

	var info model.UserInfoReadOnly
	err := db.Mysql.Model(&model.User{}).Where(map[string]any{"username": username}).Take(&info).Error
	if err != nil {
		service.HttpServerInternalError(c)
		return
	}
	c.JSON(200, model.JsonResponse{
		Code: 200,
		Msg:  "ok",
		Data: info,
	})
}
