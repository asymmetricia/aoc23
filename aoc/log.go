package aoc

import (
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var plogState struct {
	last time.Time
	init time.Time
	note string
}
var plogMu = &sync.Mutex{}

func ProgressLog(log *logrus.Entry, note string, i, total int) {
	plogMu.Lock()
	defer plogMu.Unlock()
	if note != plogState.note || plogState.init.IsZero() {
		plogState.init = time.Now()
		plogState.last = time.Time{}
		plogState.note = note
	}

	if time.Since(plogState.last) >= time.Second {
		perc := int(100 * float32(i) / float32(total))
		var eta float64
		if i > 0 {
			eta = time.Since(plogState.init).Seconds() / float64(i) * float64(total-i)
		}
		log.Printf("%d/%d (%d%%, ETA %.2fs)", i, total, perc, eta)
		plogState.last = time.Now()
	}
}
