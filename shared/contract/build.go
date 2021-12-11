package contract

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/algorand/go-algorand-sdk/logic"

	exec "github.com/koen-vr/algo-prng-roller/shared/execute"
	net "github.com/koen-vr/algo-prng-roller/shared/network"
)

func Build(list []string) error {
	for _, s := range list {
		if err := build(s); nil != err {
			return err
		}
		if err := compile(s); nil != err {
			return err
		}
	}
	return nil
}

func CleanUp(list []string) error {
	for _, s := range list {
		if err := cleanUp(s); nil != err {
			return err
		}
	}
	return nil
}

func cleanUp(name string) error {
	if _, err := exec.List([]string{"-c", fmt.Sprintf(
		"rm -f ./contracts/%s.teal", name,
	)}); nil != err {
		return err
	}
	if _, err := exec.List([]string{"-c", fmt.Sprintf(
		"rm -f ./contracts/%s.prog", name,
	)}); nil != err {
		return err
	}
	return nil
}

func build(name string) error {
	if _, err := exec.List([]string{"-c", fmt.Sprintf(
		"python3 ./contracts/%s.py > ./contracts/%s.teal", name, name,
	)}); nil != err {
		return fmt.Errorf("build %s failed: %s", name, err)
	}
	return nil
}

func compile(name string) error {
	cln, err := net.MakeClient()
	if err != nil {
		return fmt.Errorf("compile %s failed: make client: %s", name, err)
	}

	teal, err := ioutil.ReadFile(fmt.Sprintf(
		"./contracts/%s.teal", name,
	))
	if err != nil {
		return fmt.Errorf("compile %s failed: read file: %s", name, err)
	}

	chk, err := cln.TealCompile(teal).Do(context.Background())
	if err != nil {
		return fmt.Errorf("compile %s failed: compile teal: %s", name, err)
	}
	prg, err := base64.StdEncoding.DecodeString(chk.Result)
	if err != nil {
		return fmt.Errorf("compile %s failed: decode program: %s", name, err)
	}
	err = logic.CheckProgram(prg, make([][]byte, 0))
	if nil != err {
		return fmt.Errorf("compile %s failed: check program: %s", name, err)
	}

	if err = ioutil.WriteFile(fmt.Sprintf(
		"./contracts/%s.prog", name,
	), prg, os.ModePerm); nil != err {
		return fmt.Errorf("compile %s failed: write file: %s", name, err)
	}

	return nil
}
