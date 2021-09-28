package vo

import (
	"dashboard-server/internal/session/po"
	"github.com/pkg/errors"
)

type GetUserDanceReq struct {
	Start    uint64 `json:"start"`
	End      uint64 `json:"end"`
	Account  string `json:"account"`
	Username string `json:"username"`
}

type GetDanceDurationResp struct {
	Duration uint64 `json:"duration"`
}

type GetDancePerformanceResp struct {
	Wrong            int     `json:"wrong"`
	Star1            int     `json:"star1"`
	Star2            int     `json:"star2"`
	Star3            int     `json:"star3"`
	Total            int     `json:"total"`
	AvgSyncDelay     uint64  `json:"avg_sync_delay"`
	PositionAccuracy float32 `json:"position_accuracy"`
	TotalPositions   int     `json:"total_positions"`
	MoveAccuracy     float32 `json:"move_accuracy"`
}

type GetDanceBuddiesResp struct {
	Usernames []string `json:"usernames"`
}

type GetDanceOverviewResp struct {
	GetDancePerformanceResp
	GetDanceDurationResp
	GetDanceBuddiesResp
	TotalSessions int `json:"total_sessions"`
}

func GetDanceOverview(req GetUserDanceReq) (*GetDanceOverviewResp, error) {
	duration, err := GetDanceDuration(req)
	if err != nil {
		return nil, err
	}

	usernames, err := GetDanceBuddies(req)
	if err != nil {
		return nil, err
	}

	performance, err := GetDancePerformance(req)
	if err != nil {
		return nil, err
	}

	timestamps, err := po.GetUserSessionTimestamps(req.Account, req.Username, req.Start, req.End)
	totalSessions := len(timestamps)

	return &GetDanceOverviewResp{GetDancePerformanceResp: *performance,
		GetDanceBuddiesResp:  *usernames,
		GetDanceDurationResp: *duration,
		TotalSessions: totalSessions}, nil
}

func GetDanceDuration(req GetUserDanceReq) (*GetDanceDurationResp, error) {
	timestamps, err := po.GetUserSessionTimestamps(req.Account, req.Username, req.Start, req.End)
	if err != nil {
		return nil, errors.Wrap(err, "vo get user session timestamp error")
	}

	danceDuration, err := po.GetDanceDuration(req.Account, timestamps)
	if err != nil {
		return nil, errors.Wrap(err, "vo get dance duration error")
	}

	return &GetDanceDurationResp{Duration: danceDuration}, nil
}

func GetDanceBuddies(req GetUserDanceReq) (*GetDanceBuddiesResp, error) {
	timestamps, err := po.GetUserSessionTimestamps(req.Account, req.Username, req.Start, req.End)
	if err != nil {
		return nil, errors.Wrap(err, "vo get user session timestamp error")
	}

	danceBuddies, err := po.GetDanceBuddies(req.Account, timestamps, req.Username)
	if err != nil {
		return nil, errors.Wrap(err, "vo get dance buddies error")
	}

	return &GetDanceBuddiesResp{Usernames: danceBuddies}, nil
}

func GetDancePerformance(req GetUserDanceReq) (*GetDancePerformanceResp, error) {
	timestamps, err := po.GetUserSessionTimestamps(req.Account, req.Username, req.Start, req.End)
	if err != nil {
		return nil, errors.Wrap(err, "vo get user session timestamp error")
	}

	accuracyList, err := po.GetTotalAccuracy(req.Account, timestamps, req.Username)
	if err != nil {
		return nil, errors.Wrap(err, "vo get total accuracy error")
	}

	averageSyncDelay, err := po.GetAverageGroupSyncDelay(req.Account, timestamps)
	if err != nil {
		return nil, errors.Wrap(err, "vo get avg group sync delay error")
	}

	positionAccuracy, positionTotal, err := po.GetPositionAccuracy(req.Account, timestamps)
	if err != nil {
		return nil, errors.Wrap(err, "vo get position accuracy error")
	}

	total := accuracyList[0] + accuracyList[1] + accuracyList[2] + accuracyList[3]
	moveAccuracy := float32(0)
	if total > 0 {
		moveAccuracy = float32(total-accuracyList[0]) / float32(total)
	}

	return &GetDancePerformanceResp{AvgSyncDelay: averageSyncDelay,
		Wrong: accuracyList[0], Star1: accuracyList[1],
		Star2: accuracyList[2], Star3: accuracyList[3],
		Total:            accuracyList[4],
		PositionAccuracy: positionAccuracy,
		TotalPositions:   positionTotal,
		MoveAccuracy:     moveAccuracy,
		}, nil
}
