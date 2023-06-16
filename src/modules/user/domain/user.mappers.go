package domain

func (dto UserDto) ToResponse(user UserDomain) UserDto {
	dto.Id = user.Id
	dto.FirstName = user.FirstName
	dto.LastName = user.LastName
	dto.Email = user.Email

	return dto
}
