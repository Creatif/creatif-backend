package main

import (
	"github.com/schollz/progressbar/v3"
)

func generateProgressBar(num int64) chan bool {
	bar := progressbar.Default(num)
	notifier := make(chan bool)

	go func() {
		for {
			bar.Add(1)
			<-notifier
		}

	}()

	return notifier
}
