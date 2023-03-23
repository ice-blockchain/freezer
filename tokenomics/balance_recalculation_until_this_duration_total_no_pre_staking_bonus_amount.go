// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"github.com/ice-blockchain/wintr/coin"
)

//nolint:funlen,gocognit,gocyclo,revive,cyclop // .
func (s *balanceRecalculationTriggerStreamSource) processTotalNoPreStakingBonusUntilThisDurationBalanceType(
	balancesByPK map[string]*balance, details *BalanceRecalculationDetails, untilThisDurationTypeDetail, userID string,
) {
	var positiveTotalAmount *coin.ICEFlake
	positiveTotalNoPreStakingBonusBalance := s.getBalance(false, "", balancesByPK)
	positiveTotalT0NoPreStakingBonusBalance := s.getBalance(false, details.t0TypeDetail(), balancesByPK)
	positiveTotalT1NoPreStakingBonusBalance := s.getBalance(false, t1BalanceTypeDetail, balancesByPK)
	positiveTotalT2NoPreStakingBonusBalance := s.getBalance(false, t2BalanceTypeDetail, balancesByPK)
	positiveUntilThisDurationBalance := s.getBalance(false, untilThisDurationTypeDetail, balancesByPK)
	if positiveTotalNoPreStakingBonusBalance != nil && !positiveTotalNoPreStakingBonusBalance.Amount.IsNil() {
		positiveTotalAmount = positiveTotalAmount.Add(positiveTotalNoPreStakingBonusBalance.Amount)
	}
	if positiveTotalT0NoPreStakingBonusBalance != nil && !positiveTotalT0NoPreStakingBonusBalance.Amount.IsNil() {
		positiveTotalAmount = positiveTotalAmount.Add(positiveTotalT0NoPreStakingBonusBalance.Amount)
	}
	if positiveTotalT1NoPreStakingBonusBalance != nil && !positiveTotalT1NoPreStakingBonusBalance.Amount.IsNil() {
		positiveTotalAmount = positiveTotalAmount.Add(positiveTotalT1NoPreStakingBonusBalance.Amount)
	}
	if positiveTotalT2NoPreStakingBonusBalance != nil && !positiveTotalT2NoPreStakingBonusBalance.Amount.IsNil() {
		positiveTotalAmount = positiveTotalAmount.Add(positiveTotalT2NoPreStakingBonusBalance.Amount)
	}
	if !positiveTotalAmount.IsNil() {
		positiveUntilThisDurationBalance = s.getOrInitBalance(false, untilThisDurationTypeDetail, userID, balancesByPK)
		positiveUntilThisDurationBalance.Amount = positiveTotalAmount
	} else if positiveUntilThisDurationBalance != nil && !positiveUntilThisDurationBalance.Amount.IsNil() {
		positiveUntilThisDurationBalance.Amount = coin.ZeroICEFlakes()
	}
	var negativeTotalAmount *coin.ICEFlake
	negativeTotalNoPreStakingBonusBalance := s.getBalance(true, "", balancesByPK)
	negativeTotalT0NoPreStakingBonusBalance := s.getBalance(true, details.t0TypeDetail(), balancesByPK)
	negativeTotalT1NoPreStakingBonusBalance := s.getBalance(true, t1BalanceTypeDetail, balancesByPK)
	negativeTotalT2NoPreStakingBonusBalance := s.getBalance(true, t2BalanceTypeDetail, balancesByPK)
	negativeUntilThisDurationBalance := s.getBalance(true, untilThisDurationTypeDetail, balancesByPK)
	if negativeTotalNoPreStakingBonusBalance != nil && !negativeTotalNoPreStakingBonusBalance.Amount.IsNil() {
		negativeTotalAmount = negativeTotalAmount.Add(negativeTotalNoPreStakingBonusBalance.Amount)
	}
	if negativeTotalT0NoPreStakingBonusBalance != nil && !negativeTotalT0NoPreStakingBonusBalance.Amount.IsNil() {
		negativeTotalAmount = negativeTotalAmount.Add(negativeTotalT0NoPreStakingBonusBalance.Amount)
	}
	if negativeTotalT1NoPreStakingBonusBalance != nil && !negativeTotalT1NoPreStakingBonusBalance.Amount.IsNil() {
		negativeTotalAmount = negativeTotalAmount.Add(negativeTotalT1NoPreStakingBonusBalance.Amount)
	}
	if negativeTotalT2NoPreStakingBonusBalance != nil && !negativeTotalT2NoPreStakingBonusBalance.Amount.IsNil() {
		negativeTotalAmount = negativeTotalAmount.Add(negativeTotalT2NoPreStakingBonusBalance.Amount)
	}
	if !negativeTotalAmount.IsNil() {
		negativeUntilThisDurationBalance = s.getOrInitBalance(true, untilThisDurationTypeDetail, userID, balancesByPK)
		negativeUntilThisDurationBalance.Amount = negativeTotalAmount
	} else if negativeUntilThisDurationBalance != nil && !negativeUntilThisDurationBalance.Amount.IsNil() {
		negativeUntilThisDurationBalance.Amount = coin.ZeroICEFlakes()
	}
}
