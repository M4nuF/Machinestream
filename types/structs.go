package types

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
