package main

import (
	"fmt"

	"github.com/doniyorbek7376/random_string_generator/app"
)

func main() {
	values, err := app.Generate("[a-f-]{5}", 10)
	if err != nil {
		panic(err)
	}
	for _, value := range values {
		fmt.Println(value)
	}
}
