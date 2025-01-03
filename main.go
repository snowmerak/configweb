package main

import (
	"github.com/snowmerak/configweb/internal/cli"
)

func main() {
	if err := cli.NewApp().Run(); err != nil {
		panic(err)
	}
}
