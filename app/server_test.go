package app

import (
	"time"
)

func ExampleWorkerContainer() {
	var wc WorkerContainer
	wc.GFork(func(wc WorkerContainer) {
		for {
			select {
			case <-time.After(3 * time.Second):
				// select other channel, do something to get error.
				if err := error(nil); err != nil {
					// when got none-recoverable error, notify container to quit.
					wc.Quit()
					return
				}
			case <-wc.QC():
				// when got a quit signal, break the loop.
				// and must notify the container again for other workers
				// in container to quit.
				wc.Quit()
				return
			}
		}
	})
}
