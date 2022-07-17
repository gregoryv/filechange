package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gregoryv/filechange"
)

func main() {
	root, _ := os.Getwd()
	s := filechange.NewSensor(root, func(modified ...string) {
		for _, f := range modified {
			fmt.Println(f)
		}
	})
	s.Recursive = true
	s.Run(context.Background())
}
