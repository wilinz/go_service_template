package welcome

import (
	"github.com/gin-gonic/gin"
	"server_template/model"
	"server_template/service/account"
)

func WelcomeHandler(c *gin.Context) {
	logged, _ := account.IsLoggedWithResponse(c)
	if !logged {
		return
	}
	c.JSON(200, model.JsonResponse[any]{
		Code: 200,
		Msg:  "hello,you are logged in!",
		Data: nil,
	})
}
