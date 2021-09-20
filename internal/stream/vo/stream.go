package vo

type PostStreamCommandReq struct {
	Start            bool   `json:"start"`
	SessionTimestamp uint64 `json:"session_timestamp"`
	AccountName      string `json:"account_name"`
	Username1        string `json:"username1"`
	Username2        string `json:"username2"`
	Username3        string `json:"username3"`
}
