package main

import (
	"fmt"
	"os"

	"github.com/torrefatto/scredenziato/helpers"
)

func main() {
	h, err := helpers.NewFileBasedStore()
	if err != nil {
		fmt.Print(err)
		os.Exit(2)
	}

	list, err := h.List()
	if err != nil {
		fmt.Print(err)
		os.Exit(3)
	}

	fmt.Print(list)
}
