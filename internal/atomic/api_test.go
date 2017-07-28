package atomic_test

import (
	"errors"
	chaining "jecolasurdo/go-chaining"
	"jecolasurdo/go-chaining/injectionbehavior"
	"jecolasurdo/go-chaining/internal/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ApplyUnaryIface_PreviousError_IgnoresAction(t *testing.T) {
	a := new(atomic.Context)
	c := new(chaining.Context)
	timesActionWasCalled := 0
	action := func(interface{}) (interface{}, error) {
		timesActionWasCalled++
		return nil, nil
	}

	c.LocalError = errors.New("test error")
	a.ApplyUnaryIface(action, &chaining.ActionArg{}, c)

	assert.Equal(t, 0, timesActionWasCalled)
}

func Test_ApplyUnaryIface_NoPreviousError_ExecutesAction(t *testing.T) {
	a := new(atomic.Context)
	c := new(chaining.Context)
	timesActionWasCalled := 0
	action := func(interface{}) (interface{}, error) {
		timesActionWasCalled++
		return nil, nil
	}

	a.ApplyUnaryIface(action, &chaining.ActionArg{}, c)

	assert.Equal(t, 1, timesActionWasCalled)
}

func Test_ApplyUnaryIface_NoPreviousError_BehaviorIsNotSpecified_InjectsPreviousValue(t *testing.T) {
	a := new(atomic.Context)
	c := new(chaining.Context)
	injectedValue := ""
	action := func(value interface{}) (interface{}, error) {
		injectedValue = value.(string)
		return nil, nil
	}
	argWithBehaviorNotSpecified := chaining.ActionArg{}
	simulatedValueOfPreviousActionInChain := "somevalue"
	c.PreviousActionResult = simulatedValueOfPreviousActionInChain

	a.ApplyUnaryIface(action, &argWithBehaviorNotSpecified, c)

	assert.Equal(t, simulatedValueOfPreviousActionInChain, injectedValue)
}

func Test_ApplyUnaryIface_NoPreviousError_BehaviorIsUsePrevious_InjectsPreviousValue(t *testing.T) {
	a := new(atomic.Context)
	c := new(chaining.Context)
	injectedValue := ""
	action := func(value interface{}) (interface{}, error) {
		injectedValue = value.(string)
		return nil, nil
	}
	argWithSpecifiedBehavior := chaining.ActionArg{
		Behavior: injectionbehavior.InjectPreviousResult,
	}
	simulatedValueOfPreviousActionInChain := "somevalue"
	c.PreviousActionResult = simulatedValueOfPreviousActionInChain

	a.ApplyUnaryIface(action, &argWithSpecifiedBehavior, c)

	assert.Equal(t, simulatedValueOfPreviousActionInChain, injectedValue)
}

// func Test_ApplyUnaryIface_NoPreviousError_BehaviorIsOverridePrevious_InjectsSuppliedValue(t *testing.T) {
// 	d := new(chaining.Context)
// 	injectedValue := ""
// 	action := func(value interface{}) (interface{}, error) {
// 		injectedValue = value.(string)
// 		return nil, nil
// 	}
// 	valueSubmittedThroughArg := "valueFromArg"
// 	argWithSpecifiedBehavior := chaining.ActionArg{
// 		Behavior: injectionbehavior.InjectSuppliedValue,
// 		Value:    valueSubmittedThroughArg,
// 	}
// 	simulatedValueOfPreviousActionInChain := "previousValue"
// 	d.PreviousActionResult = simulatedValueOfPreviousActionInChain

// 	d.ApplyUnaryIface(action, argWithSpecifiedBehavior)

// 	assert.Equal(t, valueSubmittedThroughArg, injectedValue)
// }

// func Test_ApplyUnaryIface_NoPreviousError_ForAnySpecifiedBehavior_SetsPreviousActionResult(t *testing.T) {
// 	d := new(chaining.Context)
// 	expectedReturnValue := "expectedReturnValue"
// 	action := func(value interface{}) (interface{}, error) { return expectedReturnValue, nil }
// 	arg := chaining.ActionArg{
// 		Value: "valueFromArg",
// 	}

// 	d.PreviousActionResult = nil
// 	arg.Behavior = injectionbehavior.InjectSuppliedValue
// 	d.ApplyUnaryIface(action, arg)
// 	assert.Equal(t, expectedReturnValue, d.PreviousActionResult)

// 	d.PreviousActionResult = nil
// 	arg.Behavior = injectionbehavior.InjectPreviousResult
// 	d.ApplyUnaryIface(action, arg)
// 	assert.Equal(t, expectedReturnValue, d.PreviousActionResult)

// 	d.PreviousActionResult = nil
// 	arg.Behavior = injectionbehavior.NotSpecified
// 	d.ApplyUnaryIface(action, arg)
// 	assert.Equal(t, expectedReturnValue, d.PreviousActionResult)
// }
