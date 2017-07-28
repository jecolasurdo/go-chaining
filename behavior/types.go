// Package behavior provides mechanism for controlling how values are threaded through the execution chain.
package behavior

// InjectionOption instructs the system on how to thread values through the chain of actions.
type InjectionOption int

const (
	// NotSpecified means that no behavior has been declared. This value is assumed to mean UsePrevious
	NotSpecified InjectionOption = iota

	// InjectPreviousResult means that your intention is to inject the value supplied by the previous action in the chain.
	// If no previous action was present in the chain, nil is injected.
	// If a value is supplied in the Argument to the chain function, it is ignored.
	InjectPreviousResult

	// InjectSuppliedValue means that your intention is to inject the value supplied in the Argument.
	// If a previous value is presnet in the chain, it is ignored.
	InjectSuppliedValue
)
