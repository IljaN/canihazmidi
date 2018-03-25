package main

import (
	"canihazmidi/src/prolink2midi/midi"
	"fmt"
	"go.evanpurkhiser.com/prolink"
)

type sync struct {
	mon *prolink.CDJStatusMonitor
	clk *midi.Clock
}

func newSynchronizer(mon *prolink.CDJStatusMonitor, clk *midi.Clock) *sync {
	s := &sync{mon: mon, clk: clk}
	s.clk.Stop()
	s.mon.OnStatusUpdate(prolink.StatusHandlerFunc(s.sync))

	return s
}

func (s *sync) sync(status *prolink.CDJStatus) {
	fmt.Println(status.String())
	if !status.IsMaster {
		return
	}

	//s.clk.SetBpm(status.TrackBPM)

	if status.PlayState == prolink.PlayStatePlaying && status.BeatInMeasure == 1 && !s.clk.Playing {
		s.clk.Start()
		return
	}

	if status.PlayState == prolink.PlayStatePaused {
		s.clk.Stop()
	}
}
