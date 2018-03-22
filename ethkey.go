package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/console"
	"gopkg.in/urfave/cli.v1"
)

// getPassPhrase retrieves the password associated with an account, either fetched
// from a list of preloaded passphrases, or requested interactively from the user.
func getPassPhrase(prompt string, confirmation bool, password string) string {
	// If a list of passwords was supplied, retrieve from them
	if password != "" {
		return password
	}
	// Otherwise prompt the user for the password
	if prompt != "" {
		fmt.Println(prompt)
	}
	password, err := console.Stdin.PromptPassword("Passphrase: ")
	if err != nil {
		utils.Fatalf("Failed to read passphrase: %v", err)
	}
	if confirmation {
		confirm, err := console.Stdin.PromptPassword("Repeat passphrase: ")
		if err != nil {
			utils.Fatalf("Failed to read passphrase confirmation: %v", err)
		}
		if password != confirm {
			utils.Fatalf("Passphrases do not match")
		}
	}
	return password
}

func main() {
	app := cli.NewApp()
	app.Name = "ethkey"
	app.Usage = "create Ethereum accounts as encrypted JSON keyfiles"
	app.Version = "0.1"
	app.Commands = []cli.Command{
		cli.Command{
			Name:    "new",
			Aliases: []string{"n"},
			Usage:   "new account",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "key-store",
					Usage:  "path to key store directory",
					EnvVar: "ETH_KEYSTORE",
				},
				cli.StringFlag{
					Name:   "passphrase-file",
					Usage:  "path to file containing account passphrase",
					EnvVar: "ETH_PASSWORD",
				},
			},
			Action: func(c *cli.Context) error {
				dir := c.String("key-store")

				if dir == "" {
					wd, err := os.Getwd()
					if err != nil {
						return cli.NewExitError("ethkey: failed to get working directory", 1)
					}
					dir = wd
				}

				passphrase := ""

				if c.String("passphrase-file") != "" {
					passphraseFile, err := ioutil.ReadFile(c.String("passphrase-file"))
					if err != nil {
						return cli.NewExitError("ethsign: failed to read passphrase file", 1)
					}
					passphrase = strings.TrimSuffix(string(passphraseFile), "\n")
				}

				password := getPassPhrase("Your new account is locked with a password. Please give a password. Do not forget this password.", true, passphrase)

				ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)

				acc, err := ks.NewAccount(password)

				if err != nil {
					return cli.NewExitError("ethkey: failed to create account", 1)
				}

				fmt.Println(acc.Address.String())

				return nil
			},
		},
	}

	app.Run(os.Args)
}
