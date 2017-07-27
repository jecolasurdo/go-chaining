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

// FlushChain returns the current local error and resets the context back to its default state.
func (d *DeferredErrorContext) FlushChain() error {
	localError := d.LocalError
	d.LocalError = nil
	d.PreviousActionResult = nil
	return localError
}

//TODO: Add new methods:
// FlushError gets renamed to Flush, resets the context returns a value also (always, even if nil)
// ChainBool(action func(interface{}) (bool, error), value interface{}, overridePrevious bool = false) bool
// ChainVoid(action func(interface{}) error, value interface{}, overridePrevious bool = false)
// ChainF(action func(interface{}) (interface{}, error), value interface{}, overridePrevious bool = false)
// How ChainF works:
//   value is passed into the closure as its argument.
//   The closure is executed and returns an interface and error.
//   Both the interface and error are cached in the context object
//   The next time ChainF is executed, the value from the previous method is used as the input for the next.
