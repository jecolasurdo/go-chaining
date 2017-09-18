This package was a naive experiment from when I was in the early stages of learning Go.
I have since learned more effective error handling strategies, and don't recommend the use of this package. I keep it here for posterity, however.

# go-chaining

go-chaining is a small go package that allows you to defer error handling until the end of a chain of methods.
The package also allows you to cury the output of one method to the input of the next method in the chain.

## Documentation

Package documentation can be read at https://godoc.org/github.com/jecolasurdo/go-chaining

## Background

This package came about as an excercise to see if I could develop an API to streamline error handling in a chain of methods.

The goal was to create something that would execute a series of methods until one of the methods returned an error, at which point the remainder of the methods in the chain would not be executed and the error could be returned.

In order to setup error curying, it also made sense to setup a mechanism for curying the result of one method into the input of the next method.

In the end, while the package works as designed, I find two things make it a bit unwieldy to use.
 - It relies heavily on the empty interface (`interface{}`) to handle argument currying in a somewhat generic fashion.
   - I don't really like that you lose compile time type checking with this package, but such is life without true generics.
 - Because Go doesn't support method overloading, I had to deal with some degree of method expansion.
   - Due to this, giving the methods names that are descriptive while also not too wordy or long was extremely difficult.
   - A constant battle began in trying to find a good balance between method name clarity and the expressiveness of the code consuming the methods.
   - In the end, I don't think I got either of those things to a good point
