package model

import "time"

// UserRateLimit represents the data structure for storing rate limit information.
type UserRateLimit struct {
	Id           string    `json:"id"`
	UserId       string    `json:"userId"`
	RequestCount int       `json:"requestCount"`
	RateLimit    int       `json:"rateLimit"`
	Timestamp    time.Time `json:"timestamp"`
}
