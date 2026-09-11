package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Priority int

const (
	Low Priority = iota
	Medium
	High
)

func (p Priority) String() string {
	switch p {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	}
	return "Invalid Priority"
}

func (p Priority) IsValid() bool {
	return p >= Low && p <= High
}

func ParsePriority(s string) (Priority, error) {
	lower := strings.ToLower(s)
	switch lower {
	case "l", "low":
		return Low, nil
	case "m", "medium":
		return Medium, nil
	case "h", "high":
		return High, nil
	}
	return 0, fmt.Errorf("Invalid priority: %v, must be low/l, medium/m, or high/h", lower)
}

// MarshalJSON converts Priority to JSON string
func (p Priority) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

// UnmarshalJSON converts JSON string to Priority
func (p *Priority) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	priority, err := ParsePriority(s)
	if err != nil {
		return err
	}

	*p = priority
	return nil
}
