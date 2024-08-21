package api

import (
	"testing"

	db "github.com/EmilioCliff/payment-polling-app/authentication-service/db/sqlc"
	"github.com/EmilioCliff/payment-polling-app/authentication-service/utils"
	"github.com/stretchr/testify/require"
)

func newServerTest(t *testing.T, store *db.Queries) *Server {
	config, err := utils.LoadConfig("../")
	require.NoError(t, err)

	maker, err := utils.NewJWTMaker("../utils/my_rsa_key", "../utils/my_rsa_key.pub.pem")
	require.NoError(t, err)

	server, err := NewServer(config, store, *maker)
	require.NoError(t, err)

	return server
}
