package state

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusRunning Status = "running"
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
)

// State is the state of the task
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

	Status Status `json:"status"`

	// LastError is the last error of the task
	LastError string `json:"last_error"`
}

// newState creates a new state for the task
func NewState(parameters map[string]any) *State {
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
func Deserialize(body []byte) (*State, error) {
	var state State
	err := json.Unmarshal(body, &state)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

// Serialize serializes the state to a byte array
func Serialize(state *State) ([]byte, error) {
	return json.Marshal(state)
}

// EncodeParameters encodes the parameters to a JSON string
func (s *State) EncodeParameters() (string, error) {
	parameters, err := json.Marshal(s.Parameters)
	if err != nil {
		return "", err
	}
	return string(parameters), nil
}
