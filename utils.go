package main

func contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

type noStateError struct{}

func (n noStateError) Error() string {
	return "no available state for symbol"
}

type finalStateError struct{}

func (n finalStateError) Error() string {
	return "automaton is in final state"
}
