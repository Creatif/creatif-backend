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
	evn := app.NewEvent(dispachableEvent.Type(), dispachableEvent.Data())
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
			fmt.Errorf("Event system fail: %w", res.Error)
		}

		fmt.Printf("Found %d events to be processed", len(events))
		for _, evn := range events {
			if evn.Type == constants.FileNotRemovedEvent {
				var realEvent FileNotRemovedEvent
				// TODO: log failure to future logging system
				if err := json.Unmarshal(evn.Data, &realEvent); err != nil {
					fmt.Errorf("Event system fail: %w", err)
				}

				fmt.Println("Unmarshaled FileNotRemovedEvent: ", string(evn.Data))

				// TODO: log failure to future logging system
				if err := os.Remove(realEvent.FilePath); err != nil {
					fmt.Errorf("Event system fail: %w", err)
				}

				if res := storage.Gorm().Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", (app.Event{}).TableName()), evn.ID); res.Error != nil {
					fmt.Errorf("Event system fail: %w", res.Error)
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
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("Execute event system")
				checkEvents()
			}
		}
	}()
}
