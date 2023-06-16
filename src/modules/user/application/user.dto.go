package application

type AddUserDto struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUserDto struct {
	Id        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type ChangeUserPasswordDto struct {
	Id       int64  `json:"id"`
	Password string `json:"password"`
}
