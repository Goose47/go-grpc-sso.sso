package tests

import (
	ssov1 "github.com/Goose47/go-grpc-sso.protos/gen/go/sso"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-grpc-sso/tests/suite"
	"testing"
	"time"
)

const (
	appID     = 1
	appSecret = "test-secret"

	passDefaultLen = 10
)

func TestRegisterLogin_LoginHappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	regResp, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, regResp.UserId)

	loginResp, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    appID,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, loginResp.Token)

	token := loginResp.Token
	loginTime := time.Now()

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, regResp.UserId, claims["uid"].(int64))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, claims["app_id"].(int))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
