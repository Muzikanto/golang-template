package dto

type CreateUserDto struct {
	Test string `json:"test" form:"test" binding:"required"`
}
