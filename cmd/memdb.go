package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/flinnb/memdb/memdb"
)

func main() {
	fmt.Println("Welcome to MemDB!")

	memdb.Init()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSpace(input)
		parsedInput := parseInput(input)
		if len(parsedInput) == 0 {
			// If the user presses enter, just print an empty line and move on
			// It shouldn't actually be possible to get here, but it's included
			// just in case.
			fmt.Println("")
		} else {
			cmd := parsedInput[0]
			cmd = strings.ToUpper(cmd)

			switch cmd {
			case "":
				// If the user presses enter, just print an empty line and move on
				fmt.Println("")
			case "SET":
				if len(parsedInput) != 3 {
					fmt.Println("ERROR: You must provide a key *and* a value")
				} else {
					memdb.Set(parsedInput[1], parsedInput[2])
				}
			case "GET":
				if len(parsedInput) != 2 {
					fmt.Println("ERROR: You must provide a key")
				} else {
					fmt.Println(memdb.Get(parsedInput[1]))
				}
			case "DELETE":
				if len(parsedInput) != 2 {
					fmt.Println("ERROR: You must provide a key")
				} else {
					memdb.Delete(parsedInput[1])
				}
			case "COUNT":
				if len(parsedInput) != 2 {
					fmt.Println("ERROR: You must provide a value")
				} else {
					fmt.Println(memdb.Count(parsedInput[1]))
				}
			case "END":
				fmt.Println("Quitting...")
				return
			case "BEGIN":
				memdb.Begin()
			case "ROLLBACK":
				err := memdb.Rollback()
				if err != nil {
					fmt.Println(err.Error())
				}
			case "COMMIT":
				memdb.Commit()
			case "HELP":
				fmt.Println(helpText)
			default:
				fmt.Println("ERROR: invalid command.  Type `HELP` to see the valid list of commands")
			}
		}
	}

}

const helpText = `
SET [name] [value]
  Sets the name in the database to the given value
GET [name]
  Prints the value for the given name. If the value is not in the database,
  prints NULL
DELETE [name]
  Deletes the value from the database
COUNT [value]
  Returns the number of names that have the given value assigned to them. If
  that value is not assigned anywhere, prints 0
END
  Exits the database
BEGIN
  Begins a new transaction
ROLLBACK
  Rolls back the most recent transaction. If there is no transaction to
  rollback, prints TRANSACTION NOT FOUND error
COMMIT
  Commits *all* of the open transactions`

var wsRegex = regexp.MustCompile(`\s+`)

func parseInput(cmd string) []string {
	return wsRegex.Split(cmd, -1)
}
