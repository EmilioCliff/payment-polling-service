package mock

import (
	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
)

var _ services.GrpcInterface = (*MockGrpcService)(nil)

type MockGrpcService struct {
	RegisterUserViagRPCFunc func(services.RegisterUserRequest) (int, services.RegisterUserResponse)
	LoginUserViagRPCFunc    func(services.LoginUserRequest) (int, services.LoginUserResponse)
}

func (m *MockGrpcService) RegisterUserViagRPC(req services.RegisterUserRequest) (int, services.RegisterUserResponse) {
	return m.RegisterUserViagRPCFunc(req)
}

func (m *MockGrpcService) LoginUserViagRPC(req services.LoginUserRequest) (int, services.LoginUserResponse) {
	return m.LoginUserViagRPCFunc(req)
}
