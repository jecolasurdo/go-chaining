package errorhandler

// DeferredErrorContext is a mechanism for deferring execution of methods if an error condition has been received.
type DeferredErrorContext struct {
	LocalError error
}

// TryBool executes a func() (bool, error) if no previous Trys have resulted in an error.
// If previous Trys have resulted in an error, the action is ignored, not executed, and false is returned.
// Because false is returned when an action is ignored (rather than halting execution), it is important
// to ensure any downstream methods are also wrapped in Try methods, so they are also ignored.
func (d *DeferredErrorContext) TryBool(action func() (bool, error)) bool {
	if d.LocalError == nil {
		result, err := action()
		d.LocalError = err
		return result
	}
	return false
}

// TryVoid executes a func() error action if no previous Trys have resulted in a error.
// If previous Trys have resulted in an error, the action is ignored and not executed.
func (d *DeferredErrorContext) TryVoid(action func() error) {
	if d.LocalError == nil {
		d.LocalError = action()
	}
}

// FlushError returns the current local error and resets the error back to nil.
func (d *DeferredErrorContext) FlushError() error {
	localError := d.LocalError
	d.LocalError = nil
	return localError
}
