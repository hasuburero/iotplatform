package main

import (
	"fmt"
	"os"
	"rm/api"
	"sync"
)
import (
	"github.com/hasuburero/util/panic"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ./main address port")
		fmt.Println("Ex: ./main \"\" 8080")
		return
	}
	fmt.Printf("Addr: \"%s\"\n", os.Args[1])
	fmt.Printf("Port: \"%s\"\n", os.Args[2])
	var wg sync.WaitGroup
	wg.Add(1)
	panic.Start()
	api.Start("", "8080")
	wg.Wait()
}
