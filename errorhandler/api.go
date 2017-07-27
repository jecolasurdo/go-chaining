package errorhandler

import (
	"jecolasurdo/go-deferrederrors/errorhandler/injectionbehavior"
)

// TryBool executes a func() (bool, error) if no previous Trys have resulted in an error.
// If previous Trys have resulted in an error, the action is ignored, not executed, and false is returned.
// Because false is returned when an action is ignored (rather than halting execution), it is important
// to ensure any downstream methods are also wrapped in Try methods, so they are also ignored.
func (d *DeferredErrorContext) TryBool(action func() (bool, error)) bool {
	if d.LocalError != nil {
		return false
	}

	result, err := action()
	d.LocalError = err
	return result
}

// TryVoid executes a func() error action if no previous Trys have resulted in a error.
// If previous Trys have resulted in an error, the action is ignored and not executed.
func (d *DeferredErrorContext) TryVoid(action func() error) {
	if d.LocalError == nil {
		d.LocalError = action()
	}
}

// ChainF does something
func (d *DeferredErrorContext) ChainF(action func(interface{}) (interface{}, error), arg ActionArg) {
	if d.LocalError != nil {
		return
	}
	var valueToInject interface{}
	if arg.Behavior == injectionbehavior.InjectSuppliedValue {
		valueToInject = arg.Value
	} else {
		valueToInject = d.PreviousActionResult
	}
	result, err := action(valueToInject)
	d.LocalError = err
	d.PreviousActionResult = result
}

// Flush returns the context's error and final result, and resets the context back to its default state.
func (d *DeferredErrorContext) Flush() (interface{}, error) {
	localError := d.LocalError
	finalResult := d.PreviousActionResult
	d.LocalError = nil
	d.PreviousActionResult = nil
	return finalResult, localError
}

// NullaryVoid executes an action which takes no arguments and returns only an error.
//
// Since the supplied action returns no value aside from an error, the context will supply nil as a pseudo-result
// for the supplied action. The context will then apply that nil to the next action called within the context,
// unless the override behavior for the next action dictates otherwise.
//
// The error returned by the supplied action is also applied to the current context.
// If error is not nil, subsequent actions executed within the same context will be ignored.
func (d *DeferredErrorContext) NullaryVoid(action func() error) {}

// NullaryIface executes an action which takes no arguments and returns a tuple of (interface{}, error).
//
// The context will apply the supplied action's interface{} result to the next action called within the context,
// unless the override behavior for the next action dictates otherwise.
//
// The error returned by the supplied action is also applied to the current context.
// If error is not nil, subsequent actions executed within the same context will be ignored.
func (d *DeferredErrorContext) NullaryIface(action func() (interface{}, error)) {}

// UnaryVoid executes an action which takes one argument returns only an error.
//
// The single interface{} argument accepted by the action can be supplied via the arg parameter.
//
// Since the supplied action returns no value aside from an error, the context will supply nil as a pseudo-result
// for the supplied action. The context will then apply that nil to the next action called within the context,
// unless the override behavior for the next action dictates otherwise.
//
// The error returned by the supplied action is also applied to the current context.
// If error is not nil, subsequent actions executed within the same context will be ignored.
func (d *DeferredErrorContext) UnaryVoid(action func(interface{}) error, arg ActionArg) {}

// UnaryIface executes an action which takes one argument, and returns a tuple of (interface{}, error).
//
// The single interface{} argument accepted by the action can be supplied via the arg parameter.
//
// The context will apply the supplied action's interface{} result to the next action called within the context,
// unless the override behavior for the next action dictates otherwise.
//
// The error returned by the supplied action is also applied to the current context.
// If error is not nil, subsequent actions executed within the same context will be ignored.
func (d *DeferredErrorContext) UnaryIface(action func(interface{}) (interface{}, error), arg ActionArg) {
}

// NullaryBool executes an action which takes no arguments and returns a tuple of (bool, error).
//
// The context will apply the supplied action's bool result to the next action called within the context,
// unless the override behavior for the next action dictates otherwise.
//
// The error returned by the supplied action is also applied to the current context.
// If error is not nil, subsequent actions executed within the same context will be ignored.
//
// In addition to threading the (bool, error) tuple into the current context, NullaryBool itself also returns a bool.
// This is useful for inlining the method in boolean statements.
func (d *DeferredErrorContext) NullaryBool(action func() (bool, error)) bool { return false }

// UnaryBool executes an action which takes one argument and returns a tuple of (bool, error).
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
func (d *DeferredErrorContext) UnaryBool(action func(interface{}) (bool, error), arg ActionArg) bool {
	return false
}
