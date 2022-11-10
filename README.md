Filechange provides a sensor of file modifications

## Quick start

    $ go install github.com/gregoryv/filechange/cmd/sense@latest
    $ sense -h
    Usage: sense [OPTIONS] [DIR]
    
    Options
        -i, --interval : 2s
        -r, --recursive
        -s, --script : "./.onchange.sh"
        -w, --write-example-script
        -h, --help

    
