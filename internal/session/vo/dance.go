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
	Wrong        int    `json:"wrong"`
	Star1        int    `json:"star1"`
	Star2        int    `json:"star2"`
	Star3        int    `json:"star3"`
	AvgSyncDelay uint64 `json:"avg_sync_delay"`
}

type GetDanceBuddiesResp struct {
	Usernames []string `json:"usernames"`
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

	danceBuddies, err := po.GetDanceBuddies(req.Account, timestamps)
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

	averageSyncDelay, err := po.GetAverageGroupSyncDelayByUser(req.Account, timestamps)
	if err != nil {
		return nil, errors.Wrap(err, "vo get avg group sync delay error")
	}

	return &GetDancePerformanceResp{AvgSyncDelay: averageSyncDelay,
		Wrong: accuracyList[0], Star1: accuracyList[1], Star2: accuracyList[2], Star3: accuracyList[3]}, nil
}
