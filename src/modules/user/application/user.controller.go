package application

import (
	"github.com/gin-gonic/gin"
	"go-backend-template/src/api/utils"
)

type UserController struct {
	Engine  *gin.Engine
	Service UserService
}

func (r *UserController) Init() {
	r.Engine.GET("/user/test", r.login)

	r.Engine.GET("/user/add", r.addUser)
	r.Engine.GET("/users/me", r.authenticate, r.getMe)
	r.Engine.PUT("/users/me", r.authenticate, r.updateMe)
	r.Engine.PATCH("/users/me/password", r.authenticate, r.changeMyPassword)
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

	if err := utils.BindBody(&addUserDto, c); err != nil {
		utils.ErrorResponse(err, nil, true).Reply(c)
		return
	}

	user, err := r.Service.Add(utils.ContextWithReqInfo(c), addUserDto)

	if err != nil {
		utils.ErrorResponse(err, nil, true).Reply(c)
		return
	}

	utils.OkResponse(user).Reply(c)
}

func (r *UserController) updateMe(c *gin.Context) {
	var updateUserDto UpdateUserDto

	reqInfo := utils.GetReqInfo(c)
	updateUserDto.Id = reqInfo.UserId

	if err := utils.BindBody(&updateUserDto, c); err != nil {
		utils.ErrorResponse(err, nil, true).Reply(c)
		return
	}

	err := r.Service.Update(utils.ContextWithReqInfo(c), updateUserDto)
	if err != nil {
		utils.ErrorResponse(err, nil, true).Reply(c)
		return
	}

	utils.OkResponse(nil).Reply(c)
}

func (r *UserController) changeMyPassword(c *gin.Context) {
	var changeUserPasswordDto ChangeUserPasswordDto

	reqInfo := utils.GetReqInfo(c)
	changeUserPasswordDto.Id = reqInfo.UserId

	if err := utils.BindBody(&changeUserPasswordDto, c); err != nil {
		utils.ErrorResponse(err, nil, true).Reply(c)
		return
	}

	err := r.Service.ChangePassword(utils.ContextWithReqInfo(c), changeUserPasswordDto)
	if err != nil {
		utils.ErrorResponse(err, nil, true).Reply(c)
		return
	}

	utils.OkResponse(nil).Reply(c)
}

func (r *UserController) getMe(c *gin.Context) {
	reqInfo := utils.GetReqInfo(c)

	user, err := r.Service.GetById(utils.ContextWithReqInfo(c), reqInfo.UserId)
	if err != nil {
		utils.ErrorResponse(err, nil, true).Reply(c)
		return
	}

	utils.OkResponse(user).Reply(c)
}

type Test struct {
	Value int64 `json:"value"`
}

func (r *UserController) login(c *gin.Context) {
	utils.OkResponse(Test{Value: 1}).Reply(c)
}
