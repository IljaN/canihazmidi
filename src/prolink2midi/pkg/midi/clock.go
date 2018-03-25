// Provides a midi master clock and various other utils for working with midi
package midi

import (
	"os"
	"time"
)

// Midi standard commands
const (
	Start    = 0xFA
	Stop     = 0xFC
	Tick     = 0xF8
	Continue = 0xFB
)

// Pulses per quarter note
const ppqn = 24
const uSecInMin = 6000000

type Clock struct {
	midiOut   chan []byte
	pulseRate chan time.Duration
	device    *os.File
	Playing   bool
}

// Create a new midi clock. The clock starts to send tick events as soon as it is created.
// device, err := os.OpenFile("/device/snd/midiC1D0", os.O_WRONLY, 0664)
// if err != nil {
//	log.Fatal(err)
// }
//
// clk := midi.NewClock(device)
// clk.SetBpm(120.00)
// clk.Start()
func NewClock(device *os.File) *Clock {
	clk := &Clock{
		device:    device,
		midiOut:   make(chan []byte),
		pulseRate: make(chan time.Duration),
		Playing:   false,
	}

	go clk.run()

	return clk
}

// Change the BPM of the clock
func (clk *Clock) SetBpm(bpm float32) {
	clk.pulseRate <- bpmToPulseInterval(bpm)
}

// Send MIDI sequencer start event
func (clk *Clock) Start() {
	clk.midiOut <- []byte{Start}
	clk.Playing = true
}

// Send MIDI sequencer stop event
func (clk *Clock) Stop() {
	clk.midiOut <- []byte{Stop}
	clk.Playing = false
}

// Send MIDI sequencer stop event
func (clk *Clock) Continue() {
	clk.midiOut <- []byte{Continue}
	clk.Playing = true
}

func (clk *Clock) run() {
	pulseRate := bpmToPulseInterval(120)
	var t VarTicker
	t.SetDuration(pulseRate)

	go func() {
		for range t.C {
			clk.device.Write([]byte{Tick})
		}
	}()

	for {
		select {
		case newPulseRate := <-clk.pulseRate:
			t.SetDuration(newPulseRate)
		case cmd := <-clk.midiOut:
			clk.device.Write(cmd)
		}
	}

}

// Converts bpm to a 24ppqn pulse interval in microseconds
func bpmToPulseInterval(bpm float32) time.Duration {
	return time.Duration((uSecInMin/(bpm/10))/ppqn) * time.Microsecond
}
