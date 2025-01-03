package main

import (
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/patrykferenc/eecoin/internal/transaction/domain/transaction"
	"github.com/patrykferenc/eecoin/internal/transaction/net/http"
	"github.com/patrykferenc/eecoin/internal/transaction/query"
	"github.com/patrykferenc/eecoin/internal/wallet/domain/wallet"
	"github.com/urfave/cli/v2"
)

const remote = "http://localhost:8080"

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

				wl, err := wallet.ReadWalletFromDirectoryEcdsa(configPath, &passphrase)
				if err != nil {
					slog.Error("Cannot read wallet")
					os.Exit(1)
				}
				if wl.MainId == nil {
					mainId, err := wallet.NewEcdsaKey()
					if err != nil {
						slog.Error("Cannot generate main identity")
						os.Exit(1)
					}

					err = wl.SetMainIdentity(&mainId)
					if err != nil {
						slog.Error("Cannot set main identity")
						os.Exit(1)
					}
				}
				return wl.ExportWalletEcdsa(configPath, &passphrase)
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

				wl, err := wallet.ReadWalletFromDirectoryEcdsa(configPath, &passphrase)
				if err != nil {
					slog.Error("Cannot read wallet")
					os.Exit(1)
				}
				newKey, err := wallet.NewEcdsaKey()
				if err != nil {
					slog.Error("Cannot generate key")
					os.Exit(1)
				}
				err = wl.Add(newKey)
				if err != nil {
					slog.Error("Cannot add key")
					os.Exit(1)
				}

				err = wl.ExportWalletEcdsa(configPath, &passphrase)
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
				wl, err := wallet.ReadWalletFromDirectoryEcdsa(configPath, &passphrase)
				if err != nil {
					slog.Error("Cannot read wallet")
					os.Exit(1)
				}
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
		{
			Name:  "balance",
			Usage: "list balance",
			Action: func(c *cli.Context) error {
				if c.Args().Len() != 2 {
					slog.Error("Two arguments needed : <config file path> <passphrase>")
					os.Exit(1)
				}
				configPath := c.Args().Get(0)
				passphrase := c.Args().Get(1)

				unspentRepo := http.NewUnspentOutputsRepository(remote)
				balance := query.NewGetBalance(unspentRepo)
				wl, _ := wallet.ReadWalletFromDirectoryEcdsa(configPath, &passphrase)
				selfPub := wl.MainId.Public
				marshalled, err := x509.MarshalPKIXPublicKey(selfPub)
				if err != nil {
					fmt.Printf("Cannot marshal public key: %s\n", err)
				}
				selfAddr := hex.EncodeToString(marshalled)
				b, err := balance.GetBalance(query.GetBalanceRequest{Address: selfAddr})
				if err != nil {
					fmt.Printf("Cannot get balance: %s\n", err)
				}
				fmt.Printf("Balance: %d\n", b)
				return nil
			},
		},
		{
			Name:  "transfer",
			Usage: "transfer coins",
			Action: func(c *cli.Context) error {
				if c.Args().Len() != 2 {
					slog.Error("Four arguments needed : <index from list beginning with 0> <integer amount> <config file path> <passphrase>")
					os.Exit(1)
				}
				index, err := strconv.Atoi(c.Args().Get(0))
				if err != nil {
					slog.Error("Cannot parse index")
					os.Exit(1)
				}
				amount, err := strconv.Atoi(c.Args().Get(1))
				if err != nil {
					slog.Error("Cannot parse amount")
					os.Exit(1)
				}
				configPath := c.Args().Get(2)
				passphrase := c.Args().Get(3)
				wl, err := wallet.ReadWalletFromDirectoryEcdsa(configPath, &passphrase)
				if err != nil {
					slog.Error("Cannot read wallet")
					os.Exit(1)
				}

				if !wl.Keys[index].Present {
					slog.Error("Key not present")
					os.Exit(1)
				}
				reciverPub := wl.Keys[index].Key.Public
				recMarshalled, err := x509.MarshalPKIXPublicKey(reciverPub)
				if err != nil {
					fmt.Printf("Cannot marshal public key: %s\n", err)
				}
				recieverAddr := hex.EncodeToString(recMarshalled)

				selfPub := wl.MainId.Public
				marshalled, err := x509.MarshalPKIXPublicKey(selfPub)
				if err != nil {
					fmt.Printf("Cannot marshal public key: %s\n", err)
				}
				selfAddr := hex.EncodeToString(marshalled)

				unspentRepo := http.NewUnspentOutputsRepository(remote)
				tr, err := transaction.New(recieverAddr, selfAddr, amount, wl.MainId.Private(), unspentRepo)
				if err != nil {
					return err
				}
				fmt.Printf("Transaction created: %s\n", tr.ID())

				return nil
			},
		},
	}
}
