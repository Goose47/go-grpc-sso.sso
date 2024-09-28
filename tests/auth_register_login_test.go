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

	assert.Equal(t, regResp.UserId, int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	regResp1, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, regResp1.UserId)

	regResp2, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.Error(t, err)
	assert.Empty(t, regResp2)
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegister_FailCases(t *testing.T) {
	testCases := []struct {
		name        string
		email       string
		pass        string
		expectedErr string
	}{
		{
			name:        "Register with empty password",
			email:       gofakeit.Email(),
			pass:        "",
			expectedErr: "password is required",
		},
		{
			name:        "Register with empty email",
			email:       "",
			pass:        randomFakePassword(),
			expectedErr: "email is required",
		},
		{
			name:        "Register with both empty",
			email:       "",
			pass:        "",
			expectedErr: "email is required",
		},
	}
	ctx, st := suite.New(t)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			regResp, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tt.email,
				Password: tt.pass,
			})

			require.Error(t, err)
			assert.Empty(t, regResp)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	testCases := []struct {
		name        string
		email       string
		pass        string
		appID       int32
		expectedErr string
	}{
		{
			name:        "Login with empty email",
			email:       "",
			pass:        randomFakePassword(),
			appID:       appID,
			expectedErr: "email is required",
		},
		{
			name:        "Login with empty pass",
			email:       gofakeit.Email(),
			pass:        "",
			appID:       appID,
			expectedErr: "password is required",
		},
		{
			name:        "Login with empty app id",
			email:       gofakeit.Email(),
			pass:        randomFakePassword(),
			appID:       0,
			expectedErr: "app id is required",
		},
		{
			name:        "Login with invalid credentials",
			email:       gofakeit.Email(),
			pass:        randomFakePassword(),
			appID:       appID,
			expectedErr: "invalid credentials",
		},
	}

	ctx, st := suite.New(t)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			loginResp, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.pass,
				AppId:    tt.appID,
			})

			require.Error(t, err)
			assert.Empty(t, loginResp)
			assert.ErrorContains(t, err, tt.expectedErr)
		})
	}
}
