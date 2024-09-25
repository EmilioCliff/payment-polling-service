package mock

import (
	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
	"github.com/gin-gonic/gin"
)

var _ services.HttpInterface = (*MockHttpService)(nil)

type MockHttpService struct {
	RegisterUserViaHttpFunc func(services.RegisterUserRequest) (int, gin.H)
	LoginUserViaHttpFunc    func(services.LoginUserRequest) (int, gin.H)
}

func (m *MockHttpService) RegisterUserViaHttp(req services.RegisterUserRequest) (int, gin.H) {
	return m.RegisterUserViaHttpFunc(req)
}

func (m *MockHttpService) LoginUserViaHttp(req services.LoginUserRequest) (int, gin.H) {
	return m.LoginUserViaHttpFunc(req)
}
