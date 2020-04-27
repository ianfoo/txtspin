# txtspin

A little Go package that adds a spinner and an optional message to
tasks that may take a minute to run.

A number of default Spinners are defined, or you can create one with a slice of
strings as the frames of the animation. Run `example/example.go` to see a quick
demo of the defined styles, or experiment with different patterns and speeds.

[![Asciicast demo: click to visit](demo.gif)](https://asciinema.org/a/xWh9tMnX0sXwhNFcyRq0nqMyW)


### Tip

For best flexibility, the `Spinner` function takes a function that takes takes
no parameters and returns nothing. If the operation you want to call takes
arguments, or returns values, just wrap it in a small function that closes over
values that you can inspect after the operation has finished.
```
var (
	err    error
	result Result
	f             = func() { result, err = MyPossiblyErroringOperation(arg) }
	spinner       = txtspin.New(txtspin.StyleDefault)
)
spinner(f, "running slow op!")

```


### Caveats

This should only be used if it's known that the output is bound for a TTY.
You can use [go-isatty](https://github.com/mattn/go-isatty) to determine
this.

This was written for fun late one night and carries the commensurate
guarantees of stability and well-testedness.
