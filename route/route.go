package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"server_template/common"
	_ "server_template/docs"
	"server_template/service/account"
	"server_template/service/appversion"
	"server_template/service/feedback"
	"server_template/service/proxy"
	"server_template/service/welcome"
)

func Run(host string, port int) {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	r.Use(func(c *gin.Context) {
		/*	log.Println(c.Request.Header)
			c.Next()
			log.Println(c.Writer.Header())*/
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Static(common.UploadDir, common.UploadSavePath)
	accountRouter := r.Group("/account")
	accountRouter.POST("/login", account.LoginHandler)
	accountRouter.POST("/login_with_code", account.LoginWithCodeHandler)
	accountRouter.POST("/register", account.RegisterHandler)
	accountRouter.POST("/verify", account.VerificationCodeHandler)
	accountRouter.PUT("/password/reset", account.ResetPasswordHandler)
	accountRouter.PUT("/password/change", account.ChangePasswordHandler)
	accountRouter.PUT("/info", account.SetUserInfoHandler)
	accountRouter.GET("/info", account.GetUserInfoHandler)
	accountRouter.DELETE("/logout", account.LogoutHandler)

	r.GET("/welcome", welcome.WelcomeHandler)

	r.POST("/feedback", feedback.PostFeedbackHandler)

	r.GET("/app_version", appversion.GetAppVersion)

	r.Any("/proxy", proxy.HttpProxyHandler)

	err := r.Run(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}
}
