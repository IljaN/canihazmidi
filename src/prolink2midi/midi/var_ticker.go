package midi

import (
	"time"
)

// Variable frequency ticker from https://groups.google.com/d/msg/golang-nuts/TMX7D6XrTHk/Ode7de5KBgAJ
type VarTicker struct {
	C    <-chan time.Time
	ch   chan<- time.Time
	t    *time.Ticker
	done chan bool
}

func (t *VarTicker) SetDuration(d time.Duration) {
	if t.t != nil {
		t.t.Stop()
		close(t.done)
	} else {
		var ticker = make(chan time.Time)
		t.C = ticker
		t.ch = ticker
	}
	t.done = make(chan bool)
	t.t = time.NewTicker(d)
	go func(out chan<- time.Time, in <-chan time.Time, done <-chan bool) {
		for {
			select {
			case tick := <-in:
				out <- tick
			case <-done:
				return
			}
		}
	}(t.ch, t.t.C, t.done)
}
