package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/gregoryv/cmdline"
	"github.com/gregoryv/filechange"
)

func main() {
	var (
		cli       = cmdline.NewBasicParser()
		root, _   = os.Getwd()
		recursive = cli.Flag("-r, --recursive")
		script    = cli.Option("-s, --script").String("./.onchange.sh")
		dir       = cli.NamedArg("DIR").String(root)
	)
	cli.Parse()

	s := filechange.NewSensor(dir, func(modified ...string) {
		out, err := exec.Command(script, modified...).CombinedOutput()
		if err != nil {
			fmt.Println(err)
		}
		os.Stdout.Write(out)
	})
	s.Recursive = recursive
	s.Run(context.Background())
}
