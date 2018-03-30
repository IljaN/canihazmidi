package main

import (
	"canihazmidi/src/prolink2midi/midi"
	"go.evanpurkhiser.com/prolink"
)

type Sync struct {
	mon               *prolink.CDJStatusMonitor
	clk               *midi.Clock
	currentBPM        float32
	lastBeatInMeasure uint8
}

func newSynchronizer(mon *prolink.CDJStatusMonitor, clk *midi.Clock) *Sync {
	s := &Sync{mon: mon, clk: clk}
	s.clk.Stop()
	s.mon.OnStatusUpdate(prolink.StatusHandlerFunc(s.sync))

	return s
}

func (s *Sync) setBpm(bpm float32) {
	if s.currentBPM != bpm {
		s.clk.SetBpm(bpm)
		s.currentBPM = bpm
	}
}

func (s *Sync) sync(status *prolink.CDJStatus) {
	if !status.IsMaster {
		return
	}

	if status.PlayState == prolink.PlayStatePaused {
		if s.clk.Playing {
			s.clk.Stop()
		}
		s.setBpm(calcEffectiveBpm(status.TrackBPM, status.SliderPitch))

		return
	}

	if !s.clk.Playing && status.PlayState == prolink.PlayStatePlaying && status.BeatInMeasure == 1 && status.Beat > 0 {
		s.clk.Start()
	}

}

func calcEffectiveBpm(bpm float32, pitch float32) float32 {
	return bpm + ((bpm / 100.000) * pitch)
}
