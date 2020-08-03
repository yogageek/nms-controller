package model

import "fmt"

type ErrRes struct {
	Message     string `json:"error,omitempty"`
	Description string `json:"description,omitempty"`
}

func (e *ErrRes) Error() string {
	return fmt.Sprintf("%s:%s:", e.Message, e.Description)
}

// not used
type RedisMetric struct {
	UID    string  `json:"uid"`
	Value  float64 `json:"value"`
	Labels []Label `json:"labels"`
}
