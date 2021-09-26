package po

import (
	"dashboard-server/dbutils"
	"dashboard-server/internal/stream/po"
	"github.com/pkg/errors"
	"math"
)

func GetDanceDuration(account string, timestamps []uint64) (uint64, error) {
	type Result struct {
		Total uint64
	}

	var result Result
	if err := dbutils.GetDB().Model(&Session{}).Select("sum(session_duration) as total").
		Where(&Session{CreatedByAccount: account}).Where("timestamp in ?", timestamps).
		Group("created_by_account").
		First(&result).Error; err != nil {
		return 0, errors.Wrapf(err, "db get total time danced error")
	}

	return result.Total, nil
}

func GetTotalAccuracy(account string, timestamps []uint64, username string) ([]int, error) {
	type Result struct {
		Accuracy float32
		Count    int
	}

	var result []*Result

	if err := dbutils.GetDB().Model(&po.SensorData{}).Select("accuracy, count(accuracy) as count").
		Where(&po.SensorData{AccountName: account}).Where("session_timestamp in ?", timestamps).
		Where("username = ?", username).
		Group("accuracy").Find(&result).Error; err != nil {
		return nil, errors.Wrapf(err, "db get total accuracy error")
	}

	accuracyList := []int{0, 0, 0, 0}
	for _, r := range result {
		accuracyList[int(r.Accuracy)] = r.Count
	}

	return accuracyList, nil
}

func GetAverageGroupSyncDelayByUser(account string, timestamps []uint64) (uint64, error) {
	type Result struct {
		Average float64
	}

	var result Result

	if err := dbutils.GetDB().Model(&po.SyncDelay{}).
		Select("AVG(delay) as average").
		Where(&po.SyncDelay{Account: account}).Where("session_timestamp in ?", timestamps).
		Group("account").
		First(&result).Error; err != nil {
		return 0, errors.Wrapf(err, "db get average sync delay")
	}

	return uint64(math.Round(result.Average)), nil
}

func GetDanceBuddies(account string, timestamps []uint64) ([]string, error) {
	type Result struct {
		Username string
	}

	var results []*Result

	if err := dbutils.GetDB().Model(&UserSession{}).Select("DISTINCT(username)").
		Where(&UserSession{Account: account}).Where("session_timestamp in ?", timestamps).
		Find(&results).Error; err != nil {
		return nil, errors.Wrap(err, "db get danced with user failed")
	}

	usernames := make([]string, 0)
	for _, result := range results {
		usernames = append(usernames, result.Username)
	}
	return usernames, nil
}
