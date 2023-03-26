/*
This work is licensed and released under GNU GPL v3 or any other later versions.
The full text of the license is below/ found at <http://www.gnu.org/licenses/>

(c) 2023 Rocket Pool Pty Ltd. Modified under GNU GPL v3.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package wallet

import (
	"github.com/urfave/cli"

	cliutils "github.com/stader-labs/stader-node/shared/utils/cli"
)

// Register commands
func RegisterCommands(app *cli.App, name string, aliases []string) {
	app.Commands = append(app.Commands, cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage the node wallet",
		Subcommands: []cli.Command{

			{
				Name:      "status",
				Aliases:   []string{"s"},
				Usage:     "Get the node wallet status",
				UsageText: "stader-cli wallet status",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return getStatus(c)

				},
			},

			{
				Name:      "recover",
				Aliases:   []string{"r"},
				Usage:     "Recover a node wallet from a mnemonic phrase",
				UsageText: "stader-cli wallet recover [options]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "password, p",
						Usage: "The password to secure the wallet with (if not already set)",
					},
					cli.StringFlag{
						Name:  "mnemonic, m",
						Usage: "The mnemonic phrase to recover the wallet from",
					},
					cli.BoolFlag{
						Name:  "skip-validator-key-recovery, k",
						Usage: "Recover the node wallet, but do not regenerate its validator keys",
					},
					cli.StringFlag{
						Name:  "derivation-path, d",
						Usage: "Specify the derivation path for the wallet.\nOmit this flag (or leave it blank) for the default of \"m/44'/60'/0'/0/%d\" (where %d is the index).\nSet this to \"ledgerLive\" to use Ledger Live's path of \"m/44'/60'/%d/0/0\".\nSet this to \"mew\" to use MyEtherWallet's path of \"m/44'/60'/0'/%d\".\nFor custom paths, simply enter them here.",
					},
					cli.UintFlag{
						Name:  "wallet-index, i",
						Usage: "Specify the index to use with the derivation path when recovering your wallet",
						Value: 0,
					},
					cli.StringFlag{
						Name:  "address, a",
						Usage: "If you are recovering a wallet that was not generated by the Stadernode and don't know the derivation path or index of it, enter the address here. The Stadernode will search through its library of paths and indices to try to find it.",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Validate flags
					if c.String("password") != "" {
						if _, err := cliutils.ValidateNodePassword("password", c.String("password")); err != nil {
							return err
						}
					}
					if c.String("mnemonic") != "" {
						if _, err := cliutils.ValidateWalletMnemonic("mnemonic", c.String("mnemonic")); err != nil {
							return err
						}
					}

					// Run
					return recoverWallet(c)

				},
			},

			{
				Name:      "init",
				Aliases:   []string{"i"},
				Usage:     "Initialize the node wallet",
				UsageText: "stader-cli wallet init [options]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "password, p",
						Usage: "The password to secure the wallet with (if not already set)",
					},
					cli.BoolFlag{
						Name:  "confirm-mnemonic, c",
						Usage: "Automatically confirm the mnemonic phrase",
					},
					cli.StringFlag{
						Name:  "derivation-path, d",
						Usage: "Specify the derivation path for the wallet.\nOmit this flag (or leave it blank) for the default of \"m/44'/60'/0'/0/%d\" (where %d is the index).\nSet this to \"ledgerLive\" to use Ledger Live's path of \"m/44'/60'/%d/0/0\".\nSet this to \"mew\" to use MyEtherWallet's path of \"m/44'/60'/0'/%d\".\nFor custom paths, simply enter them here.",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Validate flags
					if c.String("password") != "" {
						if _, err := cliutils.ValidateNodePassword("password", c.String("password")); err != nil {
							return err
						}
					}

					// Run
					return initWallet(c)

				},
			},

			{
				Name:      "export",
				Aliases:   []string{"e"},
				Usage:     "Export the node wallet in JSON format",
				UsageText: "stader-cli wallet export",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return exportWallet(c)

				},
			},
		},
	})
}
