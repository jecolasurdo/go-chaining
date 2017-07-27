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

func Test_ChainF_NoPreviousError_ForAnySpecifiedBehavior_SetsPreviousActionResult(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	expectedReturnValue := "expectedReturnValue"
	action := func(value interface{}) (interface{}, error) { return expectedReturnValue, nil }
	arg := errorhandler.ActionArg{
		Value: "valueFromArg",
	}

	d.PreviousActionResult = nil
	arg.Behavior = injectionbehavior.InjectSuppliedValue
	d.ChainF(action, arg)
	assert.Equal(t, expectedReturnValue, d.PreviousActionResult)

	d.PreviousActionResult = nil
	arg.Behavior = injectionbehavior.InjectPreviousResult
	d.ChainF(action, arg)
	assert.Equal(t, expectedReturnValue, d.PreviousActionResult)

	d.PreviousActionResult = nil
	arg.Behavior = injectionbehavior.NotSpecified
	d.ChainF(action, arg)
	assert.Equal(t, expectedReturnValue, d.PreviousActionResult)
}

func Test_Flush_Normally_ResetsContext(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	d.LocalError = errors.New("test error")
	d.PreviousActionResult = "Not nil"

	d.Flush()

	assert.Nil(t, d.LocalError)
	assert.Nil(t, d.PreviousActionResult)
}

func Test__Normally_ReturnsErrorAndFinalResult(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	d.LocalError = errors.New("test error")
	expectedFinalResult := "FinalResult"
	d.PreviousActionResult = expectedFinalResult

	result, err := d.Flush()

	assert.NotNil(t, err)
	assert.Equal(t, expectedFinalResult, result)
}
