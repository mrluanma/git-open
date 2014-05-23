package main

import (
	"flag"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"os"
)

func print_usage(cmd string) {
	fmt.Printf("Usage: %s <git repo dir>\n", cmd)
}

func main() {
	flag.Parse()

	remotes, err := DetectRemote(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	for _, r := range remotes {
		url, err := MangleURL(r.Url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "remote:%s, %s\n", r.Name, err.Error())
			continue
		}

		open.Run(url)
		break
	}
}
