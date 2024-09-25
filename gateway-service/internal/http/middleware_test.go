package http

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/EmilioCliff/payment-polling-app/gateway-service/pkg"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	req *http.Request,
	tokenMaker pkg.JWTMaker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	token, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	req.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestAuthenticationMiddleware(t *testing.T) {
	testCases := []struct {
		name             string
		setAuthorization func(req *http.Request, t *testing.T, tokenMaker pkg.JWTMaker)
		checkResponse    func(t *testing.T, recorder httptest.ResponseRecorder)
	}{
		{
			name: "valid_token",
			setAuthorization: func(req *http.Request, t *testing.T, tokenMaker pkg.JWTMaker) {
				addAuthorization(t, req, tokenMaker, "Bearer", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "missing_authorization_header",
			setAuthorization: func(_ *http.Request, _ *testing.T, _ pkg.JWTMaker) {
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "unsuported_authorization",
			setAuthorization: func(req *http.Request, t *testing.T, tokenMaker pkg.JWTMaker) {
				addAuthorization(t, req, tokenMaker, "Unsupported", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "missing_authorization_type",
			setAuthorization: func(req *http.Request, t *testing.T, tokenMaker pkg.JWTMaker) {
				addAuthorization(t, req, tokenMaker, "", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "expired_token",
			setAuthorization: func(req *http.Request, t *testing.T, tokenMaker pkg.JWTMaker) {
				addAuthorization(t, req, tokenMaker, "Bearer", "user", -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				panic(err)
			}

			publicKey := &privateKey.PublicKey

			testServer := NewHttpServer(pkg.JWTMaker{
				PublicKey:  publicKey,
				PrivateKey: privateKey,
			})
			testServer.router.GET(
				"/test-auth",
				authenticationMiddleware(testServer.maker),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, "/test-auth", nil)
			require.NoError(t, err)

			tc.setAuthorization(request, t, testServer.maker)
			testServer.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, *recorder)
		})
	}
}
