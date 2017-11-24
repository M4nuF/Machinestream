package main

type PayloadStatus string
type EventStatus string

type Event struct {
	Topic string `json:"topic"`
	Ref string `json:"page"`
	Payload Payload `json:"payload"`
	Event EventStatus `json:"event"`
	JoinRef string `json:"join_ref"`
}

type Payload struct {
	MachineId string `json:"machine_id"`
	Id string `json:"id"`
	Timestamp string `json:"timestamp"`
	Status PayloadStatus `json:"status"`
}

const (
	StatusIdle = "idle"
	StatusRunning = "running"
	StatusFinished = "finished"
	StatusErrorred = "errorred"
)


//{
//	"topic":"events",
//	"ref":null,
//	"payload": {
//		"timestamp":"2017-11-23T20:54:03.055792Z",
//		"status":"running",
//		"machine_id":"db9eb448-214b-481f-96fe-d1b883ec11a7",
//		"id":"ebc150d4-66b7-4580-b820-f33ca112737c"
//	},
//	"join_ref":null,
//	"event":"new"
//}
