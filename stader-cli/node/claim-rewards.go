package node

import (
	"fmt"
	"math/big"

	"github.com/stader-labs/stader-node/shared/services/gas"
	"github.com/stader-labs/stader-node/shared/services/stader"
	cliutils "github.com/stader-labs/stader-node/shared/utils/cli"
	"github.com/stader-labs/stader-node/stader-lib/utils/eth"
	"github.com/urfave/cli"
)

func ClaimRewards(c *cli.Context) error {
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

	// Print what network we're on
	err = cliutils.PrintNetwork(staderClient)

	// Check if we can Withdraw El Rewards
	canClaimRewardsResponse, err := staderClient.CanClaimRewards()
	if err != nil {
		return err
	}
	if canClaimRewardsResponse.NoRewards {
		fmt.Println("No rewards to claim.")
		return nil
	}

	sdStatusResponse, err := staderClient.GetSDStatus(big.NewInt(0))
	if err != nil {
		return err
	}

	sdStatus := sdStatusResponse.SDStatus

	totalFee := sdStatus.AccumulatedInterest

	if canClaimRewardsResponse.NonTerminalValidators == 0 {
		if sdStatus.SdCollateralCurrentAmount.Cmp(totalFee) < 0 {
			fmt.Printf("Since you exited all your validators, you require to have at least %s SD in your self bonded position to claim your rewards.", eth.DisplayAmountInUnits(totalFee, "sd"))
			return nil
		}
	} else {
		// if withdrawableInEth < claimsBalance, then there is an existing utilization position
		if canClaimRewardsResponse.ClaimsBalance.Cmp(canClaimRewardsResponse.WithdrawableInEth) != 0 {
			if sdStatusResponse.SDStatus.SdUtilizerLatestBalance.Cmp(big.NewInt(0)) > 0 {
				fmt.Printf("You need to first pay %s and close the utilization position to get back your funds. Execute the following command to repay your utilized SD stader-cli repay-sd --amount <SD amount> \n", eth.DisplayAmountInUnits(totalFee, "sd"))

				fmt.Printf("Based on the current Health Factor, you can claim upto %s.\n", eth.DisplayAmountInUnits(canClaimRewardsResponse.WithdrawableInEth, "eth"))

				fmt.Printf("Note: Please repay your utilized SD by using the following command to claim the remaining ETH: stader-cli sd repay --amount <amount of SD to be repaid>.\n")

				if !cliutils.Confirm("Are you sure you want to proceed?\n\n") {
					fmt.Println("Cancelled.")
					return nil
				}
			}
		}
	}

	err = gas.AssignMaxFeeAndLimit(canClaimRewardsResponse.GasInfo, staderClient, c.Bool("yes"))
	if err != nil {
		return err
	}

	// Prompt for confirmation
	if !(c.Bool("yes") || cliutils.Confirm(fmt.Sprintf(
		"Are you sure you want to send rewards to your operator reward address?"))) {
		fmt.Println("Cancelled.")
		return nil
	}

	// Withdraw El Rewards
	res, err := staderClient.ClaimRewards()
	if err != nil {
		return err
	}
	fmt.Printf("Withdrawing %s Rewards to Operator Reward Address: %s\n\n", eth.DisplayAmountInUnits(res.RewardsClaimed, "eth"), res.OperatorRewardAddress)
	cliutils.PrintTransactionHash(staderClient, res.TxHash)
	if _, err = staderClient.WaitForTransaction(res.TxHash); err != nil {
		return err
	}

	// Log & return
	fmt.Printf("Successful withdrawal of %s to Operator Reward Address: %s\n\n", eth.DisplayAmountInUnits(res.RewardsClaimed, "eth"), res.OperatorRewardAddress)
	return nil
}
