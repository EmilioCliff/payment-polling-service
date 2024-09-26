package mock

import (
	"github.com/EmilioCliff/payment-polling-app/gateway-service/internal/services"
)

var _ services.HttpInterface = (*MockHttpService)(nil)

type MockHttpService struct {
	RegisterUserViaHttpFunc func(services.RegisterUserRequest) (int, services.RegisterUserResponse)
	LoginUserViaHttpFunc    func(services.LoginUserRequest) (int, services.LoginUserResponse)
}

func (m *MockHttpService) RegisterUserViaHttp(req services.RegisterUserRequest) (int, services.RegisterUserResponse) {
	return m.RegisterUserViaHttpFunc(req)
}

func (m *MockHttpService) LoginUserViaHttp(req services.LoginUserRequest) (int, services.LoginUserResponse) {
	return m.LoginUserViaHttpFunc(req)
}
