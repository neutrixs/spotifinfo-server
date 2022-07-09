package runtimedb

import "golang.org/x/exp/slices"

type EachState struct {
	Name string
	Scope string
}

type states struct {
	states []EachState
}

func (s states) Get() []EachState {
	return s.states
}

func (s states) GetIndex(stateName string) int {
	return slices.IndexFunc(s.states, func(state EachState) bool {
		return state.Name == stateName
	})
}

func (s states) Add(state EachState) {
	s.states = append(s.states, state)
}

var InitStates = &states {
	states: make([]EachState, 0),
}