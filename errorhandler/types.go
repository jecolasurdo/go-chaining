package errorhandler

import "jecolasurdo/go-deferrederrors/errorhandler/injectionbehavior"

// DeferredErrorContext is a mechanism for deferring execution of methods if an error condition has been received.
type DeferredErrorContext struct {
	LocalError     error
	injectionValue interface{}
}

// ActionArg is the information passed into a chain function that describe the intended behavior.
type ActionArg struct {
	Value    interface{}
	Behavior injectionbehavior.InjectionBehavior
}
