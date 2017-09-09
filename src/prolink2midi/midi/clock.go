// Provides a midi master clock and various other utils for working with midi
package midi

import (
	"os"
	"time"
)

const (
	START    = 0xfa
	STOP     = 0xfc
	TICK     = 0xf8
	CONTINUE = 0xfb
)

type Clock struct {
	cmd       chan []byte
	pulseRate chan time.Duration
	dev       *os.File
}

// Create a new midi clock. The clock starts to send tick events as soon as it is created.
//dev, err := os.OpenFile("/dev/snd/midiC1D0", os.O_WRONLY, 0664)
//if err != nil {
//	log.Fatal(err)
//}
//
//clk := midi.NewClock(dev)
//clk.SetBpm(120.00)
//clk.Start()
func NewClock(midiDevice *os.File) *Clock {
	clk := new(Clock)
	clk.dev = midiDevice
	clk.cmd = make(chan []byte)
	clk.pulseRate = make(chan time.Duration)

	go clk.run()

	return clk
}

// Change the BPM of the clock
func (clk *Clock) SetBpm(bpm float64) {
	clk.pulseRate <- bpmToPulseInterval(bpm)
}

// Send MIDI sequencer start event
func (clk *Clock) Start() {
	clk.cmd <- []byte{START}
}

// Send MIDI sequencer stop event
func (clk *Clock) Stop() {
	clk.cmd <- []byte{STOP}
}

// Send MIDI sequencer stop event
func (clk *Clock) Continue() {
	clk.cmd <- []byte{CONTINUE}
}

func (clk *Clock) run() {
	pulseRate := bpmToPulseInterval(120)
	tick := []byte{TICK}
	var t VarTicker
	t.SetDuration(pulseRate * time.Microsecond)
	go func() {
		for range t.C {
			clk.dev.Write(tick)
		}
	}()

	for {
		select {
		case newPulseRate := <-clk.pulseRate:
			t.SetDuration(newPulseRate * time.Microsecond)
		case cmd := <-clk.cmd:
			clk.dev.Write(cmd)
		}
	}

}

func bpmToPulseInterval(bpm float64) time.Duration {
	return time.Duration((6000000 / (bpm / 10)) / 24)
}
