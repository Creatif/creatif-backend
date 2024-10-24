package main

import (
	"github.com/schollz/progressbar/v3"
)

func generateProgressBar(num int) (chan bool, chan bool) {
	bar := progressbar.Default(int64(num))
	notifier := make(chan bool)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-notifier:
				bar.Add(1)
			case <-done:
				bar.Close()
			}
		}

	}()

	return notifier, done
}
