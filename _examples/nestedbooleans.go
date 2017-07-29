package examples

// import (
// 	chaining "jecolasurdo/go-chaining"
// 	"jecolasurdo/go-chaining/behavior"
// )

// type foo struct{}

// func (f *foo) somethingIsTrue() (bool, error)     { return true, nil }
// func (f *foo) somethingElseIsTrue() (bool, error) { return true, nil }
// func (f *foo) doSomething() error                 { return nil }
// func (f *foo) makeItRain() error                  { return nil }

// // Nested if statement with standard go error handling.
// func (f *foo) flar() error {
// 	somethingIsTrue, err := f.somethingIsTrue()
// 	if err != nil {
// 		return err
// 	}
// 	if somethingIsTrue {
// 		somethingElseIsTrue, err := f.somethingElseIsTrue()
// 		if err != nil {
// 			return err
// 		}
// 		if somethingElseIsTrue {
// 			err := f.doSomething()
// 			if err != nil {
// 				return err
// 			}
// 			err = f.makeItRain()
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	} else {
// 		err := f.makeItRain()
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// // Same nested if statements with deferred error handling.
// func (f *foo) chainedFlar() error {
// 	c := new(chaining.Context)
// 	if c.NBool(f.somethingIsTrue, behavior.NotSpecified) {
// 		if c.NBool(f.somethingElseIsTrue, behavior.NotSpecified) {
// 			c.N(f.doSomething, behavior.NotSpecified)
// 			c.N(f.makeItRain, behavior.NotSpecified)
// 		}
// 	} else {
// 		c.N(f.makeItRain, behavior.NotSpecified)
// 	}
// 	_, err := c.Flush()
// 	return err
// }
