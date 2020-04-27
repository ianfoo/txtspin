package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/ianfoo/txtspin"
)

var (
	tasklen  uint
	speed    uint
	cycles   uint
	style    string
	framestr string
	framesep string
)

func main() {
	flag.UintVar(&tasklen, "tasklen", 2000, "duration of tasks in milliseconds")
	flag.UintVar(&speed, "speed", 0, "time in milliseconds between frames")
	flag.UintVar(&cycles, "cycles", 1, "number of times cycle should repeat (0=forever)")
	flag.StringVar(&style, "style", "", "comma-separated names of styles to use for messages")
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
			"it'll work this time",
			"finishing up",
			"okay, we should be finished after this, the longest part",
			"no, wait, there's something else that takes longer, my bad",
		}
	)
	for cont() {
		for i, spinner := range spinners {
			var (
				task = func() { time.Sleep(time.Duration(tasklen) * time.Millisecond) }
				msg  = messages[i%len(messages)]
			)
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
	var (
		spinnerMap = map[string]txtspin.Spinner{
			"spinner": txtspin.New(txtspin.StyleSpinner, frameSpeed),
			"chompy":  txtspin.New(txtspin.StyleChompy, frameSpeed),
			"ooh":     txtspin.New(txtspin.StyleOoh, frameSpeed),
			"eyes":    txtspin.New(txtspin.StyleEyes, frameSpeed),
			"blink":   txtspin.New(txtspin.StyleBlink, frameSpeed),

			// I think Flow and PingPong look a little better faster than
			// the single-char-width spinners, so scale frameSpeed to make
			// them run a bit faster.
			"flow":     txtspin.New(txtspin.StyleFlow, frameSpeed/2),
			"pingpong": txtspin.New(txtspin.StylePingPong, frameSpeed/2),
		}
		spinners []txtspin.Spinner
		styles   = strings.Split(style, ",")
	)
	for _, s := range styles {
		if spinner, ok := spinnerMap[s]; ok {
			spinners = append(spinners, spinner)
		}
	}
	if len(spinners) > 0 {
		return spinners
	}
	all := []txtspin.Spinner{
		spinnerMap["spinner"],
		spinnerMap["chompy"],
		spinnerMap["ooh"],
		spinnerMap["eyes"],
		spinnerMap["blink"],
		spinnerMap["flow"],
		spinnerMap["pingpong"],
	}
	return all
}

func mkWait(d time.Duration) func() error {
	return func() error {
		time.Sleep(d)
		return nil
	}
}
