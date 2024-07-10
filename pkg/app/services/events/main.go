package events

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/lib/constants"
	"creatif/pkg/lib/storage"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func DispatchEvent(dispachableEvent DispachableEvent) {
	evn := app.NewEvent(dispachableEvent.Project(), dispachableEvent.Type(), dispachableEvent.Data())
	// TODO: log failure somewhere
	storage.Gorm().Create(&evn)
}

func checkEvents() {
	offset := 0
	limit := 1000
	currentLen := -1
	var events []app.Event
	for currentLen != 0 {
		sql := fmt.Sprintf("SELECT id, type, data FROM %s OFFSET ? LIMIT ?", (app.Event{}).TableName())
		if res := storage.Gorm().Raw(sql, offset, limit).Scan(&events); res.Error != nil {
			fmt.Errorf("Event system fail: %w\n", res.Error)
		}

		for _, evn := range events {
			if evn.Type == constants.FileNotRemovedEvent {
				var realEvent FileNotRemovedEvent
				// TODO: log failure to future logging system
				if err := json.Unmarshal(evn.Data, &realEvent); err != nil {
					fmt.Errorf("Event system fail: %w\n", err)
				}

				// if the path does not exist, its a false event, remove the event
				_, err := os.Stat(realEvent.FilePath)
				if os.IsNotExist(err) {
					if res := storage.Gorm().Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", (app.Event{}).TableName()), evn.ID); res.Error != nil {
						fmt.Errorf("Event system fail: %w\n", res.Error)
					}

					continue
				}

				if err != nil {
					// TODO: Unexpected error, log to somewhere
				}

				// TODO: log failure to future logging system
				// remove the file associated with event
				if err := os.Remove(realEvent.FilePath); err != nil {
					fmt.Errorf("Event system fail: %w\n", err)
					continue
				}

				// remove the event itself
				if res := storage.Gorm().Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", (app.Event{}).TableName()), evn.ID); res.Error != nil {
					fmt.Errorf("Event system fail: %w\n", res.Error)
				}
			}

			if evn.Type == constants.PublicDirectoryNotRemovedEvent {
				var realEvent FileNotRemovedEvent
				// TODO: log failure to future logging system
				if err := json.Unmarshal(evn.Data, &realEvent); err != nil {
					fmt.Errorf("Event system fail: %w\n", err)
				}

				// TODO: log failure to future logging system
				if err := os.RemoveAll(realEvent.FilePath); err != nil {
					fmt.Errorf("Event system fail: %w\n", err)
					continue
				}
			}
		}

		offset += limit
		currentLen = len(events)
		events = make([]app.Event, 0)
	}
}

func RunEvents() {
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				checkEvents()
			}
		}
	}()
}
