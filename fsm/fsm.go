package fsm

import "errors"

type State int

// Guard provides protection against transitioning to the goal State.
// Returning true/false indicates if the transition is permitted or not.
type Guard func(subject Stater, goal State) bool

var (
	InvalidTransition = errors.New("invalid transition")
)

// Transition is the change between States
type Transition interface {
	Origin() State
	Exit() State
}

// T implements the Transition interface; it provides a default
// implementation of a Transition.
type T struct {
	O, E State
}

func (t T) Origin() State { return t.O }
func (t T) Exit() State   { return t.E }

// RuleSet stores the rules for the state machine.
type RuleSet map[Transition][]Guard

// AddRule adds Guards for the given Transition
func (r RuleSet) AddRule(t Transition, guards ...Guard) {
	for _, guard := range guards {
		r[t] = append(r[t], guard)
	}
}

// AddTransition adds a transition with a default rule
func (r RuleSet) AddTransition(t Transition) {
	r.AddRule(t, func(subject Stater, goal State) bool {
		return subject.CurrentState() == t.Origin()
	})
}

// CreateRuleSet will establish a ruleset with the provided transitions.
// This eases initialization when storing within another structure.
func CreateRuleSet(transitions ...Transition) RuleSet {
	r := RuleSet{}

	for _, t := range transitions {
		r.AddTransition(t)
	}

	return r
}

// Permitted determines if a transition is allowed.
// This occurs in parallel.
// NOTE: Guards are not halted if they are short-circuited for some
// transition. They may continue running *after* the outcome is determined.
func (r RuleSet) Permitted(subject Stater, goal State) bool {
	attempt := T{subject.CurrentState(), goal}

	if guards, ok := r[attempt]; ok {
		outcome := make(chan bool)

		for _, guard := range guards {
			go func(g Guard) {
				outcome <- g(subject, goal)
			}(guard)
		}

		for range guards {
			select {
			case o := <-outcome:
				if !o {
					return false
				}
			}
		}

		return true // All guards passed
	}
	return false // No rule found for the transition
}

// Stater can be passed into the FSM. The Stater is reponsible for setting
// its own default state. Behavior of a Stater without a State is undefined.
type Stater interface {
	CurrentState() State
	SetState(State)
}

// Machine is a pairing of Rules and a Subject.
// The subject or rules may be changed at any time within
// the machine's lifecycle.
type Machine struct {
	Rules   *RuleSet
	Subject Stater
}

// Transition attempts to move the Subject to the Goal state.
func (m Machine) Transition(goal State) error {
	if m.Rules.Permitted(m.Subject, goal) {
		m.Subject.SetState(goal)
		return nil
	}

	return InvalidTransition
}

// New initializes a machine
func New(rules *RuleSet, subject Stater) *Machine {
	m := &Machine{
		Rules:   rules,
		Subject: subject,
	}
	return m
}
