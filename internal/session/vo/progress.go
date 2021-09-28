package vo

import (
	"dashboard-server/internal/session/po"
	"fmt"
	"github.com/pkg/errors"
	"math"
	"time"
)

type ProgressUnit struct {
	Count int `json:"count"` //y axis
	Timestamp string `json:"timestamp"` //x axis
	Label string `json:"label"`
}

type GetDanceProgressResp struct {
	WrongMove []ProgressUnit `json:"wrong_move"`
	Star1 []ProgressUnit `json:"star_1"`
	Star2 []ProgressUnit `json:"star_2"`
	Star3 []ProgressUnit `json:"star_3"`
	CorrectMove []ProgressUnit `json:"correct_move"`
	WrongPosition []ProgressUnit `json:"wrong_position"`
	CorrectPosition []ProgressUnit `json:"correct_position"`
	GroupSyncDelay []ProgressUnit `json:"group_sync_delay"`
}

func GetDanceProgress(req GetUserDanceReq) (*GetDanceProgressResp, error) {
	resp := GetDanceProgressResp{
		WrongMove: make([]ProgressUnit, 0),
		Star1: make([]ProgressUnit, 0),
		Star2: make([]ProgressUnit, 0),
		Star3: make([]ProgressUnit, 0),
		CorrectMove: make([]ProgressUnit, 0),
		WrongPosition: make([]ProgressUnit, 0),
		CorrectPosition: make([]ProgressUnit, 0),
		GroupSyncDelay: make([]ProgressUnit, 0),
	}
	timestamps, err := po.GetUserSessionTimestamps(req.Account, req.Username, req.Start, req.End)
	if err != nil {
		return nil, errors.Wrap(err, "get dance progress")
	}

	for _, epochMs := range timestamps {
		timestamp := time.Unix(int64(epochMs/1000), 0).Format("02-Jan-2006 15:04:05")

		performance, err := GetDancePerformance(GetUserDanceReq{Start: epochMs, End: epochMs, Username: req.Username, Account: req.Account})
		if err != nil {
			return nil, errors.Wrap(err, "get dance progress")
		}
		positionLabel :=  fmt.Sprintf("Position Accuracy: %.1f%% (%v/%v)", performance.PositionAccuracy * 100,
			int(math.Round(float64(performance.PositionAccuracy * float32(performance.TotalPositions)))),
			performance.TotalPositions)
		resp.CorrectPosition = append(resp.CorrectPosition,
			ProgressUnit{Count: int(math.Round(float64(performance.PositionAccuracy * float32(performance.TotalPositions)))), Timestamp: timestamp,
				Label: positionLabel})
		resp.WrongPosition = append(resp.WrongPosition,
			ProgressUnit{Count: int(math.Round(float64((1 - performance.PositionAccuracy) * float32(performance.TotalPositions)))), Timestamp: timestamp,
				Label: positionLabel})

		moveLabel := fmt.Sprintf("Move Accuracy: %.1f%% (%v/%v)", float32(performance.Total - performance.Wrong) / float32(performance.Total) * 100, performance.Total - performance.Wrong, performance.Total)
		resp.WrongMove = append(resp.WrongMove,
			ProgressUnit{Count: performance.Wrong, Timestamp: timestamp,
				Label: moveLabel})
		resp.Star1 = append(resp.Star1,
			ProgressUnit{Count: performance.Star1, Timestamp: timestamp,
				Label: moveLabel})
		resp.Star2 = append(resp.Star2,
			ProgressUnit{Count: performance.Star2, Timestamp: timestamp,
				Label: moveLabel})
		resp.Star3 = append(resp.Star3,
			ProgressUnit{Count: performance.Star3, Timestamp: timestamp,
				Label: moveLabel})
		resp.CorrectMove = append(resp.CorrectMove,
			ProgressUnit{Count: performance.Total - performance.Wrong, Timestamp: timestamp,
				Label: moveLabel})

		resp.GroupSyncDelay = append(resp.GroupSyncDelay,
			ProgressUnit{Count: int(performance.AvgSyncDelay) , Timestamp: timestamp,
				Label: fmt.Sprintf("Avg Group Sync Delay: %vms", performance.AvgSyncDelay)})
	}

	return &resp, nil
}