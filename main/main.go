package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dan-frohlich/dice"
)

func main() {

	r := dice.NewRoller()

	if len(os.Args) > 1 {
		for i, expr := range os.Args {
			switch (i) {
			case 0:
				continue
			default:
				result, _, err := r.Roll(expr)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				fmt.Print(result, " ")
			}
		}
		fmt.Println()
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Dice Roller Shell")
	fmt.Println("---------------------")

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
