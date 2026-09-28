package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	//1 hashpasswordครั้งที่1
	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	//2 hash passwordเดิมครั้งที่2(ต้องได้hashต่างกันจากsalt)
	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)

	//3 check correct password
	err = CheckPassword(password, hashedPassword1)
	require.NoError(t, err)

	//4 check รหัสผ่านที่ผิด (ต้องส่งErrMismatchedHashAndPassword)
	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
