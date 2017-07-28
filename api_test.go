package chaining_test

import (
	"errors"
	"testing"

	"jecolasurdo/go-chaining"
	"jecolasurdo/go-chaining/injectionbehavior"

	"github.com/stretchr/testify/assert"
)

func Test_ApplyNullaryBool_ActionTrueNoError_ReturnsTrue(t *testing.T) {
	d := new(chaining.Context)
	trueAction := func() (bool, error) { return true, nil }
	result := d.ApplyNullaryBool(trueAction, injectionbehavior.NotSpecified)
	assert.True(t, result)
}

func Test_ApplyNullaryBool_ActionFalseNoError_ReturnsFalse(t *testing.T) {
	d := new(chaining.Context)
	falseAction := func() (bool, error) { return false, nil }
	result := d.ApplyNullaryBool(falseAction, injectionbehavior.NotSpecified)
	assert.False(t, result)
}

func Test_ApplyNullaryBool_PreviousError_IgnoresAction(t *testing.T) {
	d := new(chaining.Context)
	timesActionWasCalled := 0
	action := func() (bool, error) {
		timesActionWasCalled++
		return true, nil
	}

	d.LocalError = errors.New("test error")
	d.ApplyNullaryBool(action, injectionbehavior.NotSpecified)

	assert.Equal(t, 0, timesActionWasCalled)
}

func Test_ApplyNullaryBool_PreviousError_ReturnsFalse(t *testing.T) {
	d := new(chaining.Context)
	trueAction := func() (bool, error) {
		return true, nil
	}

	d.LocalError = errors.New("test error")
	result := d.ApplyNullaryBool(trueAction, injectionbehavior.NotSpecified)

	assert.False(t, result)
}

func Test_ApplyNullaryBool_ActionErrors_SetsLocalError(t *testing.T) {
	d := new(chaining.Context)
	errorAction := func() (bool, error) { return false, errors.New("test error") }

	d.ApplyNullaryBool(errorAction, injectionbehavior.NotSpecified)

	assert.NotNil(t, d.LocalError)
}

func Test_ApplyNullary_NoPreviousError_ExecutesAction(t *testing.T) {
	d := new(chaining.Context)
	timesActionWasCalled := 0
	action := func() error {
		timesActionWasCalled++
		return nil
	}

	d.ApplyNullary(action, injectionbehavior.NotSpecified)

	assert.Equal(t, 1, timesActionWasCalled)
}

func Test_ApplyNullary_PreviousError_IgnoresAction(t *testing.T) {
	d := new(chaining.Context)
	timesActionWasCalled := 0
	action := func() error {
		timesActionWasCalled++
		return nil
	}

	d.LocalError = errors.New("test error")
	d.ApplyNullary(action, injectionbehavior.NotSpecified)

	assert.Equal(t, 0, timesActionWasCalled)
}

func Test_ApplyNullary_ActionErrors_SetsLocalError(t *testing.T) {
	d := new(chaining.Context)
	errorAction := func() error { return errors.New("test error") }

	d.ApplyNullary(errorAction, injectionbehavior.NotSpecified)

	assert.NotNil(t, d.LocalError)
}

func Test_ApplyUnaryIface_PreviousError_IgnoresAction(t *testing.T) {
	d := new(chaining.Context)
	timesActionWasCalled := 0
	action := func(interface{}) (interface{}, error) {
		timesActionWasCalled++
		return nil, nil
	}

	d.LocalError = errors.New("test error")
	d.ApplyUnaryIface(action, chaining.ActionArg{})

	assert.Equal(t, 0, timesActionWasCalled)
}

func Test_ApplyUnaryIface_NoPreviousError_ExecutesAction(t *testing.T) {
	d := new(chaining.Context)
	timesActionWasCalled := 0
	action := func(interface{}) (interface{}, error) {
		timesActionWasCalled++
		return nil, nil
	}

	d.ApplyUnaryIface(action, chaining.ActionArg{})

	assert.Equal(t, 1, timesActionWasCalled)
}

