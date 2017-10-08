package converter

import (
	"canihazmidi/src/prolink2midi/converter/midi"
	"fmt"
	"go.evanpurkhiser.com/prolink"
	"log"
	"os"
)

const midi_out_dev = "/dev/snd/midiC1D0"

type Converter struct {
	midi *midi.Clock
	net  *prolink.Network
	mon  *prolink.CDJStatusMonitor
}

func NewConverter() *Converter {
	c := new(Converter)
	dev := openMidiDev(midi_out_dev)
	c.midi = midi.NewClock(dev)
	log.Printf("Opened midi device clock: %v", midi_out_dev)

	net, err := prolink.Connect(prolink.Config{VirtualCDJID: 0x04})
	if err != nil {
		panic(err)
	}
	c.net = net
	log.Printf("Network initialized...")
	c.mon = net.CDJStatusMonitor()

	c.mon.OnStatusUpdate(prolink.StatusHandlerFunc(c.statusChange))

	return c
}

func openMidiDev(devicePath string) *os.File {
	m, err := os.OpenFile(devicePath, os.O_WRONLY, 0664)
	if err != nil {
		log.Fatal(err)
	}

	return m
}

func (c *Converter) statusChange(status *prolink.CDJStatus) {
	fmt.Println(status.String())

}
