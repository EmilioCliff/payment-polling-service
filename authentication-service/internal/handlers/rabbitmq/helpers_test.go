package rabbitmq_test

import (
	"net/http"
	"testing"

	"github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
)

func TestGRPCServer_convertPkgError (t *testing.T) {
	r := NewTestRabbitConn()

	tests := []struct{
		name string
		err *pkg.Error
		want int
	}{
		{
			name: "already_exists",
			err: &pkg.Error{
				Code: pkg.ALREADY_EXISTS_ERROR,
			},
			want: http.StatusConflict,
		},
		{
			name: "internal_error",
			err: &pkg.Error{
				Code: pkg.INTERNAL_ERROR,
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "invalid_error",
			err: &pkg.Error{
				Code: pkg.INVALID_ERROR,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "not_found",
			err: &pkg.Error{
				Code: pkg.NOT_FOUND_ERROR,
			},
			want: http.StatusNotFound,
		},
		{
			name: "not_implemented",
			err: &pkg.Error{
				Code: pkg.NOT_IMPLEMENTED_ERROR,
			},
			want: http.StatusNotImplemented,
		},
		{
			name: "authentication_error",
			err: &pkg.Error{
				Code: pkg.AUTHENTICATION_ERROR,
			},
			want: http.StatusUnauthorized,
		},
		{
			name: "default",
			err: &pkg.Error{
				Code: "system_error",
			},
			want: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := r.rabbitConn.ConvertPkgError(tc.err)
			if got != tc.want {
				t.Errorf("convertPkgError() = %v, want %v", got, tc.want)
			}
		})
	}
}