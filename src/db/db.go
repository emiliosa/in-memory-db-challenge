package db

import (
	"fmt"
	"in-memory-challenge/src/utils"
	"os"
	"strconv"
	"strings"
)

const Commands = `Supported commands:
SET <name> <value> – Set the variable name to the value <value>. Neither variable names nor values will contain spaces.
(Name and value are case-sensitive strings with no spaces)

GET <name> – Print out the value of the variable name, or NULL if that variable is not set.
UNSET <name> – Unset the variable name, making it just like that variable was never set.
NUMEQUALTO <value> – Print out the number of variables that are currently set to value. If no variables equal that value, print 0.
END – Exit the program. Your program will always receive this as its last command. Commands will be fed to your program one at a time, with each command on its own line. Any output that your program generates should end with a newline character.

Supported transactions:
BEGIN – Open a new transaction block. Transaction blocks can be nested; a BEGIN can be issued inside of an existing block.
ROLLBACK – Undo all of the commands issued in the most recent transaction block, and close the block. Print nothing if successful, or print NO TRANSACTION if no transaction is in progress.
COMMIT – Close all open transaction blocks, permanently applying the changes made in them. Print nothing if successful, or print NO TRANSACTION if no transaction is in progress.`

var Storage = map[string]string{}
var Transactions = NestedTransactions{}

type Executable interface {
	Execute() (result string)
	Undo()
}

// Set "SET" command structure
type Set struct {
	Key   string
	Value string

	oldValue string
	exists   bool // Store true if Key already exists is DB
}

// Get "GET" command structure
type Get struct {
	Key string
}

// Unset "UNSET" command structure
type Unset struct {
	Key string

	oldValue string
	exists   bool
}

// NumEqualTo "NUMEQUALTO" command structure
type NumEqualTo struct {
	Value string
}

// Transaction group of commands
type Transaction struct {
	commands []Executable
}

// NestedTransactions group of transactions
type NestedTransactions struct {
	transactions []*Transaction
}

// AddCommand Add command to last opened transaction
func (t *NestedTransactions) AddCommand(cmd Executable) {
	tr, err := t.getLast()
	if err != nil {
		return
	}
	tr.AddCommand(cmd)
}

// Begin Open new transaction
func (t *NestedTransactions) Begin() *Transaction {
	tr := &Transaction{}
	t.transactions = append(t.transactions, tr)
	return tr
}

// Return last opened transaction and delete it from the list
func (t *NestedTransactions) popLast() (*Transaction, error) {
	if len(t.transactions) == 0 {
		return nil, utils.ErrNoTransaction
	}

	lastTrIndex := len(t.transactions) - 1
	tr, transactions := t.transactions[lastTrIndex], t.transactions[:lastTrIndex]
	t.transactions = transactions
	return tr, nil
}

// Commit all opened transactions
func (t *NestedTransactions) Commit() error {
	if len(t.transactions) == 0 {
		return utils.ErrNoTransaction
	}
	t.transactions = []*Transaction{}
	return nil
}

// Return last opened transaction
func (t *NestedTransactions) getLast() (*Transaction, error) {
	if len(t.transactions) == 0 {
		return nil, utils.ErrNoTransaction
	}

	lastTrIndex := len(t.transactions) - 1
	tr := t.transactions[lastTrIndex]
	return tr, nil
}

// Rollback last opened transaction
func (t *NestedTransactions) Rollback() error {
	tr, err := t.popLast()
	if err != nil {
		return err
	}
	tr.Rollback()
	return nil
}

// Rollback Undo all commands in transaction
func (tr *Transaction) Rollback() {
	for i := len(tr.commands) - 1; i >= 0; i-- {
		cmd := tr.commands[i]
		cmd.Undo()
	}
}

// AddCommand Add command in transaction
func (tr *Transaction) AddCommand(cmd Executable) {
	tr.commands = append(tr.commands, cmd)
}

// Execute Set value to the key in storage
func (cmd *Set) Execute() (result string) {
	oldValue, ok := Storage[cmd.Key]
	if ok {
		cmd.oldValue = oldValue
		cmd.exists = true
	}
	Storage[cmd.Key] = cmd.Value
	return
}

// Undo Back key in previous state
func (cmd *Set) Undo() {
	if cmd.exists {
		Storage[cmd.Key] = cmd.oldValue
	} else {
		delete(Storage, cmd.Key)
	}
}

// Execute Delete key from storage
func (cmd *Unset) Execute() (result string) {
	oldValue, ok := Storage[cmd.Key]
	if ok {
		cmd.exists = true
		cmd.oldValue = oldValue
		delete(Storage, cmd.Key)
	}
	return
}

// Undo deleting key from storage
func (cmd *Unset) Undo() {
	if cmd.exists {
		Storage[cmd.Key] = cmd.oldValue
	}
}

// Execute Get value by key from storage
func (cmd *Get) Execute() (result string) {
	if value, exists := Storage[cmd.Key]; exists {
		return value
	} else {
		return "NULL"
	}
}

// Undo Do nothing, just for implementing Executable interface
func (cmd *Get) Undo() {}

func (cmd *NumEqualTo) Execute() (result string) {
	count := 0
	for _, value := range Storage {
		if value == cmd.Value {
			count++
		}
	}
	return strconv.Itoa(count)
}

// Undo Do nothing, just for implementing Executable interface
func (cmd *NumEqualTo) Undo() {}

func Run(line string) (string, error) {
	args := strings.Split(line, " ")
	cmd := strings.ToUpper(args[0])
	params := args[1:]

	switch cmd {
	case "END":
		os.Exit(0)
	case "?":
	case "HELP":
		fmt.Println(Commands)
	case "GET":
		if len(params) < 1 {
			return utils.NotEnoughArguments(cmd), nil
		}

		cmd := &Get{Key: params[0]}
		fmt.Println(cmd.Execute())
	case "SET":
		if len(params) < 2 {
			return utils.NotEnoughArguments(cmd), nil
		}
		cmd := &Set{Key: params[0], Value: params[1]}
		Transactions.AddCommand(cmd)
		cmd.Execute()
	case "NUMEQUALTO":
		if len(params) < 1 {
			return utils.NotEnoughArguments(cmd), nil
		}

		cmd := &NumEqualTo{Value: params[0]}
		fmt.Println(cmd.Execute())
	case "UNSET":
		if len(params) < 1 {
			return utils.NotEnoughArguments(cmd), nil
		}

		cmd := &Unset{Key: params[0]}
		Transactions.AddCommand(cmd)
		cmd.Execute()
	case "BEGIN":
		Transactions.Begin()
	case "ROLLBACK":
		err := Transactions.Rollback()
		if err != nil {
			fmt.Println(err)
		}
	case "COMMIT":
		err := Transactions.Commit()
		if err != nil {
			fmt.Println(err)
		}
	default:
		return utils.UnknownCommand(cmd), nil
	}

	return "", nil
}
