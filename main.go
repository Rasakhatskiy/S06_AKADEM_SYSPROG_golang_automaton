package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readAutomatonFromFile(path string) (automaton, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return automaton{}, err
	}
	r := bufio.NewReader(strings.NewReader(string(bytes)))
	return readAutomaton(r)
}

func main() {
	args := os.Args[1:]
	if len(args) != 3 {
		fmt.Println("Usage: config_file_path word_1 word_2")
		fmt.Println("Example: config.txt babbababab bbbaabbba")
	}

	a, err := readAutomatonFromFile(args[0])
	fmt.Println(a)
	fmt.Println(err)

}
