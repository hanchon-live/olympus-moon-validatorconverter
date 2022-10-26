package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

type ConsensusPubKey struct {
	Key string `json:"key"`
}

type Validator struct {
	ConsensusPubKey ConsensusPubKey `json:"consensus_pubkey"`
	OperatorAddress string          `json:"operator_address"`
}

type ValidatorsResponse struct {
	Validators []Validator `json:"validators"`
}

const (
	Bech32Prefix         = "evmos"
	Bech32PrefixAccAddr  = Bech32Prefix
	Bech32PrefixAccPub   = Bech32Prefix + sdk.PrefixPublic
	Bech32PrefixValAddr  = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	Bech32PrefixValPub   = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	Bech32PrefixConsAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	Bech32PrefixConsPub  = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

func pubkeyToValCons(pubkey string) (string, error) {
	sdk.GetConfig().SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	sdk.GetConfig().SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	sdk.GetConfig().SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)

	pubkeyBytes, err := b64.StdEncoding.DecodeString(pubkey)
	if err != nil {
		return "", err
	}
	address := crypto.Address(tmhash.SumTruncated(pubkeyBytes))
	conAddress := sdk.ConsAddress(address)
	return conAddress.String(), nil
}

func main() {
	content, err := os.ReadFile("validatorset.json")
	if err != nil {
		panic("validatorset.json file doesn't exists")
	}

	var validators ValidatorsResponse
	err = json.Unmarshal(content, &validators)
	if err != nil {
		panic("error on json unmarshal")
	}

	total := len(validators.Validators)

	res := `"res": [`
	for i, v := range validators.Validators {
		res = res + `{"operator_address:"` + v.OperatorAddress + `,`
		address, err := pubkeyToValCons(v.ConsensusPubKey.Key)
		if err != nil {
			panic("error converting pubkey to address")
		}
		res = res + `"consensus_address:"` + address + `}`
		if i != total-1 {
			res = res + `,`
		}
	}
	res = res + `]`
	fmt.Println(res)
}
