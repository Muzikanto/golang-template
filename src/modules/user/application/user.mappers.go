package application

import "go-backend-template/src/modules/user/domain"

func (dto AddUserDto) ToDomain() (domain.UserDomain, error) {
	return domain.NewUser(
		dto.FirstName,
		dto.LastName,
		dto.Email,
		dto.Password,
	)
}
