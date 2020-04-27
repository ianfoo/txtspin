package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/ianfoo/txtspin"
)

var (
	baseTaskLen uint
	speed       uint
	cycles      uint
	framestr    string
	framesep    string
)

func main() {
	flag.UintVar(&baseTaskLen, "tasklen", 750, "base duration of tasks in milliseconds")
	flag.UintVar(&speed, "speed", 0, "time in milliseconds between frames")
	flag.UintVar(&cycles, "cycles", 1, "number of times cycle should repeat (0=forever)")
	flag.StringVar(&framestr, "frames", "", "string of frames in animation sequence")
	flag.StringVar(&framesep, "framesep", "", "separator between frames in frame string")
	flag.Parse()

	if err := run(); err != nil {
		log.SetPrefix("error: ")
		log.SetFlags(0)
		log.Fatal(err)
	}
}

func run() error {
	cont := func() bool { return true }
	if cycles > 0 {
		var cycleCount uint
		cont = func() bool {
			defer func() { cycleCount++ }()
			return cycleCount < cycles
		}
	}
	var (
		spinners = getSpinners()
		messages = []string{
			"doing a long thing",
			"doing the next operation",
			"third time's a charm",
			"part four, should be finished after this, the longest part",
			"no, wait, there's one more thing that's longer, chump, lol",
		}
	)
	for cont() {
		for i, msg := range messages {
			task := func() {
				dur := time.Duration(baseTaskLen*uint(i+1)) * time.Millisecond
				time.Sleep(dur)
			}
			spinner := spinners[i%len(spinners)]
			spinner(task, msg)
		}
	}
	return nil
}

func getSpinners() []txtspin.Spinner {
	frameSpeed := time.Duration(speed) * time.Millisecond
	if frameSpeed == 0 {
		frameSpeed = txtspin.SpeedDefault
	}
	if framestr != "" {
		frames := strings.Split(framestr, framesep)
		return []txtspin.Spinner{txtspin.NewCustom(frames, frameSpeed)}
	}
	return []txtspin.Spinner{
		txtspin.New(txtspin.StyleSpinner, frameSpeed),
		txtspin.New(txtspin.StyleChompy, frameSpeed),
		txtspin.New(txtspin.StyleOoh, frameSpeed),

		// I think Flow and PingPong look a little better faster than
		// the single-char-width spinners, so scale frameSpeed to make
		// them run a bit faster.
		txtspin.New(txtspin.StyleFlow, frameSpeed/2),
		txtspin.New(txtspin.StylePingPong, frameSpeed/2),
	}
}

func mkWait(d time.Duration) func() error {
	return func() error {
		time.Sleep(d)
		return nil
	}
}
