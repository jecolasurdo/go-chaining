package chaining

import "jecolasurdo/go-chaining/injectionbehavior"
import "errors"

// Flush returns the context's error and final result, and resets the context back to its default state.
func (c *Context) Flush() (interface{}, error) {
	localError := c.LocalError
	finalResult := c.PreviousActionResult
	c.LocalError = nil
	c.PreviousActionResult = nil
	return finalResult, localError
}

// ApplyNullary executes an action which takes no arguments and returns only an error.
//
// Since the supplied action returns no value aside from an error, the context will supply nil as a pseudo-result
// for the supplied action. The context will then apply that nil to the next action called within the context,
// unless the override behavior for the next action dictates otherwise.
//
// The error returned by the supplied action is also applied to the current context.
// If error is not nil, subsequent actions executed within the same context will be ignored.
func (c *Context) ApplyNullary(action func() error) {
	if c.LocalError == nil {
		c.LocalError = action()
	}
}

// ApplyNullaryIface executes an action which takes no arguments and returns a tuple of (interface{}, error).
//
// The context will apply the supplied action's interface{} result to the next action called within the context,
// unless the override behavior for the next action dictates otherwise.
//
// The error returned by the supplied action is also applied to the current context.
// If error is not nil, subsequent actions executed within the same context will be ignored.
func (c *Context) ApplyNullaryIface(action func() (interface{}, error)) {
	c.LocalError = errors.New("ApplyNullaryIface not implemented")
	c.PreviousActionResult = nil
}

// ApplyUnary executes an action which takes one argument returns only an error.
//
// The single interface{} argument accepted by the action can be supplied via the arg parameter.
//
// Since the supplied action returns no value aside from an error, the context will supply nil as a pseudo-result
// for the supplied action. The context will then apply that nil to the next action called within the context,
// unless the override behavior for the next action dictates otherwise.
//
// The error returned by the supplied action is also applied to the current context.
// If error is not nil, subsequent actions executed within the same context will be ignored.
func (c *Context) ApplyUnary(action func(interface{}) error, arg ActionArg) {
	c.LocalError = errors.New("ApplyUnary not implemented")
	c.PreviousActionResult = nil
}

// ApplyUnaryIface executes an action which takes one argument, and returns a tuple of (interface{}, error).
//
// The single interface{} argument accepted by the action can be supplied via the arg parameter.
//
// The context will apply the supplied action's interface{} result to the next action called within the context,
// unless the override behavior for the next action dictates otherwise.
//
// The error returned by the supplied action is also applied to the current context.
// If error is not nil, subsequent actions executed within the same context will be ignored.
func (c *Context) ApplyUnaryIface(action func(interface{}) (interface{}, error), arg ActionArg) {
	if c.LocalError != nil {
		return
	}
	var valueToInject interface{}
	if arg.Behavior == injectionbehavior.InjectSuppliedValue {
		valueToInject = arg.Value
	} else {
		valueToInject = c.PreviousActionResult
	}
	result, err := action(valueToInject)
	c.LocalError = err
	c.PreviousActionResult = result
}

// ApplyNullaryBool executes an action which takes no arguments and returns a tuple of (bool, error).
//
// The context will apply the supplied action's bool result to the next action called within the context,
// unless the override behavior for the next action dictates otherwise.
//
// The error returned by the supplied action is also applied to the current context.
// If error is not nil, subsequent actions executed within the same context will be ignored.
//
// In addition to threading the (bool, error) tuple into the current context, NullaryBool itself also returns a bool.
// This is useful for inlining the method in boolean statements.
func (c *Context) ApplyNullaryBool(action func() (bool, error)) bool {
	if c.LocalError != nil {
		return false
	}

	result, err := action()
	c.LocalError = err
	return result
}

// ApplyUnaryBool executes an action which takes one argument and returns a tuple of (bool, error).
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
func (c *Context) ApplyUnaryBool(action func(interface{}) (bool, error), arg ActionArg) bool {
	c.LocalError = errors.New("ApplyUnaryBool not implemented")
	c.PreviousActionResult = nil
	return false
}
