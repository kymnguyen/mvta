package event

type UserAuthorizedEvent struct {
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	Timestamp int64  `json:"timestamp"`
}
