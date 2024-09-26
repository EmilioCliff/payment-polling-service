package services

type HttpInterface interface {
	RegisterUserViaHttp(RegisterUserRequest) (int, RegisterUserResponse)
	LoginUserViaHttp(LoginUserRequest) (int, LoginUserResponse)
}
