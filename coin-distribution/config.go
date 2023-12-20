// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/log"
)

func (cfg *config) EnsureValid() {
	if cfg.Workers == 0 {
		log.Panic("workers must be > 0")
	}
	if cfg.Ethereum.ChainID == 0 {
		log.Panic("ethereum.chainID must be > 0")
	}
	if cfg.Ethereum.RPC == "" {
		log.Panic("ethereum.rpc must not be empty")
	}
	if cfg.Ethereum.PrivateKey == "" {
		log.Panic("ethereum.privateKey must not be empty")
	}
	_, err := crypto.HexToECDSA(cfg.Ethereum.PrivateKey)
	log.Panic(errors.Wrap(err, "ethereum.privateKey is invalid")) //nolint:revive,nolintlint //.

	if cfg.Ethereum.ContractAddress == "" {
		log.Panic("ethereum.contractAddress must not be empty")
	}
}
