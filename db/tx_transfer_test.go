package db

//testTransferTx testความถูกต้องการโอนเงิน Concurrent Transactions
//ป้องกัน Race condition และเงินหายระหว่างทาง

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	//1 Create 2 account for test และยอดเงินตั้งต้น
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">> ยอดเงินหลังโอน:", account1.Balance, account2.Balance)

	//2 ปล่อยgoroutine โอนพร้อมกัน 5ครั้ง ครั้งละ 10บาท
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	//check result จากแต่ละ Goroutine
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err) //

		result := <-results
		require.NotEmpty(t, result)

		//3.1 check row data in Transfer table
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)

		//3.2 check Entries (เงินเข้า ออก)
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)

		//
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)

		//3.3 check Balance คงเหลือ แต่ละบัญชีแต่ละรอบ
		diff1 := account1.Balance - result.FromAccount.Balance
		diff2 := result.ToAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2) //ยอดเงินที่หักออก ต้องเท่ากับยอดที่เพิ่มเข้า
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) //ผลต่างต้องหารด้วยจำนวนโอนลงตัว

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k) //แต่ละรอบ kต้องไม่ซ้ำกัน (พิสูจไม่มีRace condition)
		existed[k] = true

	}

	//4 check ยอดเงินคงเหลือ สุดท้ายในฐานข้อมูลจริง
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> ยอดเงินหลังโอน:", updatedAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	//1 จำลองการโอนเงินสวนทางกันทั้งหมด10ธุรกรรม

	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		//สลับทิศทางครึ่งหนึ่ง :ให้2โอนคืนให้บัญชี1
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()

	}

	//2 ตรวจresult ทุกTransaction ต้องsuccesไร้ปัญหาDeadlock 100%
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err) //ต้องไม่มี pq: deadlock detect สักครั้งเดียว
	}

	//3 check ยอดเงินคงเหลือสุดท้าย โอนไป5ครั้งอนกลับ5ครั้ง ยอดต้องเท่าเดิม

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

}
