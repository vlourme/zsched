package zsched

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type stateStatus string

const (
	StatusPending stateStatus = "pending"
	StatusRunning stateStatus = "running"
	StatusSuccess stateStatus = "success"
	StatusFailed  stateStatus = "failed"
)

// State is the State of the task
type State struct {
	// id is the id of the state
	ID uuid.UUID `json:"id"`

	// taskID is the id of the task
	TaskID uuid.UUID `json:"task_id"`

	// parentID is the id of the parent task
	ParentID uuid.UUID `json:"parent_id,omitempty"`

	// parameters is the parameters for the task
	Parameters map[string]any `json:"parameters"`

	// InitializedAt is the time the state was initialized
	InitializedAt time.Time `json:"initialized_at"`

	// StartedAt is the time the task started executing
	StartedAt time.Time `json:"started_at"`

	// Iteration is the current iteration of the task
	Iterations int `json:"iterations"`

	Status stateStatus `json:"status"`

	// LastError is the last error of the task
	LastError string `json:"last_error"`
}

// newState creates a new state for the task
func newState(parameters map[string]any) *State {
	return &State{
		TaskID:        uuid.New(),
		Parameters:    parameters,
		InitializedAt: time.Now(),
		Status:        StatusPending,
		Iterations:    0,
		LastError:     "",
	}
}

// recoverState recovers the state from the body
func deserializeState(body []byte) (*State, error) {
	var state State
	err := json.Unmarshal(body, &state)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

// Serialize serializes the state to a byte array
func (s *State) Serialize() ([]byte, error) {
	return json.Marshal(s)
}

// EncodeParameters encodes the parameters to a JSON string
func (s *State) EncodeParameters() (string, error) {
	parameters, err := json.Marshal(s.Parameters)
	if err != nil {
		return "", err
	}
	return string(parameters), nil
}

// GetStr returns the string value of the parameter by name
func (s *State) GetStr(name string, defaultValue ...string) string {
	value, ok := s.Parameters[name]
	if !ok {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
	return value.(string)
}

// GetFloat returns the float value of the parameter by name
func (s *State) GetFloat(name string, defaultValue ...float64) float64 {
	value, ok := s.Parameters[name].(float64)
	if !ok {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return value
}

// GetBool returns the bool value of the parameter by name
func (s *State) GetBool(name string, defaultValue ...bool) bool {
	value, ok := s.Parameters[name].(bool)
	if !ok {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}
	return value
}

// GetInt returns the int value of the parameter by name
func (s *State) GetInt(name string, defaultValue ...int) int {
	if len(defaultValue) > 0 {
		return int(s.GetFloat(name, float64(defaultValue[0])))
	}
	return int(s.GetFloat(name))
}
