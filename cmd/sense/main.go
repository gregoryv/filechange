package main

import (
	"context"
	"fmt"

	"github.com/gregoryv/filechange"
)

func main() {
	s := new(filechange.Sensor)
	s.Visit = func(modified ...string) {
		for _, f := range modified {
			fmt.Println(f)
		}
	}
	s.Recursive = true
	s.Run(context.Background())
}
