package po

import (
	"dashboard-server/dbutils"
	"dashboard-server/internal/stream/po"
	"github.com/pkg/errors"
)

type Session struct {
	CreatedByAccount string `gorm:"primaryKey"`
	Timestamp        uint64 `gorm:"primaryKey;autoIncrement:false"`
	Title            string
	SessionDuration  uint64
}

type UserSession struct {
	Username         string
	Account          string
	SessionTimestamp uint64
	Session          Session `gorm:"foreignKey:Account,SessionTimestamp;references:CreatedByAccount,Timestamp"`
}

func BindUserSession(
	account string,
	sessionTimestamp uint64,
	username string) (*UserSession, error) {

	userSession := UserSession{
		Session:  Session{CreatedByAccount: account, Timestamp: sessionTimestamp},
		Username: username,
	}

	if err := dbutils.GetDB().Create(&userSession).Error; err != nil {
		return nil, errors.Wrapf(err, "db create user session '%v' failed", userSession)
	}

	return &userSession, nil
}

func CreateSession(
	createdByAccount string,
	sessionTimestamp uint64,
	sessionName string,
	batchData *[]po.SensorData,
	sessionDuration uint64,
) (*Session, int64, error) {
	timestamp := sessionTimestamp
	session := Session{
		CreatedByAccount: createdByAccount,
		Title:            sessionName,
		Timestamp:        timestamp,
		SessionDuration:  sessionDuration,
	}

	if err := dbutils.GetDB().Create(&session).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "create session '%v' failed", session)
	}

	rowsAffected, err := po.CreateBatchSensorDataFromStructs(batchData)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "db batch create sensor data for session '%v' failed", session)
	}

	return &session, rowsAffected, nil
}

/**
Tab selection to specify user of interest

Session Metrics to be Shown:
--> set scope affects all metrics, default to of all time [ok do it]

Single statistic: <social media slide show statistic mode>
	total sessions danced
	total moves danced
	total time danced [ok do it]
	average accuracy (correctness) [ok do it]
	average group synchronisation [ok do it]
	stars collected [ok do it]
	danced with (users) [ok do it]
	local placement
	global placement
Trends:
	average group synchronisation delay [ok do it]
	total correctness/accuracy bar chart blue/red for right, wrong on same chart.
	Or can use same colour scheme with hovered percentages
	points accumulated^ can be represented above									[ok do it]
	--> set division filter by average/total of 'each month', 'each day', 'each session', 'each move (in a session)' [ok do it]
*/

func GetUserSessionTimestamps(account string, username string, timeStart uint64, timeEnd uint64) ([]uint64, error) {
	type Result struct {
		SessionTimestamp uint64
	}

	integerTimestamps := make([]uint64, 0)
	var results []*Result

	err := dbutils.GetDB().Model(&UserSession{}).Select("session_timestamp").
		Where(&UserSession{Account: account, Username: username}).
		Where("session_timestamp >= ?",timeStart).Where("session_timestamp <= ?", timeEnd).
		Find(&results).Error
	if err != nil {
		return nil, errors.Wrapf(err, "db get user session timestamps error")
	}

	for _, s := range results {
		integerTimestamps = append(integerTimestamps, s.SessionTimestamp)
	}

	return integerTimestamps, nil
}
