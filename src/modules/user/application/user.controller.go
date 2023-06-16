package application

import (
	"github.com/gin-gonic/gin"
	http "go-backend-template/src/utils/http"
)

type UserController struct {
	*http.Router
	Service UserService
}

func (r *UserController) Init() {
	r.RouterGroup.GET("/add", r.addUser)
	r.RouterGroup.GET("/users/me", r.authenticate, r.getMe)
	r.RouterGroup.PUT("/users/me", r.authenticate, r.updateMe)
	r.RouterGroup.PATCH("/users/me/password", r.authenticate, r.changeMyPassword)
}

func (r *UserController) authenticate(c *gin.Context) {
	//token := c.Request.Header.Get("Authorization")

	//userId, err := r.authService.VerifyAccessToken(token)
	//if err != nil {
	//	response := utils.ErrorResponse(err, nil, true)
	//	c.AbortWithStatusJSON(response.Status, response)
	//}
	//
	//utils.SetUserId(c, userId)
}

func (r *UserController) addUser(c *gin.Context) {
	var addUserDto AddUserDto

	if err := http.BindBody(&addUserDto, c); err != nil {
		http.ErrorResponse(err, nil, true).Reply(c)
		return
	}

	user, err := r.Service.Add(http.ContextWithReqInfo(c), addUserDto)

	if err != nil {
		http.ErrorResponse(err, nil, true).Reply(c)
		return
	}

	http.OkResponse(user).Reply(c)
}

func (r *UserController) updateMe(c *gin.Context) {
	var updateUserDto UpdateUserDto

	reqInfo := http.GetReqInfo(c)
	updateUserDto.Id = reqInfo.UserId

	if err := http.BindBody(&updateUserDto, c); err != nil {
		http.ErrorResponse(err, nil, true).Reply(c)
		return
	}

	err := r.Service.Update(http.ContextWithReqInfo(c), updateUserDto)
	if err != nil {
		http.ErrorResponse(err, nil, true).Reply(c)
		return
	}

	http.OkResponse(nil).Reply(c)
}

func (r *UserController) changeMyPassword(c *gin.Context) {
	var changeUserPasswordDto ChangeUserPasswordDto

	reqInfo := http.GetReqInfo(c)
	changeUserPasswordDto.Id = reqInfo.UserId

	if err := http.BindBody(&changeUserPasswordDto, c); err != nil {
		http.ErrorResponse(err, nil, true).Reply(c)
		return
	}

	err := r.Service.ChangePassword(http.ContextWithReqInfo(c), changeUserPasswordDto)
	if err != nil {
		http.ErrorResponse(err, nil, true).Reply(c)
		return
	}

	http.OkResponse(nil).Reply(c)
}

func (r *UserController) getMe(c *gin.Context) {
	reqInfo := http.GetReqInfo(c)

	user, err := r.Service.GetById(http.ContextWithReqInfo(c), reqInfo.UserId)
	if err != nil {
		http.ErrorResponse(err, nil, true).Reply(c)
		return
	}

	http.OkResponse(user).Reply(c)
}
