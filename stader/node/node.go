/*
This work is licensed and released under GNU GPL v3 or any other later versions.
The full text of the license is below/ found at <http://www.gnu.org/licenses/>

(c) 2023 Rocket Pool Pty Ltd. Modified under GNU GPL v3. [0.3.0-exit]

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
package node

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli"

	"github.com/stader-labs/stader-node/shared/services"
	"github.com/stader-labs/stader-node/shared/services/config"
	"github.com/stader-labs/stader-node/shared/utils/log"
)

// TODO - increase the tasksInterval to 5m
var tasksInterval, _ = time.ParseDuration("10s")
var taskCooldown, _ = time.ParseDuration("10s")

const (
	MaxConcurrentEth1Requests = 200

	ManageFeeRecipientColor = color.FgHiCyan
	ErrorColor              = color.FgRed
)

// Register node command
func RegisterCommands(app *cli.App, name string, aliases []string) {
	app.Commands = append(app.Commands, cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   "Run Stader node activity daemon",
		Action: func(c *cli.Context) error {
			return run(c)
		},
	})
}

// Run daemon
func run(c *cli.Context) error {

	// Handle the initial fee recipient file deployment
	err := deployDefaultFeeRecipientFile(c)
	if err != nil {
		return err
	}

	// Configure
	configureHTTP()

	// Initialize tasks
	manageFeeRecipient, err := newManageFeeRecipient(c, log.NewColorLogger(ManageFeeRecipientColor))
	if err != nil {
		return err
	}

	// Initialize loggers
	errorLog := log.NewColorLogger(ErrorColor)

	// Wait group to handle the various threads
	wg := new(sync.WaitGroup)
	wg.Add(1)

	// Run task loop
	go func() {
		for {
			// Check the EC status
			err := services.WaitEthClientSynced(c, false) // Force refresh the primary / fallback EC status
			if err != nil {
				errorLog.Println(err)
			} else {
				// Check the BC status
				err := services.WaitBeaconClientSynced(c, false) // Force refresh the primary / fallback BC status
				if err != nil {
					errorLog.Println(err)
				} else {
					// Manage the fee recipient for the node
					if err := manageFeeRecipient.run(); err != nil {
						errorLog.Println(err)
					}
					time.Sleep(taskCooldown)
				}
			}
			time.Sleep(tasksInterval)
		}
		wg.Done()
	}()

	// Wait for both threads to stop
	wg.Wait()
	return nil

}

// Configure HTTP transport settings
func configureHTTP() {

	// The guardian daemon makes a large number of concurrent RPC requests to the Eth1 client
	// The HTTP transport is set to cache connections for future re-use equal to the maximum expected number of concurrent requests
	// This prevents issues related to memory consumption and address allowance from repeatedly opening and closing connections
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = MaxConcurrentEth1Requests

}

// Copy the default fee recipient file into the proper location
func deployDefaultFeeRecipientFile(c *cli.Context) error {

	cfg, err := services.GetConfig(c)
	if err != nil {
		return err
	}

	feeRecipientPath := cfg.StaderNode.GetFeeRecipientFilePath()
	_, err = os.Stat(feeRecipientPath)
	if os.IsNotExist(err) {
		// Make sure the validators dir is created
		validatorsFolder := filepath.Dir(feeRecipientPath)
		err = os.MkdirAll(validatorsFolder, 0755)
		if err != nil {
			return fmt.Errorf("could not create validators directory: %w", err)
		}

		// Create the file
		var defaultFeeRecipientFileContents string
		if cfg.IsNativeMode {
			// Native mode needs an environment variable definition
			defaultFeeRecipientFileContents = fmt.Sprintf("%s=%s", config.FeeRecipientEnvVar, cfg.StaderNode.GetEthxTokenAddress().Hex())
		} else {
			// Docker and Hybrid just need the address itself
			defaultFeeRecipientFileContents = cfg.StaderNode.GetEthxTokenAddress().Hex()
		}
		err := ioutil.WriteFile(feeRecipientPath, []byte(defaultFeeRecipientFileContents), 0664)
		if err != nil {
			return fmt.Errorf("could not write default fee recipient file to %s: %w", feeRecipientPath, err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking fee recipient file status: %w", err)
	}

	return nil

}
