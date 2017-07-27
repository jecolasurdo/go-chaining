package errorhandler_test

import (
	"errors"
	"jecolasurdo/go-deferrederrors/errorhandler/injectionbehavior"
	"testing"

	"jecolasurdo/go-deferrederrors/errorhandler"

	"github.com/stretchr/testify/assert"
)

func Test_TryBool_ActionTrueNoError_ReturnsTrue(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	trueAction := func() (bool, error) { return true, nil }
	result := d.TryBool(trueAction)
	assert.True(t, result)
}

func Test_TryBool_ActionFalseNoError_ReturnsFalse(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	falseAction := func() (bool, error) { return false, nil }
	result := d.TryBool(falseAction)
	assert.False(t, result)
}

func Test_TryBool_PreviousError_IgnoresAction(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	timesActionWasCalled := 0
	action := func() (bool, error) {
		timesActionWasCalled++
		return true, nil
	}

	d.LocalError = errors.New("test error")
	d.TryBool(action)

	assert.Equal(t, 0, timesActionWasCalled)
}

func Test_TryBool_PreviousError_ReturnsFalse(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	trueAction := func() (bool, error) {
		return true, nil
	}

	d.LocalError = errors.New("test error")
	result := d.TryBool(trueAction)

	assert.False(t, result)
}

func Test_TryBool_ActionErrors_SetsLocalError(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	errorAction := func() (bool, error) { return false, errors.New("test error") }

	d.TryBool(errorAction)

	assert.NotNil(t, d.LocalError)
}

func Test_TryVoid_NoPreviousError_ExecutesAction(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	timesActionWasCalled := 0
	action := func() error {
		timesActionWasCalled++
		return nil
	}

	d.TryVoid(action)

	assert.Equal(t, 1, timesActionWasCalled)
}

func Test_TryVoid_PreviousError_IgnoresAction(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	timesActionWasCalled := 0
	action := func() error {
		timesActionWasCalled++
		return nil
	}

	d.LocalError = errors.New("test error")
	d.TryVoid(action)

	assert.Equal(t, 0, timesActionWasCalled)
}

func Test_TryVoid_ActionErrors_SetsLocalError(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	errorAction := func() error { return errors.New("test error") }

	d.TryVoid(errorAction)

	assert.NotNil(t, d.LocalError)
}

func Test_ChainF_PreviousError_IgnoresAction(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	timesActionWasCalled := 0
	action := func(interface{}) (interface{}, error) {
		timesActionWasCalled++
		return nil, nil
	}

	d.LocalError = errors.New("test error")
	d.ChainF(action, errorhandler.ActionArg{})

	assert.Equal(t, 0, timesActionWasCalled)
}

func Test_ChainF_NoPreviousError_ExecutesAction(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	timesActionWasCalled := 0
	action := func(interface{}) (interface{}, error) {
		timesActionWasCalled++
		return nil, nil
	}

	d.ChainF(action, errorhandler.ActionArg{})

	assert.Equal(t, 1, timesActionWasCalled)
}

func Test_ChainF_NoPreviousError_BehaviorIsNotSpecified_InjectsPreviousValue(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	injectedValue := ""
	action := func(value interface{}) (interface{}, error) {
		injectedValue = value.(string)
		return nil, nil
	}
	argWithBehaviorNotSpecified := errorhandler.ActionArg{}
	simulatedValueOfPreviousActionInChain := "somevalue"
	d.PreviousActionResult = simulatedValueOfPreviousActionInChain

	d.ChainF(action, argWithBehaviorNotSpecified)

	assert.Equal(t, simulatedValueOfPreviousActionInChain, injectedValue)
}

func Test_ChainF_NoPreviousError_BehaviorIsUsePrevious_InjectsPreviousValue(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	injectedValue := ""
	action := func(value interface{}) (interface{}, error) {
		injectedValue = value.(string)
		return nil, nil
	}
	argWithSpecifiedBehavior := errorhandler.ActionArg{
		Behavior: injectionbehavior.InjectPreviousResult,
	}
	simulatedValueOfPreviousActionInChain := "somevalue"
	d.PreviousActionResult = simulatedValueOfPreviousActionInChain

	d.ChainF(action, argWithSpecifiedBehavior)

	assert.Equal(t, simulatedValueOfPreviousActionInChain, injectedValue)
}

func Test_ChainF_NoPreviousError_BehaviorIsOverridePrevious_InjectsSuppliedValue(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	injectedValue := ""
	action := func(value interface{}) (interface{}, error) {
		injectedValue = value.(string)
		return nil, nil
	}
	valueSubmittedThroughArg := "valueFromArg"
	argWithSpecifiedBehavior := errorhandler.ActionArg{
		Behavior: injectionbehavior.InjectSuppliedValue,
		Value:    valueSubmittedThroughArg,
	}
	simulatedValueOfPreviousActionInChain := "previousValue"
	d.PreviousActionResult = simulatedValueOfPreviousActionInChain

	d.ChainF(action, argWithSpecifiedBehavior)

	assert.Equal(t, valueSubmittedThroughArg, injectedValue)
}

// Test_ChainF_NoPreviousError_ForAnySpecifiedBehavior_SetsInjectionValueWithOutput

func Test_FlushChain_Normally_ResetsErrorToNil(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	d.LocalError = errors.New("test error")
	d.FlushChain()
	assert.Nil(t, d.LocalError)
}

func Test_FlushChain_Normally_ReturnsAnyEror(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	d.LocalError = errors.New("test error")
	err := d.FlushChain()
	assert.NotNil(t, err)
}

// ChainF(action func(interface{}) (interface{}, error), value interface{}, overridePrevious bool = false)
// How ChainF works:
//   value is passed into the closure as its argument.
//   The closure is executed and returns an interface and error.
//   Both the interface and error are cached in the context object
//   The next time ChainF is executed, the value from the previous method is used as the input for the next.
