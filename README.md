# In-Memory Database Challenge

**[PDF challenge](https://github.com/emiliosa/in-memory-db-challenge/blob/main/CodingChallenge.pdf)**.

Im-memory database implementation made with Golang.

## Requirements:

- Go version 1.16+

## Run:

- To run, type `go run main.go`

## Commands:

`SET <name> <value>`
- Set the variable **name** to the value **value**. Neither variable names nor values will contain spaces.

`GET <name>`
- Print out the value of the variable name, or NULL if that variable is not set.

`UNSET <name>`
- Unset the variable name, making it just like that variable was never set.

`NUMEQUALTO <value>`
- Print out the number of variables that are currently set to value. If no variables equal that value, print 0.

`END`
- Exit the program. Your program will always receive this as its last command. Commands will be fed to your program one at a time, with each command on its own line. Any output that your program generates should end with a newline character.

`HELP`
- Prints command list.

## Supported transactions:

`BEGIN`
- Open a new transaction block. Transaction blocks can be nested; a BEGIN can be issued inside of an existing block.

`ROLLBACK`
- Undo all of the commands issued in the most recent transaction block, and close the block. Print nothing if successful, or print NO TRANSACTION if no transaction is in progress.

`COMMIT`
- Close all open transaction blocks, permanently applying the changes made in them. Print nothing if successful, or print NO TRANSACTION if no transaction is in progress.

## Run tests:

- To run tests, type `go test ./src/db/... -coverpkg=./... -v`

### TODO:

- Add Dockerfile
- Refactor db.go
- Complete tests coverage