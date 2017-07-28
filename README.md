# go-chaining
go-chaining is a small go package that allows you to defer error handling until the end of a chain of methods.

## Example

```go

// Some setup (pardon the compressed formatting)
type Foo struct { }
func (f *Foo) SomethingIsTrue() (bool error) { return true, nil }
func (f *Foo) SomethingElseIsTrue() (bool error) { return true, nil }
func (f *Foo) DoSomething() error { return nil }
func (f *Foo) MakeItRain() error { return nil }

// Nested if statement with standard go error handling.
func (f *Foo) Flar() error {
    somethingIsTrue, err:= f.SomethingIsTrue()
    if err != nil {
        return err
    }
    if somethingIsTrue {
        somethingElseIsTrue, err:= f.SomethingElseIsTrue()
        if err != nil {
            return err
        }
        if somethingElseIsTrue {
            err:= f.DoSomething()
            if err != nil {
                return err
            }
            err = f.MakeItRain()
            if err != nil {
                return error
            }
        }
    } else {
        err:= f.MakeItRain()
        if err != nil {
            return err
        }
    }
    return nil
}

// Same nested if statement with deferred error handling.
func (f *Foo) Flar() error {
	e := new(chaining.Context)
	if e.ExecNullaryBool(f.SomethingIsTrue) {
		if e.ExecNullaryBool(f.SomethingElseIsTrue) {
			e.ExecNullaryVoid(f.DoSomething)
			e.ExecNullaryVoid(f.MakeItRain)
		}
	} else {
		e.ExecNullaryVoid(f.MakeItRain)
	}
    _, err:= e.Flush()
	return err
}

```
