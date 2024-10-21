package main

import (
	"fmt"
	"github.com/patrykferenc/eecoin/internal/wallet/domain/wallet"
	"github.com/urfave/cli/v2"
	"log/slog"
	"os"
)

func main() {

	app := &cli.App{
		Name:     "wallet",
		Usage:    "a cryptocurrency wallet",
		Commands: setupCliCommands(),
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
func setupCliCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "init",
			Usage: "initialize a wallet",
			Action: func(c *cli.Context) error {
				if c.Args().Len() != 2 {
					slog.Error("Two arguments needed : <config file path> <passphrase>")
					os.Exit(1)
				}
				configPath := c.Args().Get(0)
				passphrase := c.Args().Get(1)

				wl, _ := wallet.ReadWalletFromDirectory(configPath, &passphrase)
				if wl.MainId == nil {
					mainId, _ := wallet.NewRsaKey()
					_ = wl.SetMainIdentity(&mainId)
				}
				return wl.ExportWallet(configPath, &passphrase)
			},
		},
		{
			Name:  "gen",
			Usage: "generate a key to file",
			Action: func(c *cli.Context) error {
				if c.Args().Len() != 2 {
					slog.Error("Two arguments needed : <config file path> <passphrase>")
					os.Exit(1)
				}
				configPath := c.Args().Get(0)
				passphrase := c.Args().Get(1)

				wl, _ := wallet.ReadWalletFromDirectory(configPath, &passphrase)
				newKey, _ := wallet.NewRsaKey()
				_ = wl.Add(newKey)
				err := wl.ExportWallet(configPath, &passphrase)
				return err
			},
		},
		{
			Name:  "list",
			Usage: "list keys",
			Action: func(c *cli.Context) error {
				if c.Args().Len() != 2 {
					slog.Error("Two arguments needed : <config file path> <passphrase>")
					os.Exit(1)
				}
				configPath := c.Args().Get(0)
				passphrase := c.Args().Get(1)
				wl, _ := wallet.ReadWalletFromDirectory(configPath, &passphrase)
				for k, v := range wl.Keys {
					fmt.Printf("pub: %s | private %t\n", k, v.Present)
				}
				if wl.MainId != nil {
					fmt.Println("Main id is set")
				} else {
					fmt.Println("Main id is not set")
				}
				return nil
			},
		},
	}
}
