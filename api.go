// Package chaining provides a mechanism for chaining functions while deferring error handling to the end of execution.
package chaining

import (
	"jecolasurdo/go-chaining/behavior"
)

// New returns an instance of a chaining Context.
func New() *Chain {
	return &Chain{
		AtomicFunc: atomic,
	}
}

// Flush returns the context's error and final result, and resets the context back to its default state.
func (c *Chain) Flush() (*interface{}, error) {
	localError := c.LocalError
	finalResult := c.PreviousActionResult
	c.LocalError = nil
	c.PreviousActionResult = nil
	return finalResult, localError
}

// N (Nullary)executes an action which takes no arguments and returns only an error.
//
// Since the supplied action returns no value aside from an error, the context will supply nil as a pseudo-result
// for the supplied action. The context will then apply that nil to the next action called within the context,
// unless the override behavior for the next action dictates otherwise.
//
// The error returned by the supplied action is also applied to the current context.
// If error is not nil, subsequent actions executed within the same context will be ignored.
func (c *Chain) N(action func() error, behavior behavior.InjectionOption) {
	restatedAction := func(val *interface{}) (*interface{}, error) {
		return nil, action()
	}
	c.AtomicFunc(c, restatedAction, ActionArg{Behavior: behavior})
}

// NIface (Nullary Interface) executes an action which takes no arguments and returns a tuple of (interface{}, error).
//
// The context will apply the supplied action's interface{} result to the next action called within the context,
// unless the override behavior for the next action dictates otherwise.
//
// The error returned by the supplied action is also applied to the current context.
// If error is not nil, subsequent actions executed within the same context will be ignored.
func (c *Chain) NIface(action func() (*interface{}, error), behavior behavior.InjectionOption) {
	restatedAction := func(val *interface{}) (*interface{}, error) {
		return action()
	}
	c.AtomicFunc(c, restatedAction, ActionArg{Behavior: behavior})
}

// U (Unary) executes an action which takes one argument returns only an error.
//
// The single interface{} argument accepted by the action can be supplied via the arg parameter.
//
// Since the supplied action returns no value aside from an error, the context will supply nil as a pseudo-result
// for the supplied action. The context will then apply that nil to the next action called within the context,
// unless the override behavior for the next action dictates otherwise.
//
// The error returned by the supplied action is also applied to the current context.
// If error is not nil, subsequent actions executed within the same context will be ignored.
func (c *Chain) U(action func(*interface{}) error, arg ActionArg) {
	restatedAction := func(val *interface{}) (*interface{}, error) {
		return nil, action(val)
	}
	c.AtomicFunc(c, restatedAction, arg)
}

// UIface (Unary Interface) executes an action which takes one argument, and returns a tuple of (interface{}, error).
//
// The single interface{} argument accepted by the action can be supplied via the arg parameter.
//
// The context will apply the supplied action's interface{} result to the next action called within the context,
// unless the override behavior for the next action dictates otherwise.
//
// The error returned by the supplied action is also applied to the current context.
// If error is not nil, subsequent actions executed within the same context will be ignored.
func (c *Chain) UIface(action func(*interface{}) (*interface{}, error), arg ActionArg) {
	c.AtomicFunc(c, action, arg)
}

// NBool (Nullary Bool) executes an action which takes no arguments and returns a tuple of (bool, error).
//
// The context will apply the supplied action's bool result to the next action called within the context,
// unless the override behavior for the next action dictates otherwise.
//
// The error returned by the supplied action is also applied to the current context.
// If error is not nil, subsequent actions executed within the same context will be ignored.
//
// In addition to threading the (bool, error) tuple into the current context, NullaryBool itself also returns a bool.
// This is useful for inlining the method in boolean statements.
func (c *Chain) NBool(action func() (*bool, error), behavior behavior.InjectionOption) bool {
	restatedAction := func(*interface{}) (*bool, error) {
		return action()
	}
	return c.UBool(restatedAction, ActionArg{Behavior: behavior})
}

// UBool (Unary Bool) executes an action which takes one argument and returns a tuple of (bool, error).
//
// The single interface{} argument accepted by the action can be supplied via the arg parameter.
//
// The context will apply the supplied action's bool result to the next action called within the context,
// unless the override behavior for the next action dictates otherwise.
//
// The error returned by the supplied action is also applied to the current context.
// If error is not nil, subsequent actions executed within the same context will be ignored.
//
// In addition to threading the (bool, error) tuple into the current context, UnaryBool itself also returns a bool.
// This is useful for inlining the method in boolean statements.
func (c *Chain) UBool(action func(*interface{}) (*bool, error), arg ActionArg) bool {
	restatedAction := func(val *interface{}) (*interface{}, error) {
		r, err := action(val)
		var result interface{} = r
		return &result, err
	}
	c.AtomicFunc(c, restatedAction, arg)
	if c.LocalError != nil {
		return false
	}
	return *((*c.PreviousActionResult).(*bool))
}

func atomic(c *Chain, action func(*interface{}) (*interface{}, error), arg ActionArg) {
	if c.LocalError != nil {
		return
	}
	var valueToInject *interface{}
	if arg.Behavior == behavior.InjectSuppliedValue {
		valueToInject = arg.Value
	} else {
		valueToInject = c.PreviousActionResult
	}
	result, err := action(valueToInject)
	c.LocalError = err
	c.PreviousActionResult = result
}
