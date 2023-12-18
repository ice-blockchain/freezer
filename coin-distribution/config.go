// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ice-blockchain/wintr/log"
)

func (cfg *config) EnsureValid() {
	if cfg.Workers == 0 {
		log.Panic("workers must be > 0")
	}
	if cfg.BatchSize == 0 || cfg.BatchSize > batchMaxSizeLimit {
		log.Panic("batchSize must be > 0 and < " + strconv.Itoa(batchMaxSizeLimit))
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
	log.Panic(err, "ethereum.privateKey is invalid") //nolint:revive,nolintlint //.

	if cfg.Ethereum.ContractAddress == "" {
		log.Panic("ethereum.contractAddress must not be empty")
	}
}
