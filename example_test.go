package filechange_test

import (
	"context"
	"fmt"
	"time"

	"github.com/gregoryv/filechange"
)

func Example() {
	s := &filechange.Sensor{
		Visit: func(modified ...string) {
			fmt.Println(modified)
		},
		Recursive: true,
		Pause:     2 * time.Second,
	}
	ctx, _ := context.WithCancel(context.Background())
	go s.Run(ctx)
}
