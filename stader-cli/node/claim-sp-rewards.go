package node

import (
	"fmt"
	"github.com/stader-labs/stader-node/shared/services/gas"
	"github.com/stader-labs/stader-node/shared/services/stader"
	cliutils "github.com/stader-labs/stader-node/shared/utils/cli"
	"github.com/urfave/cli"
	"math/big"
	"strconv"
	"strings"
)

func ClaimSpRewards(c *cli.Context, downloadMerkleProofs bool) error {
	staderClient, err := stader.NewClientFromCtx(c)
	if err != nil {
		return err
	}
	defer staderClient.Close()

	// Check and assign the EC status
	err = cliutils.CheckClientStatus(staderClient)
	if err != nil {
		return err
	}

	// prompt user to select the cycles to claim from
	canClaimSpRewards, err := staderClient.CanClaimSpRewards()
	if err != nil {
		return err
	}
	if canClaimSpRewards.SocializingPoolContractPaused {
		fmt.Println("The socializing pool contract is paused!")
		return nil
	}
	if len(canClaimSpRewards.UnclaimedCycles) == 0 {
		fmt.Println("No cycles to claim!")
		return nil
	}

	if len(canClaimSpRewards.CyclesToDownload) > 0 && !downloadMerkleProofs {
		fmt.Println("You have unclaimed rewards from cycles that you have not downloaded the merkle proofs for. Please download the merkle proofs for these cycles before claiming your rewards.")
		return nil
	}

	if downloadMerkleProofs {
		fmt.Println("Downloading the merkle proofs for the cycles you have unclaimed rewards for...")
		if len(canClaimSpRewards.CyclesToDownload) > 0 {
			if !(c.Bool("yes") || cliutils.Confirm(fmt.Sprintf(
				"Are you sure you want to download the merkle proofs for cycles %v?", canClaimSpRewards.CyclesToDownload))) {
				fmt.Println("Cancelled.")
				return nil
			}
			_, err := staderClient.DownloadSpMerkleProofs()
			if err != nil {
				return err
			}
		}
		fmt.Println("Merkle proofs downloaded!")
	}

	fmt.Println("Following are the unclaimed cycles, Please enter in a comma seperated string the cycles you want to claim rewards for:")

	cyclesToClaim := map[int64]bool{}
	for {
		// TODO - add cycle info like amount of rewards. These are UX things and can be done later
		for i, cycle := range canClaimSpRewards.UnclaimedCycles {
			fmt.Printf("%d) %d\n", i, cycle.Int64())
		}

		cycleSelection := cliutils.Prompt("Which cycles would you like to claim? Use a comma separated list (such as '1,2,3') or leave it blank to claim all cycles at once.", "^$|^\\d+(,\\d+)*$", "Invalid cycle selection")
		if cycleSelection == "" {
			for _, cycle := range canClaimSpRewards.UnclaimedCycles {
				cyclesToClaim[cycle.Int64()] = true
			}
			break
		} else {
			elements := strings.Split(cycleSelection, ",")
			allValid := true
			for _, element := range elements {
				cycle, err := strconv.ParseUint(element, 0, 64)
				if err != nil {
					fmt.Printf("Unable to parse element: %s", element)
					allValid = false
				}

				// check if unclaimedCycles contains the cycle
				found := false
				for _, unclaimedCycle := range canClaimSpRewards.UnclaimedCycles {
					if unclaimedCycle.Int64() == int64(cycle) {
						found = true
						break
					}
				}
				if !found {
					fmt.Printf("Cycle %d is not in the list of unclaimed cycles. Please enter a valid cycle number", cycle)
				} else {
					cyclesToClaim[int64(cycle)] = true
				}
			}

			if allValid {
				break
			}
		}
	}

	// convert the map to a cycle big int array
	cyclesToClaimArray := []*big.Int{}
	for cycle := range cyclesToClaim {
		cyclesToClaimArray = append(cyclesToClaimArray, big.NewInt(cycle))
	}

	// estimate gas
	fmt.Println("Estimating gas...")
	estimateGasResponse, err := staderClient.EstimateClaimSpRewardsGas(cyclesToClaimArray)
	if err != nil {
		return err
	}

	//Assign max fees
	err = gas.AssignMaxFeeAndLimit(estimateGasResponse.GasInfo, staderClient, c.Bool("yes"))
	if err != nil {
		return err
	}

	if !(c.Bool("yes") || cliutils.Confirm(fmt.Sprintf(
		"Are you sure you want to claim the rewards for cycles %v? (y/n)", cyclesToClaimArray))) {
		fmt.Println("Cancelled.")
		return nil
	}

	fmt.Printf("Claiming rewards for cycles %v\n", cyclesToClaimArray)
	res, err := staderClient.ClaimSpRewards(cyclesToClaimArray)
	if err != nil {
		return err
	}

	cliutils.PrintTransactionHash(staderClient, res.TxHash)
	_, err = staderClient.WaitForTransaction(res.TxHash)
	if err != nil {
		return err
	}
	fmt.Println("Rewards claimed! You can check your balance with the node status command.")

	return nil
}