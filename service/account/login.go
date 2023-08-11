package account

import (
	"github.com/gin-gonic/gin"
	my_sessions "server_template/common"
	"server_template/constant/code_type"
	"server_template/constant/redis_prefix"
	"server_template/db"
	"server_template/model"
	"server_template/service"
	"server_template/util"
)

// LoginHandler
// @Summary Login
// @Description Handles user login with username and password
// @Tags Account
// @Accept json
// @Produce json
// @Param loginRequest body model.LoginParameters true "Login Request"
// @Success 200 {object} model.JsonResponse[any]
// @Failure 400 {object} model.JsonResponse[any]
// @Header 200 {string} Set-Cookie "session-key=<session_key>; Path=/"
// @Router /account/login [post]
func LoginHandler(c *gin.Context) {
	var user model.LoginParameters
	err := c.Bind(&user)
	if err != nil {
		service.HttpParameterError(c)
		return
	}
	if IsPasswordCorrect(c, user.Username, user.Password) {
		HandleLogin(c, user.Username)
	}
}

// LoginWithCodeHandler
// @Summary Login with verification code
// @Description Handles user login with verification code
// @Tags Account
// @Accept json
// @Produce json
// @Param loginWithCodeRequest body model.UseCodeLoginParameters true "Login With Code Request"
// @Success 200 {object} model.JsonResponse[any]
// @Failure 400 {object} model.JsonResponse[any]
// @Header 200 {string} Set-Cookie "session-key=<session_key>; Path=/"
// @Router /account/login_with_code [post]
func LoginWithCodeHandler(c *gin.Context) {
	var user model.UseCodeLoginParameters
	err := c.Bind(&user)
	if err != nil {
		service.HttpParameterError(c)
		return
	}

	if IsVerificationCodeCorrect(c, user.VerificationCode, code_type.Login, user.Username) {
		key := util.GetKey(redis_prefix.PasswordTryCount, user.Username)
		db.Redis.Del(db.Context, key)
		HandleLogin(c, user.Username)
	}
}

func HandleLogin(c *gin.Context, username string) {
	//expirationTime := time.Now().Add(1 * time.Minute)
	session, err := my_sessions.Sessions.Get(c.Request, "session-key")
	if err != nil {
		service.HttpServerInternalError(c)
		return
	}
	session.Values["username"] = username

	err = my_sessions.Sessions.Save(c.Request, c.Writer, session)
	if err != nil {
		service.HttpServerInternalError(c)
		return
	}

	service.HttpOK(c)
}
