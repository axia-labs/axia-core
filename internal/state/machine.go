package state

import (
	"github.com/looplab/fsm"
	"github.com/sirupsen/logrus"
)

// StateManager handles the state transitions of graph elements
type StateManager struct {
	FSM *fsm.FSM
	log *logrus.Logger
}

// NodeState represents possible states for a node
const (
	StateInitial    = "initial"
	StateValidated  = "validated"
	StateProcessing = "processing"
	StateCompleted  = "completed"
	StateError      = "error"
)

// NewNodeStateManager creates a new state manager for nodes
func NewNodeStateManager(logger *logrus.Logger) *StateManager {
	return &StateManager{
		FSM: fsm.NewFSM(
			StateInitial,
			fsm.Events{
				{Name: "validate", Src: []string{StateInitial}, Dst: StateValidated},
				{Name: "process", Src: []string{StateValidated}, Dst: StateProcessing},
				{Name: "complete", Src: []string{StateProcessing}, Dst: StateCompleted},
				{Name: "error", Src: []string{"*"}, Dst: StateError},
			},
			fsm.Callbacks{
				"before_event": func(e *fsm.Event) {
					logger.WithFields(logrus.Fields{
						"from":  e.Src,
						"to":    e.Dst,
						"event": e.Event,
					}).Info("State transition starting")
				},
				"after_event": func(e *fsm.Event) {
					logger.WithFields(logrus.Fields{
						"from":  e.Src,
						"to":    e.Dst,
						"event": e.Event,
					}).Info("State transition completed")
				},
			},
		),
		log: logger,
	}
} 