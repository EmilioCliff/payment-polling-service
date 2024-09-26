package services

type GrpcInterface interface {
	RegisterUserViagRPC(RegisterUserRequest) (int, RegisterUserResponse)
	LoginUserViagRPC(LoginUserRequest) (int, LoginUserResponse)
}
