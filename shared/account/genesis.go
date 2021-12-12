package account

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/types"

	exe "github.com/koen-vr/algo-dice-roll/shared/execute"
	net "github.com/koen-vr/algo-dice-roll/shared/network"
)

func FundFromSeed(amount uint64, to string) error {
	if 0 >= amount {
		return fmt.Errorf(
			"account: invalid amount: %d", amount,
		)
	}
	if !isValidAddress(to) {
		return fmt.Errorf(
			"account: invalid to address: %s", to,
		)
	}
	from, err := getSeedAddress()
	if nil != err {
		return fmt.Errorf(
			"account: get seed address: %s", err,
		)
	}
	if !isValidAddress(from) {
		return fmt.Errorf(
			"account: invalid seed address: %s", from,
		)
	}
	if _, err := exe.List([]string{"-c", fmt.Sprintf(
		"goal -d %s clerk send -a %d -f %s -t %s",
		net.NodePath(), amount*10000000000000, from, to,
	)}); nil != err {
		return fmt.Errorf("account: clerk send: %s", err)
	}
	return nil
}

func isValidAddress(addr string) bool {
	if _, err := types.DecodeAddress(addr); nil != err {
		return false
	}
	return true
}

func getSeedAddress() (string, error) {
	return exe.List([]string{"-c", fmt.Sprintf(
		"goal account list -d %s | awk '{ print $3 }' | head -n 1",
		net.NodePath(),
	)})
}
