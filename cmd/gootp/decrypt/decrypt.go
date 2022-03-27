package decrypt

import (
	"fmt"
	"os"
	"syscall"

	"github.com/Kichiyaki/gootp/internal/andotp"
	"golang.org/x/term"

	"github.com/urfave/cli/v2"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "decrypt",
		Usage: "decrypt backup file generated by andotp",
		Action: func(c *cli.Context) error {
			var err error
			password := []byte(c.String("password"))
			if len(password) == 0 {
				password, err = getPassword()
			}
			if err != nil {
				return err
			}

			b, err := os.ReadFile(c.String("path"))
			if err != nil {
				return fmt.Errorf("something went wrong while reading file: %w", err)
			}

			result, err := andotp.Decrypt(password, b)
			if err != nil {
				return fmt.Errorf("something went wrong while decrypting file: %w", err)
			}

			fmt.Print("\n")
			fmt.Print(string(result))

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Usage:    "path to backup file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "password",
				Usage:    "encryption password",
				Required: false,
			},
		},
	}
}

func getPassword() ([]byte, error) {
	fmt.Print("Password: ")

	pass, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return nil, fmt.Errorf("term.ReadPassword: %w", err)
	}

	return pass, nil
}
