package mock

import (
	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/gin-gonic/gin"
)

var _ services.GrpcInterface = (*MockGrpcService)(nil)

type MockGrpcService struct {
	RegisterUserViagRPCFunc func(services.RegisterUserRequest) (int, gin.H)
	LoginUserViagRPCFunc    func(services.LoginUserRequest) (int, gin.H)
}

func (m *MockGrpcService) RegisterUserViagRPC(req services.RegisterUserRequest) (int, gin.H) {
	return m.RegisterUserViagRPCFunc(req)
}

func (m *MockGrpcService) LoginUserViagRPC(req services.LoginUserRequest) (int, gin.H) {
	return m.LoginUserViagRPCFunc(req)
}
