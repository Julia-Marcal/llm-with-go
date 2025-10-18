package dto

import "time"

type LLMAudit struct {
	Prompt      string
	Temperature float64
	Response    string
	Err         error
	Duration    time.Duration
	Timestamp   time.Time
}
