package models

import "encoding/binary"

const pendingStakeWithdrawalBytesSize = 36

var PendingStakeWithdrawalPrefix = GetBadgerHoldPrefix(PendingStakeWithdrawal{})

type PendingStakeWithdrawal struct {
	BatchID           Uint256
	FinalisationBlock uint32 `badgerhold:"index"`
}

func (s *PendingStakeWithdrawal) Bytes() []byte {
	data := make([]byte, pendingStakeWithdrawalBytesSize)
	copy(data[0:32], s.BatchID.Bytes())
	binary.BigEndian.PutUint32(data[32:36], s.FinalisationBlock)
	return data
}

func (s *PendingStakeWithdrawal) SetBytes(data []byte) error {
	if len(data) < pendingStakeWithdrawalBytesSize {
		return ErrInvalidLength
	}
	s.BatchID.SetBytes(data[0:32])
	s.FinalisationBlock = binary.BigEndian.Uint32(data[32:36])
	return nil
}
