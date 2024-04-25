package repository

type CreateUserResponse struct {
	*User
}

type ExistsUserResponse struct {
	Exists bool
}

type GetUserResponse struct {
	*User
}

type UpdateUserResponse struct {
	*User
}

type DeleteUserResponse struct {
	*User
}

type CreateAddressResponse struct {
	*UserAddress
}

type GetAddressResponse struct {
	*UserAddress
}

type UpdateAddressResponse struct {
	*UserAddress
}

type DeleteAddressResponse struct {
	*UserAddress
}

