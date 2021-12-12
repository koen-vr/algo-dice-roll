package contract

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"

	net "github.com/koen-vr/algo-dice-roll/shared/network"
)

type Config struct {
	Manager crypto.Account

	ClearProg    string
	ApprovalProg string

	LocalSchema  types.StateSchema
	GlobalSchema types.StateSchema
}

func GetAppId(name string) (uint64, error) {
	return loadFromJsonFile(name)
}

func saveToJsonFile(name string, id uint64) error {
	str, err := json.Marshal(id)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(fmt.Sprintf(
		"%s/%s.id", net.Path(), name,
	), str, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func loadFromJsonFile(name string) (uint64, error) {
	id := uint64(0)
	data, err := os.ReadFile(fmt.Sprintf(
		"%s/%s.id", net.Path(), name,
	))
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(data[:], &id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
