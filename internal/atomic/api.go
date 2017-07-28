package atomic

import (
	chaining "jecolasurdo/go-chaining"
	"jecolasurdo/go-chaining/injectionbehavior"
)

// Context includes all fundamental elements from which higher order elements can be composed.
type Context struct{}

// ApplyUnaryIface is the atomic function from which all higher order chaining functions can be composed.
func (a *Context) ApplyUnaryIface(action func(interface{}) (interface{}, error), arg *chaining.ActionArg, c *chaining.Context) {
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
