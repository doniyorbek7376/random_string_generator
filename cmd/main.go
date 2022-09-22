package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/doniyorbek7376/random_string_generator/app"
)

func main() {
	var regex string
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter regex for the generator: ")
	regex, _ = reader.ReadString('\n')
	regex = strings.Trim(regex, "\n")

	generator := app.NewGenerator()
	values, err := generator.Generate(regex, 10)
	if err != nil {
		panic(err)
	}
	for _, value := range values {
		fmt.Println(value)
	}
}
