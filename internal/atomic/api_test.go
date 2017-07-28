package atomic_test

import (
	"errors"
	chaining "jecolasurdo/go-chaining"
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

	arg := chaining.ActionArg{}
	c.LocalError = errors.New("test error")
	a.ApplyUnaryIface(action, &arg, c)

	assert.Equal(t, 0, timesActionWasCalled)
}

// func Test_ApplyUnaryIface_NoPreviousError_ExecutesAction(t *testing.T) {
// 	d := new(chaining.Context)
// 	timesActionWasCalled := 0
// 	action := func(interface{}) (interface{}, error) {
// 		timesActionWasCalled++
// 		return nil, nil
// 	}

// 	d.ApplyUnaryIface(action, chaining.ActionArg{})

// 	assert.Equal(t, 1, timesActionWasCalled)
// }

// func Test_ApplyUnaryIface_NoPreviousError_BehaviorIsNotSpecified_InjectsPreviousValue(t *testing.T) {
// 	d := new(chaining.Context)
// 	injectedValue := ""
// 	action := func(value interface{}) (interface{}, error) {
// 		injectedValue = value.(string)
// 		return nil, nil
// 	}
// 	argWithBehaviorNotSpecified := chaining.ActionArg{}
// 	simulatedValueOfPreviousActionInChain := "somevalue"
// 	d.PreviousActionResult = simulatedValueOfPreviousActionInChain

// 	d.ApplyUnaryIface(action, argWithBehaviorNotSpecified)

// 	assert.Equal(t, simulatedValueOfPreviousActionInChain, injectedValue)
// }

// func Test_ApplyUnaryIface_NoPreviousError_BehaviorIsUsePrevious_InjectsPreviousValue(t *testing.T) {
// 	d := new(chaining.Context)
// 	injectedValue := ""
// 	action := func(value interface{}) (interface{}, error) {
// 		injectedValue = value.(string)
// 		return nil, nil
// 	}
// 	argWithSpecifiedBehavior := chaining.ActionArg{
// 		Behavior: injectionbehavior.InjectPreviousResult,
// 	}
// 	simulatedValueOfPreviousActionInChain := "somevalue"
// 	d.PreviousActionResult = simulatedValueOfPreviousActionInChain

// 	d.ApplyUnaryIface(action, argWithSpecifiedBehavior)

// 	assert.Equal(t, simulatedValueOfPreviousActionInChain, injectedValue)
// }

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