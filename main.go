package main

import (
	"io/ioutil"
	"log"
	"os"
)

func main() {
	arg := os.Args[1:]

	if len(arg) != 1 {
		log.Fatal("usage: cal <filename>")
	}

	f, err := ioutil.ReadFile(arg[0])
	if err != nil {
		panic(err)
	}
	run(string(f))
}

func run(s string) {
	scanner := NewScanner(s)
	for _, i := range scanner.tokens {
		println(i.lexeme)
	}
	parser := NewParser(scanner.tokens)
	parser.PrintAST()
}
