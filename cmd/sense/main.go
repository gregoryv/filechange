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
		interval  = cli.Option("-i, --interval").Duration("2s")
		recursive = cli.Flag("-r, --recursive")
		script    = cli.Option("-s, --script").String("./.onchange.sh")
		writeEx   = cli.Flag("-w, --write-example-script")
		dir       = cli.NamedArg("DIR").String(root)
	)
	cli.Parse()

	if writeEx {
		os.WriteFile(script, []byte(myscript), 0755)
		return
	}

	s := filechange.NewSensor(dir, func(modified ...string) {
		out, err := exec.Command(script, modified...).CombinedOutput()
		if err != nil {
			fmt.Println(err)
		}
		os.Stdout.Write(out)
	})
	s.Recursive = recursive
	s.Interval = interval
	s.Run(context.Background())
}

const myscript = `#!/bin/bash -e
path=$1
dir=$(dirname "$path")
filename=$(basename "$path")
extension="${filename##*.}"
nameonly="${filename%.*}"

case $extension in
    go)
        goimports -w $path
        ;;
esac

go test ./...
`
