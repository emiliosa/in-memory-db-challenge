package main

import (
	"bufio"
	"fmt"
	"in-memory-challenge/src/db"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf(">> ")

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := strings.TrimSpace(scanner.Text())
		output, err := db.Run(line)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(output)
	}
}
