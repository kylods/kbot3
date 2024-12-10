package models

import "sync"

type UpdateCurrentTrack struct {
	Title    string
	Artist   string
	Duration uint
	UUID     string
}

type QueueObject struct {
	QueueUUID string `json:"queueUUID"`
	FileUUID  string `json:"-"`
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Duration  uint   `json:"duration"`
	QueuedBy  string `json:"queuedBy"`
}

type Queue struct {
	Tracks []QueueObject `json:"tracks"`
	Mutex  sync.Mutex
}
