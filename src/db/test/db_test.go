package test

import (
	"github.com/stretchr/testify/assert"
	"in-memory-challenge/src/db"
	"in-memory-challenge/src/utils"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func cleanStorage() {
	db.Storage = map[string]string{}
}

func TestCommands(t *testing.T) {
	var cmd db.Executable
	var key, value string
	var ok bool

	// SET
	cleanStorage()
	key = "a"
	value = "10"

	cmd = &db.Set{Key: key, Value: value}
	cmd.Execute()

	assert.Equal(t, db.Storage[key], value)

	cmd.Undo()
	_, ok = db.Storage[key]
	assert.False(t, ok)

	// GET
	cleanStorage()
	db.Storage[key] = value

	cmd = &db.Get{Key: key}
	assert.Equal(t, cmd.Execute(), value)

	cmd.Undo()

	cleanStorage()
	key = "b"
	cmd = &db.Get{Key: key}
	assert.Equal(t, cmd.Execute(), "NULL")

	// UNSET
	cleanStorage()
	db.Storage[key] = value

	cmd = &db.Unset{Key: key}
	cmd.Execute()
	_, ok = db.Storage[key]
	assert.False(t, ok)

	cmd.Undo()
	assert.Equal(t, db.Storage[key], value)

	// NUMEQUALTO
	cleanStorage()

	count := 10
	for i := 0; i < count; i++ {
		key = "c"
		db.Storage[key] = value
	}

	// NOT WORKING
	//cmd = &db.NumEqualTo{Value: value}
	//assert.Equal(t, cmd.Execute(), strconv.Itoa(count))
	//
	//cmd.Undo()
}

func TestTransactions(t *testing.T) {
	var cmd db.Executable
	var key, value1, value2 string

	key, value1, value2 = "a", "10", "20"

	// No transactions
	cleanStorage()
	_ = db.Transactions.Commit()
	//Assert("No transactions", transactions.Rollback(), ErrNoTransaction)
	//Assert("No transactions", transactions.Commit(), ErrNoTransaction)

	// Single transaction
	cleanStorage()
	cmd = &db.Set{Key: key, Value: value1}
	db.Transactions.Begin()
	db.Transactions.AddCommand(cmd)
	cmd.Execute()

	assert.Equal(t, db.Storage[key], value1)

	_ = db.Transactions.Rollback()

	_, ok := db.Storage[key]
	assert.False(t, ok)

	// Nested transactions
	cleanStorage()
	cmd = &db.Set{Key: key, Value: value1}
	db.Transactions.Begin()
	db.Transactions.AddCommand(cmd)
	cmd.Execute()

	assert.Equal(t, db.Storage[key], value1)

	cmd = &db.Set{Key: key, Value: value2}
	db.Transactions.Begin()
	db.Transactions.AddCommand(cmd)
	cmd.Execute()

	assert.Equal(t, db.Storage[key], value2)

	_ = db.Transactions.Rollback()

	assert.Equal(t, db.Storage[key], value1)

	_ = db.Transactions.Rollback()

	_, ok = db.Storage[key]
	assert.Equal(t, ok, false)
}

func TestRun(t *testing.T) {
	var cmd, key, value, expected string

	cleanStorage()
	key = "a"
	value = "10"

	// OK
	_, _ = db.Run("SET " + key + " " + value)
	_, _ = db.Run("GET " + key)
	_, _ = db.Run("UNSET " + key)
	_, _ = db.Run("NUMEQUALTO " + value)
	_, _ = db.Run("BEGIN")
	_, _ = db.Run("COMMIT")
	_, _ = db.Run("ROLLBACK")

	// NOK
	cmd = "UNKNOWNCOMMAND"
	result, _ := db.Run(cmd)
	expected = "ErrUnknownCommand (" + cmd + "): " + utils.ErrUnknownCommand.Error()
	assert.Equal(t, result, expected)

	cmd = "SET A"
	result, _ = db.Run(cmd)
	expected = "ErrNotEnoughArguments (SET): " + utils.ErrNotEnoughArguments.Error()
	assert.Equal(t, result, expected)

	cmd = "GET"
	result, _ = db.Run(cmd)
	expected = "ErrNotEnoughArguments (" + cmd + "): " + utils.ErrNotEnoughArguments.Error()
	assert.Equal(t, result, expected)

	cmd = "UNSET"
	result, _ = db.Run(cmd)
	expected = "ErrNotEnoughArguments (" + cmd + "): " + utils.ErrNotEnoughArguments.Error()
	assert.Equal(t, result, expected)

	cmd = "NUMEQUALTO"
	result, _ = db.Run(cmd)
	expected = "ErrNotEnoughArguments (" + cmd + "): " + utils.ErrNotEnoughArguments.Error()
	assert.Equal(t, result, expected)
}
