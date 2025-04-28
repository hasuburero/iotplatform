package main

import (
	"fmt"
	"github.com/hasuburero/util/panic"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ./main Addr Port")
		fmt.Println("Ex: ./main localhost 8080")
		return
	}
	panic.Start()
}