func Test_ApplyUnaryIface_NoPreviousError_BehaviorIsNotSpecified_InjectsPreviousValue(t *testing.T) {
	d := new(chaining.Context)
	var injectedValue interface{}
	action := func(value interface{}) (interface{}, error) {
		injectedValue = value
		return nil, nil
	}
	argWithBehaviorNotSpecified := chaining.ActionArg{}
	var simulatedValueOfPreviousActionInChain interface{} = "somevalue"
	d.PreviousActionResult = &simulatedValueOfPreviousActionInChain

	d.ApplyUnaryIface(action, argWithBehaviorNotSpecified)

	assert.Equal(t, &simulatedValueOfPreviousActionInChain, injectedValue)
}

func Test_ApplyUnaryIface_NoPreviousError_BehaviorIsUsePrevious_InjectsPreviousValue(t *testing.T) {
	d := new(chaining.Context)
	var injectedValue interface{}
	action := func(value interface{}) (interface{}, error) {
		injectedValue = value
		return nil, nil
	}
	argWithSpecifiedBehavior := chaining.ActionArg{
		Behavior: injectionbehavior.InjectPreviousResult,
	}
	var simulatedValueOfPreviousActionInChain interface{} = "somevalue"
	d.PreviousActionResult = &simulatedValueOfPreviousActionInChain

	d.ApplyUnaryIface(action, argWithSpecifiedBehavior)

	assert.Equal(t, &simulatedValueOfPreviousActionInChain, injectedValue)
}

func Test_ApplyUnaryIface_NoPreviousError_BehaviorIsOverridePrevious_InjectsSuppliedValue(t *testing.T) {
	d := new(chaining.Context)
	injectedValue := ""
	action := func(value interface{}) (interface{}, error) {
		injectedValue = value.(string)
		return nil, nil
	}
	valueSubmittedThroughArg := "valueFromArg"
	argWithSpecifiedBehavior := chaining.ActionArg{
		Behavior: injectionbehavior.InjectSuppliedValue,
		Value:    valueSubmittedThroughArg,
	}
	var simulatedValueOfPreviousActionInChain interface{} = "previousValue"
	d.PreviousActionResult = &simulatedValueOfPreviousActionInChain

	d.ApplyUnaryIface(action, argWithSpecifiedBehavior)

	assert.Equal(t, valueSubmittedThroughArg, injectedValue)
}

func Test_ApplyUnaryIface_NoPreviousError_ForAnySpecifiedBehavior_SetsPreviousActionResult(t *testing.T) {
	d := new(chaining.Context)
	var expectedReturnValue interface{} = "expectedReturnValue"
	action := func(value interface{}) (interface{}, error) { return expectedReturnValue, nil }
	arg := chaining.ActionArg{
		Value: "valueFromArg",
	}

	d.PreviousActionResult = nil
	arg.Behavior = injectionbehavior.InjectSuppliedValue
	d.ApplyUnaryIface(action, arg)
	assert.Equal(t, &expectedReturnValue, d.PreviousActionResult)

	d.PreviousActionResult = nil
	arg.Behavior = injectionbehavior.InjectPreviousResult
	d.ApplyUnaryIface(action, arg)
	assert.Equal(t, &expectedReturnValue, d.PreviousActionResult)

	d.PreviousActionResult = nil
	arg.Behavior = injectionbehavior.NotSpecified
	d.ApplyUnaryIface(action, arg)
	assert.Equal(t, &expectedReturnValue, d.PreviousActionResult)
}

func Test_Flush_Normally_ResetsContext(t *testing.T) {
	d := new(chaining.Context)
	d.LocalError = errors.New("test error")

	var notNilValue interface{} = "Not nil"
	d.PreviousActionResult = &notNilValue

	d.Flush()

	assert.Nil(t, d.LocalError)
	assert.Nil(t, d.PreviousActionResult)
}

func Test_Flush_Normally_ReturnsErrorAndFinalResult(t *testing.T) {
	d := new(chaining.Context)
	d.LocalError = errors.New("test error")
	var expectedFinalResult interface{} = "FinalResult"
	d.PreviousActionResult = &expectedFinalResult

	result, err := d.Flush()

	assert.NotNil(t, err)
	assert.Equal(t, &expectedFinalResult, result)
}
