package main

import (
	"log"
	"os"

	"github.com/Kichiyaki/gootp/cmd/gootp/decrypt"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                 "gootp",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			decrypt.NewCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln("app.Run:", err)
	}
}
