package main

import (
	"bufio"
	"dice"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Dice Roller Shell")
	fmt.Println("---------------------")

	r := dice.NewRoller()
	var text string
	for {
		fmt.Print("-> ")
		text, _ = reader.ReadString('\n')
		if isExit(text) {
			break
		}
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		result, plan, err := r.Roll(text)
		if err != nil {
			fmt.Println("<-", "ERROR", err)
			continue
		}
		fmt.Println("<-", result, ":", plan)
	}

}

func isExit(input string) bool {
	return strings.HasPrefix(input, "exit") ||
		strings.HasPrefix(input, "quit") ||
		strings.HasPrefix(input, "q")
}
