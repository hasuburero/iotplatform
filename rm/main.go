package main

import (
	"fmt"
	"rm/data"
	"sync"
)
import (
	"github.com/hasuburero/util/panic"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	panic.Start()
	wg.Wait()
}
