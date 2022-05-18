package main

import (
	"fmt"
	"github.com/Kichiyaki/gootp/internal"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
	"log"
	"os"
	"path"
	"syscall"
)

var Version = "development"

func main() {
	app, err := newApp()
	if err != nil {
		log.Fatalln("newApp:", err)
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln("app.Run:", err)
	}
}

func newApp() (*cli.App, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("couldn't get user home dir: %w", err)
	}

	return &cli.App{
		Name:    "gootp",
		Version: Version,
		Action: func(c *cli.Context) error {
			var err error
			password := []byte(c.String("password"))
			if len(password) == 0 {
				password, err = readPasswordFromStdin()
			}
			if err != nil {
				return err
			}

			b, err := os.ReadFile(c.String("path"))
			if err != nil {
				return fmt.Errorf("something went wrong while reading file: %w", err)
			}

			entries, err := internal.DecryptAsEntries(password, b)
			if err != nil {
				return fmt.Errorf("something went wrong while decrypting file: %w", err)
			}

			for _, entry := range entries {
				otp, err := internal.GenerateOTP(entry)
				if err != nil {
					log.Printf("%s - %s: %s", entry.Issuer, entry.Label, err)
					continue
				}

				log.Printf("%s - %s: %s", entry.Issuer, entry.Label, otp)
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Usage:       "path to encrypted andotp file",
				Required:    false,
				DefaultText: "$HOME/.otp_accounts.json.aes",
				Value:       path.Join(dirname, ".otp_accounts.json.aes"),
			},
			&cli.StringFlag{
				Name:     "password",
				Usage:    "encryption password",
				Required: false,
			},
		},
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			newDecryptCommand(),
		},
	}, nil
}

func newDecryptCommand() *cli.Command {
	return &cli.Command{
		Name:  "decrypt",
		Usage: "decrypt backup file generated by andotp",
		Action: func(c *cli.Context) error {
			var err error
			password := []byte(c.String("password"))
			if len(password) == 0 {
				password, err = readPasswordFromStdin()
			}
			if err != nil {
				return err
			}

			b, err := os.ReadFile(c.String("path"))
			if err != nil {
				return fmt.Errorf("something went wrong while reading file: %w", err)
			}

			result, err := internal.Decrypt(password, b)
			if err != nil {
				return fmt.Errorf("something went wrong while decrypting file: %w", err)
			}

			fmt.Print(string(result))

			return nil
		},
	}
}

func readPasswordFromStdin() ([]byte, error) {
	fmt.Print("Password: ")

	pass, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return nil, fmt.Errorf("term.ReadPassword: %w", err)
	}

	fmt.Print("\n")

	return pass, nil
}
