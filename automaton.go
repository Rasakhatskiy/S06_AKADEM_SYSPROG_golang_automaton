package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type symbol rune

type state struct {
	id          int
	transitions []transition
	isFinal     bool
}

type transition struct {
	input  symbol
	output *state
}

type automaton struct {
	states      []state
	finalStates []*state

	stateStart   *state
	stateCurrent *state
}

func readInt(r *bufio.Reader) (int, error) {
	line, err := r.ReadString('\n')
	line = strings.TrimSuffix(line, "\n")
	if err != nil {
		return -1, err
	}
	i, err := strconv.ParseInt(line, 10, 32)
	if err != nil {
		return -1, err
	}
	return int(i), nil
}

func readAutomaton(r *bufio.Reader) (automaton, error) {
	// |A|
	symbolsSize, err := readInt(r)
	fmt.Println(symbolsSize)
	if err != nil {
		return automaton{}, err
	}

	// |S|
	statesSize, err := readInt(r)
	if err != nil {
		return automaton{}, err
	}

	// s0
	startStateId, err := readInt(r)
	fmt.Println(startStateId)
	if err != nil {
		return automaton{}, err
	}

	// |F| & F
	line, err := r.ReadString('\n')
	line = strings.TrimSuffix(line, "\n")
	if err != nil {
		return automaton{}, err
	}
	strgs := strings.Split(line, " ")
	if len(strgs) == 0 {
		return automaton{}, errors.New("invalid |F| value")
	}

	finalStatesSize, err := strconv.ParseInt(strgs[0], 10, 32)
	if err != nil {
		return automaton{}, err
	}
	if len(strgs) != int(finalStatesSize)+1 {
		return automaton{}, errors.New("invalid F value")
	}

	var finalStatesIds []int
	for i := 1; i <= int(finalStatesSize); i++ {
		id, err := strconv.ParseInt(strgs[i], 10, 32)
		if err != nil {
			return automaton{}, err
		}
		finalStatesIds = append(finalStatesIds, int(id))
	}

	type rawTransition struct {
		id     int
		s      symbol
		output int
	}

	var rawTransitions []rawTransition

	for {
		line, err := r.ReadString('\n')
		line = strings.TrimSuffix(line, "\n")
		if err == io.EOF && len(line) == 0 {
			break
		}
		if err != nil && err != io.EOF {
			return automaton{}, err
		}

		strgs = strings.Split(line, " ")
		if len(strgs) != 3 {
			return automaton{}, errors.New("invalid S value")
		}

		id, err := strconv.ParseInt(strgs[0], 10, 32)
		if err != nil {
			return automaton{}, err
		}
		s := strgs[1][0]
		out, err := strconv.ParseInt(strgs[2], 10, 32)
		if err != nil {
			return automaton{}, err
		}

		rawTransitions = append(rawTransitions, rawTransition{int(id), symbol(s), int(out)})
	}

	// set states
	automat := automaton{
		states:       []state{},
		finalStates:  []*state{},
		stateStart:   nil,
		stateCurrent: nil,
	}

	for i := 0; i < statesSize; i++ {
		automat.states = append(automat.states, state{
			id:          i,
			transitions: []transition{},
			isFinal:     contains(finalStatesIds, i),
		})
		if automat.states[len(automat.states)-1].isFinal {
			automat.finalStates = append(automat.finalStates, &automat.states[len(automat.states)-1])
		}
	}

	automat.stateStart = &automat.states[startStateId]
	automat.stateCurrent = automat.stateStart

	// set transitions
	for _, trans := range rawTransitions {
		automat.states[trans.id].transitions = append(automat.states[trans.id].transitions, transition{
			input:  trans.s,
			output: &automat.states[trans.output],
		})
	}

	return automat, nil
}

func (a *automaton) goToState(s symbol) error {
	for _, t := range a.stateCurrent.transitions {
		if t.input == s {
			a.stateCurrent = t.output
			return nil
		}
	}
	return noStateError{}
}

func (a *automaton) checkFirstWord(word []symbol) (bool, error) {
	a.stateCurrent = a.stateStart

	for _, s := range word {
		err := a.goToState(s)
		if err != nil {
			return false, err
		}
	}

	if a.stateCurrent.isFinal {
		return false, finalStateError{}
	}

	return true, nil
}

func (a *automaton) checkLastWord(word []symbol) (bool, error) {
	// todo 1. check all final states for last symbol
	// todo 2. check if there is a state that goes from * through symbol to prev state
	// todo 3. continue to beginning of the word
	// todo 4. get state for first symbol
	return false, nil
}
