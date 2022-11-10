package filechange_test

import (
	"context"
	"fmt"

	"github.com/gregoryv/filechange"
)

func Example() {
	s := filechange.NewSensor(".", func(modified ...string) {
		// do something with the modified files
		fmt.Println(modified)
	})
	s.Recursive = true

	go s.Run(context.Background())
}
