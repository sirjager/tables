package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {

	password := RandomString(10)
	hashed, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)
	err = VerifyPassword(password, hashed)
	require.NoError(t, err)

	wrongPassword := RandomString(10)
	hashedWrongPass, err := HashPassword(wrongPassword)
	require.NoError(t, err)
	require.NotEmpty(t, hashedWrongPass)
	err = VerifyPassword(password, hashedWrongPass)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

}
