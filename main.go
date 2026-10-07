package main

import (
	"fmt"

	"github.com/JacksonS25/gator/internal"
)

func main() {
	cfg, err := internal.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
	}

	internal.SetUser(cfg, "Jackson")

	cfg, err = internal.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
	}

	fmt.Println("Config:", cfg)
}
