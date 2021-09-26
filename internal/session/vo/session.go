package vo

import (
	"dashboard-server/internal/session"
	"dashboard-server/internal/session/po"
	streamPo "dashboard-server/internal/stream/po"
	"github.com/pkg/errors"
)

type GetCurrentSessionResp struct {
	User1    []session.Packet   `json:"user_1"`
	User2    []session.Packet   `json:"user_2"`
	User3    []session.Packet   `json:"user_3"`
	Position []session.Position `json:"position"`
}

type UploadSessionReq struct {
	SessionTimestamp uint64                `json:"session_timestamp"`
	SessionName      string                `json:"session_name"`
	AccountName      string                `json:"account_name"`
	Username1        string                `json:"username_1"`
	Username2        string                `json:"username_2"`
	Username3        string                `json:"username_3"`
	SensorData       GetCurrentSessionResp `json:"sensor_data"`
}

type UploadSessionResp struct {
	NumSensorDataAdded int64       `json:"num_sensor_data_added"`
	Session            *po.Session `json:"session"`
}

func UploadSession(req UploadSessionReq) (*UploadSessionResp, error) {
	batchSensorData := make([]streamPo.SensorData,
		len(req.SensorData.User1)+len(req.SensorData.User2)+len(req.SensorData.User3))

	batchSensorDataPointer := 0
	sensorDataCombined := [][]session.Packet{req.SensorData.User1, req.SensorData.User2, req.SensorData.User3}
	Usernames := []string{req.Username1, req.Username2, req.Username3}
	sessionDuration := uint64(0)

	for index, userSensorData := range sensorDataCombined {
		for moveNum, sensorData := range userSensorData {
			trueAccuracy := float32(0)
			if sensorData.End == "correct" {
				trueAccuracy = sensorData.Accuracy
			}
			batchSensorData[batchSensorDataPointer] = streamPo.SensorData{
				AccountName:      req.AccountName,
				Username:         Usernames[index],
				MoveNum:          uint32(moveNum),
				SessionTimestamp: req.SessionTimestamp,
				PacketTimestamp:  sensorData.EpochMs,
				DanceMove:        sensorData.DanceMove,
				Accuracy:         trueAccuracy}
			batchSensorDataPointer++

			if moveNum == len(userSensorData)-1 && sensorData.EpochMs-userSensorData[0].EpochMs > sessionDuration {
				sessionDuration = sensorData.EpochMs - userSensorData[0].EpochMs
			}

			if index == 0 && len(sensorDataCombined[1]) > moveNum && len(sensorDataCombined[2]) > moveNum { //only once
				syncDelay := streamPo.ComputeSyncDelay([]uint64{
					sensorDataCombined[0][moveNum].EpochMs,
					sensorDataCombined[1][moveNum].EpochMs,
					sensorDataCombined[2][moveNum].EpochMs})
				_, err := streamPo.CreateSyncDelay(req.AccountName, req.SessionTimestamp, uint32(moveNum), syncDelay)
				if err != nil {
					return nil, errors.Wrapf(err, "unable to create sync delay for session")
				}
			}
		}
	}

	session, totalSensorData, err := po.CreateSession(req.AccountName, req.SessionTimestamp, req.SessionName, &batchSensorData, sessionDuration)
	if err != nil {
		return nil, errors.Wrap(err, "create session failed")
	}

	for index := range sensorDataCombined {
		if _, err := po.BindUserSession(req.AccountName, req.SessionTimestamp, Usernames[index]); err != nil {
			return nil, errors.Wrap(err, "bind user session error")
		}
	}

	return &UploadSessionResp{Session: session, NumSensorDataAdded: totalSensorData}, nil
}
