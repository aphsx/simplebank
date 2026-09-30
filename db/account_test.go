package db

import (
	"context"
	"database/sql"
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

//createRandomAccount helper funct ช้วยสร้างaccountสุ่มขึ้นมา1บัญชี

func createRandomAccount(t *testing.T) Account {

	user := createRandomUser(t) // 1. สร้าง User สุ่มขึ้นมาก่อนเสมอ

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)      //ต้องไม่มีerr
	require.NotEmpty(t, account) //ข้อมูลต้องไม่ว่าง

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID) //ID ต้องถูกรันจากdb
	require.NotZero(t, account.CreateAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.WithinDuration(t, acc1.CreateAt, acc2.CreateAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	acc1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	// เมื่อจบแล้วค้นหาอีกรอบ ต้องได้sql.ErrNoRowsกลับ
	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, acc2)

}

func TestListAccounts(t *testing.T) {
	user := createRandomUser(t)
	for _, currency := range []string{util.USD, util.EUR, util.THB} {
		_, err := testQueries.CreateAccount(context.Background(), CreateAccountParams{
			Owner:    user.Username,
			Balance:  util.RandomMoney(),
			Currency: currency,
		})
		require.NoError(t, err)
	}
	createRandomAccount(t)

	arg := ListAccountsParams{
		Owner:  user.Username,
		Limit:  2,
		Offset: 1,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 2)

	for _, account := range accounts {
		require.Equal(t, user.Username, account.Owner)
	}
}
