package po

import (
	"dashboard-server/dbutils"
	"github.com/pkg/errors"
)

type SyncDelay struct {
	Account          string
	SessionTimestamp uint64
	MoveNum          uint32
	Delay            uint64
}

func ComputeSyncDelay(timestamps []uint64) uint64 {
	maxTimestamp := uint64(0)
	minTimestamp := ^uint64(0) //90s is maximum
	for _, timestamp := range timestamps {
		if timestamp > maxTimestamp {
			maxTimestamp = timestamp
		}
		if timestamp < minTimestamp {
			minTimestamp = timestamp
		}
	}

	return maxTimestamp - minTimestamp
}

func CreateSyncDelay(
	account string,
	sessionTimestamp uint64,
	moveNum uint32,
	delay uint64,
) (*SyncDelay, error) {
	syncDelay := SyncDelay{Account: account,
		SessionTimestamp: sessionTimestamp,
		MoveNum:          moveNum,
		Delay:            delay}
	err := dbutils.GetDB().
		Create(&syncDelay).Error
	if err != nil {
		return nil, errors.Wrap(err, "db create sync delay error")
	}

	return &syncDelay, nil
}
