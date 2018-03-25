// Provides a midi master clock and various other utils for working with midi
package midi

import (
	"os"
	"time"
)

// Midi standard commands
const (
	StatusStart    = 0xfa
	StatusStop     = 0xfc
	StatusTick     = 0xf8
	StatusContinue = 0xfb
)

type Clock struct {
	cmd       chan []byte
	pulseRate chan time.Duration
	dev       *os.File
	Playing   bool
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
	clk.Playing = false

	go clk.run()
	clk.Stop()

	return clk
}

// Change the BPM of the clock
func (clk *Clock) SetBpm(bpm float32) {
	clk.pulseRate <- bpmToPulseInterval(bpm)
}

// Send MIDI sequencer start event
func (clk *Clock) Start() {
	clk.Playing = true
	clk.cmd <- []byte{StatusStart}
}

// Send MIDI sequencer stop event
func (clk *Clock) Stop() {
	clk.Playing = false
	clk.cmd <- []byte{StatusStop}
}


// Send MIDI sequencer stop event
func (clk *Clock) Continue() {
	clk.cmd <- []byte{StatusContinue}
}

func (clk *Clock) run() {
	pulseRate := bpmToPulseInterval(120)
	var t VarTicker
	t.SetDuration(pulseRate)

	go func() {
		for range t.C {
			clk.dev.Write([]byte{StatusTick})
		}
	}()

	for {
		select {
		case newPulseRate := <-clk.pulseRate:
			t.SetDuration(newPulseRate)
		case cmd := <-clk.cmd:
			clk.dev.Write(cmd)
		}
	}

}

// Pulses per quarter note
const ppqn = 24
const usec_in_min = 6000000

// Converts bpm to a 24ppqn pulse interval in microseconds
func bpmToPulseInterval(bpm float32) time.Duration {
	return time.Duration((usec_in_min/(bpm/10))/ppqn) * time.Microsecond
}
