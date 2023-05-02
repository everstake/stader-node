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
package eth1

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stader-labs/stader-node/shared/services"
	"github.com/urfave/cli"
)

// Sets the nonce of the provided transaction options to the latest nonce if requested
func CheckForNonceOverride(c *cli.Context, opts *bind.TransactOpts) error {

	customNonceString := c.GlobalString("nonce")
	if customNonceString != "" {
		customNonce, success := big.NewInt(0).SetString(customNonceString, 0)
		if !success {
			return fmt.Errorf("Invalid nonce: %s", customNonceString)
		}

		// Do a sanity check to make sure the provided nonce is for a pending transaction
		// otherwise the user is burning gas for no reason
		ec, err := services.GetEthClient(c)
		if err != nil {
			return fmt.Errorf("Could not retrieve ETH1 client: %w", err)
		}

		// Make sure it's not higher than the next available nonce
		nextNonceUint, err := ec.PendingNonceAt(context.Background(), opts.From)
		if err != nil {
			return fmt.Errorf("Could not get next available nonce: %w", err)
		}

		nextNonce := big.NewInt(0).SetUint64(nextNonceUint)
		if customNonce.Cmp(nextNonce) == 1 {
			return fmt.Errorf("Can't use nonce %s because it's greater than the next available nonce (%d).", customNonceString, nextNonceUint)
		}

		// Make sure the nonce hasn't already been included in a block
		latestProposedNonceUint, err := ec.NonceAt(context.Background(), opts.From, nil)
		if err != nil {
			return fmt.Errorf("Could not get latest nonce: %w", err)
		}

		latestProposedNonce := big.NewInt(0).SetUint64(latestProposedNonceUint)
		if customNonce.Cmp(latestProposedNonce) == -1 {
			return fmt.Errorf("Can't use nonce %s because it has already been included in a block.", customNonceString)
		}

		// It points to a pending transaction, so this is a valid thing to do
		opts.Nonce = customNonce
	}
	return nil

}
