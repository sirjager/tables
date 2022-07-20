package tokens

import (
	"testing"
	"time"

	"github.com/SirJager/tables/service/core/utils"
	"github.com/stretchr/testify/require"
)

func TestPasetoBuilder(t *testing.T) {
	builder, err := NewPasetoBuilder(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomUserName()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := builder.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = builder.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.Id)
	require.Equal(t, username, payload.User)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	builder, err := NewPasetoBuilder(utils.RandomString(32))
	require.NoError(t, err)

	token, payload, err := builder.CreateToken(utils.RandomUserName(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = builder.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)

}
