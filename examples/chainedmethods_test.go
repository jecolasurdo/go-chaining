package examples

import (
	"fmt"
	c "jecolasurdo/go-chaining"
	"jecolasurdo/go-chaining/behavior"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func addOne(value *interface{}) (*interface{}, error) {
	var result interface{} = (*value).(int) + 1
	return &result, nil
}

func multiplyBySix(value *interface{}) (*interface{}, error) {
	var result interface{} = (*value).(int) * 6
	return &result, nil
}

func convertToString(value *interface{}) (*interface{}, error) {
	var result interface{} = strconv.Itoa((*value).(int))
	return &result, nil
}

func sendToThePrinter(value *interface{}) error {
	_, err := fmt.Print(value)
	return err
}

func Test_ChainItAllTogether(t *testing.T) {
	chain := c.New()
	var initialValue interface{} = 1
	initialBehavior := c.ActionArg{
		Value:    &initialValue,
		Behavior: behavior.InjectSuppliedValue,
	}
	chain.ApplyUnaryIface(addOne, initialBehavior)
	// chain.ApplyUnaryIface(multiplyBySix, c.ActionArg{})
	// chain.ApplyUnaryIface(convertToString, c.ActionArg{})
	// chain.ApplyUnary(sendToThePrinter, c.ActionArg{})

	result, err := chain.Flush()
	assert.Equal(t, 2, (*result).(int))
	assert.Nil(t, err)
}
