package model

type RatingEvent struct {
	UserID     UserID          `json:"userId"`
	RecordID   RecordID        `json:"recordId"`
	RecordType RecordType      `json:"recordType"`
	Value      RatingValue     `json:"value"`
	EventType  RatingEventType `json:"eventType"`
}

type RatingEventType string

const (
	RatingEventTypePut    = "put"
	RatingEventTypeDelete = "delete"
)
const (
	RecordTypeMovie = RecordType("movie")
)

type RecordID string

type RecordType string

type UserID string

type RatingValue int

type Rating struct {
	RecordID   string      `json:"recordId"`
	RecordType string      `json:"recordType"`
	UserID     UserID      `json:"userId"`
	Value      RatingValue `json:"value"`
}
