package node

import (
	"github.com/urfave/cli"

	cliutils "github.com/stader-labs/stader-node/shared/utils/cli"
)

// Register commands
func RegisterCommands(app *cli.App, name string, aliases []string) {
	app.Commands = append(app.Commands, cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Manage the node",
		Subcommands: []cli.Command{

			{
				Name:      "status",
				Aliases:   []string{"s"},
				Usage:     "Get the node's status",
				UsageText: "rocketpool node status",
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
				Name:      "sync",
				Aliases:   []string{"y"},
				Usage:     "Get the sync progress of the eth1 and eth2 clients",
				UsageText: "rocketpool node sync",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return getSyncProgress(c)

				},
			},

			{
				Name:      "register",
				Aliases:   []string{"r"},
				Usage:     "Register the node with Rocket Pool",
				UsageText: "rocketpool node register [options]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "timezone, t",
						Usage: "The timezone location to register the node with (in the format 'Country/City')",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Validate flags
					if c.String("timezone") != "" {
						if _, err := cliutils.ValidateTimezoneLocation("timezone location", c.String("timezone")); err != nil {
							return err
						}
					}

					// Run
					return registerNode(c)

				},
			},

			{
				Name:      "rewards",
				Aliases:   []string{"e"},
				Usage:     "Get the time and your expected RPL rewards of the next checkpoint",
				UsageText: "rocketpool node rewards",
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return getRewards(c)

				},
			},

			{
				Name:      "set-withdrawal-address",
				Aliases:   []string{"w"},
				Usage:     "Set the node's withdrawal address",
				UsageText: "rocketpool node set-withdrawal-address [options] address",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm setting withdrawal address",
					},
					cli.BoolFlag{
						Name:  "force",
						Usage: "Force update the withdrawal address, bypassing the 'pending' state that requires a confirmation transaction from the new address",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					withdrawalAddress := c.Args().Get(0)

					// Run
					return setWithdrawalAddress(c, withdrawalAddress)

				},
			},

			{
				Name:      "confirm-withdrawal-address",
				Aliases:   []string{"f"},
				Usage:     "Confirm the node's pending withdrawal address if it has been set back to the node's address itself",
				UsageText: "rocketpool node confirm-withdrawal-address [options]",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm withdrawal address",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return confirmWithdrawalAddress(c)

				},
			},

			{
				Name:      "set-timezone",
				Aliases:   []string{"t"},
				Usage:     "Set the node's timezone location",
				UsageText: "rocketpool node set-timezone [options]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "timezone, t",
						Usage: "The timezone location to set for the node (in the format 'Country/City')",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Validate flags
					if c.String("timezone") != "" {
						if _, err := cliutils.ValidateTimezoneLocation("timezone location", c.String("timezone")); err != nil {
							return err
						}
					}

					// Run
					return setTimezoneLocation(c)

				},
			},

			{
				Name:      "swap-rpl",
				Aliases:   []string{"p"},
				Usage:     "Swap old RPL for new RPL",
				UsageText: "rocketpool node swap-rpl [options]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "amount, a",
						Usage: "The amount of old RPL to swap (or 'all')",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Validate flags
					if c.String("amount") != "" && c.String("amount") != "all" {
						if _, err := cliutils.ValidatePositiveEthAmount("swap amount", c.String("amount")); err != nil {
							return err
						}
					}

					// Run
					return nodeSwapRpl(c)

				},
			},

			{
				Name:      "stake-rpl",
				Aliases:   []string{"k"},
				Usage:     "Stake RPL against the node",
				UsageText: "rocketpool node stake-rpl [options]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "amount, a",
						Usage: "The amount of RPL to stake (or 'min', 'max', or 'all')",
					},
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm RPL stake",
					},
					cli.BoolFlag{
						Name:  "swap, s",
						Usage: "Automatically confirm swapping old RPL before staking",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Validate flags
					if c.String("amount") != "" && c.String("amount") != "min" && c.String("amount") != "max" && c.String("amount") != "all" {
						if _, err := cliutils.ValidatePositiveEthAmount("stake amount", c.String("amount")); err != nil {
							return err
						}
					}

					// Run
					return nodeStakeRpl(c)

				},
			},

			{
				Name:      "claim-rewards",
				Aliases:   []string{"c"},
				Usage:     "Claim available RPL and ETH rewards for any checkpoint you haven't claimed yet",
				UsageText: "rocketpool node claim-rpl [options]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "restake-amount, a",
						Usage: "The amount of RPL to automatically restake during claiming (or '150%%' to stake up to 150%% collateral, or 'all' for all available RPL)",
					},
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm rewards claim",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return nodeClaimRewards(c)

				},
			},

			{
				Name:      "withdraw-rpl",
				Aliases:   []string{"i"},
				Usage:     "Withdraw RPL staked against the node",
				UsageText: "rocketpool node withdraw-rpl [options]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "amount, a",
						Usage: "The amount of RPL to withdraw (or 'max')",
					},
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm RPL withdrawal",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Validate flags
					if c.String("amount") != "" && c.String("amount") != "max" {
						if _, err := cliutils.ValidatePositiveEthAmount("withdrawal amount", c.String("amount")); err != nil {
							return err
						}
					}

					// Run
					return nodeWithdrawRpl(c)

				},
			},

			{
				Name:      "deposit",
				Aliases:   []string{"d"},
				Usage:     "Make a deposit and create a minipool",
				UsageText: "rocketpool node deposit [options]",
				Flags: []cli.Flag{
					/*cli.StringFlag{
						Name:  "amount, a",
						Usage: "The amount of ETH to deposit (0, 16 or 32)",
					},*/
					//cli.StringFlag{
					//	Name:  "max-slippage, s",
					//	Usage: "The maximum acceptable slippage in node commission rate for the deposit (or 'auto'). Only relevant when the commission rate is not fixed.",
					//},
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm deposit",
					},
					cli.StringFlag{
						Name:  "salt, l",
						Usage: "An optional seed to use when generating the new minipool's address. Use this if you want it to have a custom vanity address.",
					},
					cli.StringFlag{
						Name:  "operator-name, on",
						Usage: "Name of the operator",
					},
					cli.StringFlag{
						Name:  "operator-rewarder-address, ora",
						Usage: "EL Address where operator will get rewards",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Validate flags
					if c.String("amount") != "" {
						if _, err := cliutils.ValidateDepositEthAmount("deposit amount", c.String("amount")); err != nil {
							return err
						}
					}
					if c.String("max-slippage") != "" && c.String("max-slippage") != "auto" {
						if _, err := cliutils.ValidatePercentage("maximum commission rate slippage", c.String("max-slippage")); err != nil {
							return err
						}
					}
					if c.String("salt") != "" {
						if _, err := cliutils.ValidateBigInt("salt", c.String("salt")); err != nil {
							return err
						}
					}

					// Run
					return nodeDeposit(c)

				},
			},

			{
				Name:      "send",
				Aliases:   []string{"n"},
				Usage:     "Send ETH or tokens from the node account to an address. ENS names supported.",
				UsageText: "rocketpool node send [options] amount token to",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm token send",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 3); err != nil {
						return err
					}
					amount, err := cliutils.ValidatePositiveEthAmount("send amount", c.Args().Get(0))
					if err != nil {
						return err
					}
					token, err := cliutils.ValidateTokenType("token type", c.Args().Get(1))
					if err != nil {
						return err
					}

					// Run
					return nodeSend(c, amount, token, c.Args().Get(2))

				},
			},

			{
				Name:      "set-voting-delegate",
				Aliases:   []string{"sv"},
				Usage:     "Set the address you want to use when voting on Rocket Pool governance proposals, or the address you want to delegate your voting power to.",
				UsageText: "rocketpool node set-voting-delegate address",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm delegate setting",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 1); err != nil {
						return err
					}
					delegate := c.Args().Get(0)

					// Run
					return nodeSetVotingDelegate(c, delegate)

				},
			},
			{
				Name:      "clear-voting-delegate",
				Aliases:   []string{"cv"},
				Usage:     "Remove the address you've set for voting on Rocket Pool governance proposals.",
				UsageText: "rocketpool node clear-voting-delegate",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm delegate clearing",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return nodeClearVotingDelegate(c)

				},
			},

			{
				Name:      "initialize-fee-distributor",
				Aliases:   []string{"z"},
				Usage:     "Create the fee distributor contract for your node, so you can withdraw priority fees and MEV rewards after the merge",
				UsageText: "rocketpool node initialize-fee-distributor",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm initialization gas costs",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return initializeFeeDistributor(c)

				},
			},

			{
				Name:      "distribute-fees",
				Aliases:   []string{"b"},
				Usage:     "Distribute the priority fee and MEV rewards from your fee distributor to your withdrawal address and the rETH contract (based on your node's average commission)",
				UsageText: "rocketpool node distribute-fees",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm distribution",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return distribute(c)

				},
			},

			{
				Name:      "join-smoothing-pool",
				Aliases:   []string{"js"},
				Usage:     "Opt your node into the Smoothing Pool",
				UsageText: "rocketpool node join-smoothing-pool",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm opt-in",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return joinSmoothingPool(c)

				},
			},

			{
				Name:      "leave-smoothing-pool",
				Aliases:   []string{"ls"},
				Usage:     "Leave the Smoothing Pool",
				UsageText: "rocketpool node leave-smoothing-pool",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "Automatically confirm opt-out",
					},
				},
				Action: func(c *cli.Context) error {

					// Validate args
					if err := cliutils.ValidateArgCount(c, 0); err != nil {
						return err
					}

					// Run
					return leaveSmoothingPool(c)

				},
			},

			{
				Name:      "sign-message",
				Aliases:   []string{"sm"},
				Usage:     "Sign an arbitrary message with the node's private key",
				UsageText: "rocketpool node sign-message [-m message]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "message, m",
						Usage: "The 'quoted message' to be signed",
					},
				},
				Action: func(c *cli.Context) error {
					// Run
					return signMessage(c)
				},
			},
		},
	})
}