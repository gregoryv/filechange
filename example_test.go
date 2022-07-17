package filechange_test

import (
	"context"
	"fmt"

	"github.com/gregoryv/filechange"
)

func Example() {
	s := filechange.NewSensor(".", func(modified ...string) {
		fmt.Println(modified)
	})
	s.Recursive = true

	ctx, _ := context.WithCancel(context.Background())
	go s.Run(ctx)
}
