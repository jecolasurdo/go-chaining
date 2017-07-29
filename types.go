package chaining

import "jecolasurdo/go-chaining/behavior"

// Chain is a mechanism for deferring execution of methods if an error condition has been received.
type Chain struct {
	AtomicFunc           func(*Chain, func(*interface{}) (*interface{}, error), ActionArg)
	LocalError           error
	PreviousActionResult *interface{}
}

// ActionArg is the information passed into a chain function that describe the intended behavior.
type ActionArg struct {
	Value    *interface{}
	Behavior behavior.InjectionOption
}
