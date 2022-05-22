package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime/debug"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Kichiyaki/gootp/internal"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

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

	buildInfo, _ := debug.ReadBuildInfo()

	return &cli.App{
		Name:    "gootp",
		Usage:   "Two-Factor Authentication (2FA) App compatible with andOTP file format",
		Version: buildInfo.Main.Version,
		Action: func(c *cli.Context) error {
			password, err := getPassword(c)
			if err != nil {
				return err
			}

			b, err := os.ReadFile(c.String("path"))
			if err != nil {
				return fmt.Errorf("something went wrong while reading file: %w", err)
			}

			entries, err := internal.DecryptAsEntries(b, password)
			if err != nil {
				return fmt.Errorf("something went wrong while decrypting file: %w", err)
			}

			p := tea.NewProgram(internal.NewUI(entries), tea.WithAltScreen())
			if err := p.Start(); err != nil {
				return fmt.Errorf("p.Start: %w", err)
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Aliases:     []string{"p"},
				Usage:       "path to andOTP backup file",
				Required:    false,
				DefaultText: "$HOME/.otp_accounts.json",
				Value:       path.Join(dirname, ".otp_accounts.json"),
			},
			&cli.StringFlag{
				Name:     "password",
				Usage:    "encryption password",
				Required: false,
			},
		},
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			newEncryptCommand(),
			newDecryptCommand(),
		},
	}, nil
}

func newDecryptCommand() *cli.Command {
	return &cli.Command{
		Name:   "decrypt",
		Usage:  "Decrypts the specified file",
		Action: newEncryptDecryptActionFunc(internal.Decrypt),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "output",
				Usage:    "Write to file instead of stdout",
				Aliases:  []string{"o"},
				Required: false,
			},
		},
	}
}

func newEncryptCommand() *cli.Command {
	return &cli.Command{
		Name:   "encrypt",
		Usage:  "Encrypts the specified file",
		Action: newEncryptDecryptActionFunc(internal.Encrypt),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "output",
				Usage:    "Write to file instead of stdout",
				Aliases:  []string{"o"},
				Required: false,
			},
		},
	}
}

func newEncryptDecryptActionFunc(fn func(text, password []byte) ([]byte, error)) cli.ActionFunc {
	return func(c *cli.Context) error {
		password, err := getPassword(c)
		if err != nil {
			return err
		}

		b, err := os.ReadFile(c.String("path"))
		if err != nil {
			return fmt.Errorf("something went wrong while reading file: %w", err)
		}

		result, err := fn(b, password)
		if err != nil {
			return fmt.Errorf("something went wrong while processing file: %w", err)
		}

		output := c.String("output")
		if output != "" {
			if err := os.WriteFile(output, result, 0600); err != nil {
				return fmt.Errorf("something went wrong while saving result: %w", err)
			}
		} else {
			fmt.Print(string(result))
		}

		return nil
	}
}

func getPassword(c *cli.Context) ([]byte, error) {
	password := []byte(c.String("password"))
	if len(password) == 0 {
		return readPasswordFromStdin()
	}
	return password, nil
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
