package chaining

import (
	"jecolasurdo/go-chaining/injectionbehavior"
)

// Context is a mechanism for deferring execution of methods if an error condition has been received.
type Context struct {
	AtomicFunc           func(*Context, func(interface{}) (interface{}, error), ActionArg)
	LocalError           error
	PreviousActionResult *interface{}
}

// ActionArg is the information passed into a chain function that describe the intended behavior.
type ActionArg struct {
	Value    interface{}
	Behavior injectionbehavior.InjectionBehavior
}
