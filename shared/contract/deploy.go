package contract

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"

	net "github.com/koen-vr/algo-prng-roller/shared/network"
)

func Deploy(cfg Config) error {
	optIn := true

	clearProg, err := ioutil.ReadFile(fmt.Sprintf(
		"./contracts/%s.prog", cfg.ClearProg,
	))
	if err != nil {
		return fmt.Errorf("deploy failed: %s: read file: %s", cfg.ClearProg, err)
	}
	approvalProg, err := ioutil.ReadFile(fmt.Sprintf(
		"./contracts/%s.prog", cfg.ApprovalProg,
	))
	if err != nil {
		return fmt.Errorf("deploy failed: %s: read file: %s", cfg.ApprovalProg, err)
	}

	appArgs := [][]byte{}
	accounts := []string{}
	foreignApps := []uint64{}
	foreignAssets := []uint64{}

	cln, err := net.MakeClient()
	if err != nil {
		return fmt.Errorf("deploy failed: make client: %s", err)
	}
	txnParams, err := cln.SuggestedParams().Do(context.Background())
	if err != nil {
		return fmt.Errorf("deploy failed: suggested params: %s", err)
	}

	note := []byte{}
	group := types.Digest{}
	lease := [32]byte{}
	rekeyTo := types.ZeroAddress
	extraPages := uint32(0)

	createTx, err := future.MakeApplicationCreateTxWithExtraPages(
		optIn, approvalProg, clearProg, cfg.GlobalSchema, cfg.LocalSchema,
		appArgs, accounts, foreignApps, foreignAssets, txnParams,
		cfg.Manager.Address, note, group, lease, rekeyTo, extraPages,
	)
	if err != nil {
		return fmt.Errorf("deploy failed: make create tx: %s", err)
	}

	// Enforce it or fail, a bug?
	createTx.OnCompletion = types.OptInOC

	_, signedTx, err := crypto.SignTransaction(cfg.Manager.PrivateKey, createTx)
	if err != nil {
		return fmt.Errorf("deploy failed: sign create tx: %s", err)
	}
	pendingTx, err := cln.SendRawTransaction(signedTx).Do(context.Background())
	if err != nil {
		return fmt.Errorf("deploy failed: send create tx: %s", err)
	}

	txConfirm, err := net.WaitForConfirmation(cln, pendingTx, 24, context.Background())
	if err != nil {
		return fmt.Errorf("deploy failed: confirm tx: %s", err)
	}
	if len(txConfirm.PoolError) > 0 {
		return fmt.Errorf("deploy failed: pool error: %s", txConfirm.PoolError)
	}

	if err := saveToJsonFile(cfg.ApprovalProg, txConfirm.ApplicationIndex); err != nil {
		return fmt.Errorf("contract: failed to save app: %s", err)
	}

	return nil
}
