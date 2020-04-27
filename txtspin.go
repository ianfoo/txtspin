// Package txtspin contains a builder for plain spinner-enhanced status updates.
package txtspin

import (
	"fmt"
	"strings"
	"time"
)

// SpeedDefault is the default delay between frames of the animation.
const SpeedDefault = 100 * time.Millisecond

// There are a few predefined styles that can be chosen.
const (
	StyleSpinner = "| / - \\"
	StyleChompy  = "- < C <"
	StyleOoh     = ". o O o"
	StyleFlow    = "-------- " +
		">------- " +
		">>------ " +
		">>>----- " +
		"->>>---- " +
		"-->>>--- " +
		"--->>>-- " +
		"---->>>- " +
		"----->>> " +
		"------>> " +
		"-------> " +
		"--------"
	StylePingPong = "-------- " +
		">------- " +
		">>------ " +
		">>>----- " +
		"->>>---- " +
		"-->>>--- " +
		"--->>>-- " +
		"---->>>- " +
		"----->>> " +
		"------>> " +
		"-------> " +
		"-------< " +
		"-------- " +
		"------<< " +
		"-----<<< " +
		"----<<<- " +
		"---<<<-- " +
		"--<<<--- " +
		"-<<<---- " +
		"<<<----- " +
		"<<------ " +
		"<-------"

	StyleDefault = StyleSpinner
)

// Spinner is a function that will run f and display a spinner as it runs, as
// configured by New or NewCustom.
type Spinner func(f func(), msg string)

// New returns a Spinner with the animation frames passed as a space-separated
// string.  If framesStr is empty, a default spinner animation will be used. If
// speed is 0, a default speed will be used.
func New(framesStr string, speed time.Duration) Spinner {
	if framesStr == "" {
		framesStr = StyleDefault
	}
	return NewCustom(strings.Split(framesStr, " "), speed)
}

// NewCustom returns a Spinner using the given frames and animation speed.
func NewCustom(frames []string, speed time.Duration) Spinner {
	if len(frames) == 0 {
		frames = strings.Split(StyleDefault, " ")
	}
	if speed == 0 {
		speed = SpeedDefault
	}

	// animate implements the spinning itself.
	animate := func(sig chan struct{}) {
		var (
			clearFrames = make([]string, len(frames))
			i           int
		)
		for i, f := range frames {
			clearFrames[i] = strings.Repeat("\b", len(f))
		}
		defer func() {
			var (
				clearIdx   = (i - 1 + len(frames)) % len(frames)
				clearFrame = clearFrames[clearIdx]
				blankFrame = strings.Repeat(" ", len(clearFrame))
			)
			fmt.Printf("%s%s\n", clearFrame, blankFrame)
			sig <- struct{}{}
		}()
		fmt.Print(strings.Repeat(" ", len(frames[0])))
		for {
			select {
			case <-sig:
				return
			default:
				fmt.Print(clearFrames[i])
				i = (i + 1) % len(frames)
				fmt.Print(frames[i])
				time.Sleep(time.Duration(speed))
			}
		}
	}

	// spinner is returned to the user as a function to call to write a
	// message to stdout, run a function, and display a spinner while the
	// function is running.
	spinner := func(f func(), msg string) {
		sig := make(chan struct{})
		defer close(sig)
		if msg != "" {
			fmt.Print(msg, " ")
		}
		go animate(sig)
		f()
		sig <- struct{}{}
		<-sig
	}
	return spinner
}
