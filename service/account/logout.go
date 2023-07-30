package account

import (
	"github.com/gin-gonic/gin"
	"log"
	"server_template/common"
	"server_template/model"
	"server_template/service"
)

func LogoutHandler(c *gin.Context) {
	logged, _ := IsLoggedWithResponse(c)
	if !logged {
		return
	}
	session, err := common.Sessions.Get(c.Request, "session-key")
	session.Options.MaxAge = -1
	if err = session.Save(c.Request, c.Writer); err != nil {
		service.HttpServerInternalError(c)
		log.Println("failed deleting session: ", err)
		return
	}
	c.JSON(200, model.JsonResponse{
		Code: 200,
		Msg:  "ok",
		Data: "退出登录成功",
	})
}
