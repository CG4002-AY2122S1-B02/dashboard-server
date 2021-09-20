package po

import (
	"dashboard-server/dbutils"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type Session struct {
	CreatedByAccount  string `gorm:"primaryKey"`
	CreatedByUsername string `gorm:"primaryKey"`
	Timestamp         uint64 `gorm:"primaryKey"`
	Title             string
	GroupName         string //group name
}

func CreateSession(createdByUsername string, createdByAccount string) (*Session, error) {
	timestamp := uint64(time.Now().UnixNano()) / uint64(time.Millisecond)
	session := Session{
		CreatedByAccount:  createdByAccount,
		CreatedByUsername: createdByUsername,
		Title:             fmt.Sprint(timestamp),
		Timestamp:         timestamp,
	}

	if err := dbutils.GetDB().Create(&session).Error; err != nil {
		return nil, errors.Wrapf(err, "create session '%v' failed", session)
	}

	return &session, nil
}

//func CreateSessionWith(createdBy string, groupName string) (*Session, error) {
//	timestamp := uint64(time.Now().UnixNano()) / uint64(time.Millisecond)
//	session := Session{
//		CreatedBy: createdBy,
//		Title: fmt.Sprint(timestamp),
//		Timestamp: timestamp,
//		GroupName: groupName,
//	}
//
//	if err := dbutils.GetDB().Create(&session).Error; err != nil {
//	return nil, errors.Wrapf(err, "create session '%v' failed", session)
//	}
//
//	return &session, nil
//}

func GetSession(createdByAccount, createdByUsername string, timestamp uint64) (*Session, error) {
	var (
		session Session
	)
	err := dbutils.GetDB().Where(&Session{
		CreatedByAccount: createdByAccount, CreatedByUsername: createdByUsername, Timestamp: timestamp}).First(&session).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "get session '%v' failed", session)
	}

	return &session, nil
}

func GetGroupSession(groupName string, timestamp uint64) (*Session, error) {
	var session Session
	err := dbutils.GetDB().Where(&Session{
		GroupName: groupName, Timestamp: timestamp}).First(&session).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "get session '%v' failed", session)
	}

	return &session, nil
}
