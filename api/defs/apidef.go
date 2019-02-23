package defs

type UserCredential struct {
	Username string `json:"user_name"`
	Pwd      string `json:"pwd"`
}

type VideoInfo struct {
	ID           string
	AuthorID     int
	Name         string
	DisplayCtime string
}

type Comments struct {
	ID      string
	VideoID string
	Author  string
	Content string
}

type SimpleSession struct {
	Username string
	TTL      int64
}

type SignedUp struct {
	Success   bool   `json:"success"`
	SessionID string `json:"session_id"`
}
