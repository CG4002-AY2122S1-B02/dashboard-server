package po

import (
	"dashboard-server/dbutils"
	"dashboard-server/internal/session"
	"github.com/pkg/errors"
)

type SensorData struct {
	AccountName      string `gorm:"primaryKey"`
	Username         string `gorm:"primaryKey"`
	MoveNum          uint32
	SessionTimestamp uint64 `gorm:"primaryKey;autoIncrement:false"`
	PacketTimestamp  uint64 `gorm:"primaryKey;autoIncrement:false"`
	DanceMove        string
	Accuracy         float32
}

//sync delay can be computed by comparing all packetTimestamp of sensorData with the same move_num

func CreateSensorDataStruct(
	packet session.Packet,
	accountName string,
	username string,
	sessionTimestamp uint64,
	moveNum uint32) *SensorData {
	sensorData := SensorData{
		MoveNum:          moveNum,
		Username:         username,
		AccountName:      accountName,
		SessionTimestamp: sessionTimestamp,
		PacketTimestamp:  packet.EpochMs,
		DanceMove:        packet.DanceMove,
		Accuracy:         packet.Accuracy,
	}

	return &sensorData
}

func CreateSensorDataFromStruct(sensorData *SensorData) error {

	if err := dbutils.GetDB().Create(sensorData).Error; err != nil {
		return errors.Wrapf(err, "create sensor data '%v' failed", sensorData)
	}

	return nil
}

func CreateBatchSensorDataFromStructs(sensorData *[]SensorData) (int64, error) {

	ret := dbutils.GetDB().Create(sensorData)
	if err := ret.Error; err != nil {
		return 0, errors.Wrapf(err, "create batch sensor data '%v' failed", sensorData)
	}

	return ret.RowsAffected, nil
}
