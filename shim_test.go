package comshim_test

import (
	"runtime"
	"sync"
	"testing"

	"github.com/bi-zone/go-ole/oleutil"
	"github.com/scjalliance/comshim"
)

func TestConcurrentShims(t *testing.T) {
	var maxRounds int
	if testing.Short() {
		maxRounds = 64
	} else {
		maxRounds = 256
	}

	// Vary the number of threads
	for procs := 1; procs < 11; procs++ {
		runtime.GOMAXPROCS(procs)

		// Vary the number of shims
		for rounds := 1; rounds <= maxRounds; rounds *= 2 {
			wg := sync.WaitGroup{}
			for i := 0; i < rounds; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()

					comshim.Add(1)
					defer comshim.Done()

					obj, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
					if err != nil {
						t.Error(err)
					} else {
						defer obj.Release()
					}
				}(i)
			}
			wg.Wait()
		}
	}
}
