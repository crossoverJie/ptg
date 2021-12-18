package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/crossoverJie/ptg/gui/io"
)

type (
	History struct {
		lruCache          *LruCache
		notifyChan        chan struct{}
		historyButton     *fyne.Container
		alreadyButtonList []*widget.Button
		targetInput       *widget.Entry
		requestEntry      *widget.Entry
		metadataEntry     *widget.Entry
		responseEntry     *widget.Entry
	}

	HistoryValue struct {
		Id         int
		Value      *io.Log
		MethodInfo string
	}
)

func NewHistory(size int, historyButton *fyne.Container, targetInput, requestEntry, metadataEntry, responseEntry *widget.Entry) *History {
	h := &History{
		lruCache:      NewLruList(size),
		notifyChan:    make(chan struct{}, size),
		historyButton: historyButton,
		targetInput:   targetInput,
		requestEntry:  requestEntry,
		metadataEntry: metadataEntry,
		responseEntry: responseEntry,
	}
	go h.viewHistory()
	return h

}

func (h *History) Put(k, v interface{}) {
	h.lruCache.Put(k, v)
	h.notifyChan <- struct{}{}
}

func (h *History) viewHistory() {
	for {
		select {
		case <-h.notifyChan:

			// Reset view.
			for _, button := range h.alreadyButtonList {
				h.historyButton.Remove(button)
			}
			h.alreadyButtonList = make([]*widget.Button, 0)

			// Draw view.
			for _, v := range h.lruCache.List() {
				//index := i
				historyValue := v.(*HistoryValue)
				button := widget.NewButtonWithIcon(historyValue.MethodInfo, theme.HistoryIcon(), func() {
					fmt.Println("Tapped", historyValue.Id)
					h.lruCache.Get(historyValue.Id)
					h.targetInput.SetText(historyValue.Value.Target)
					h.requestEntry.SetText(historyValue.Value.Request)
					h.metadataEntry.SetText(historyValue.Value.Metadata)
					h.responseEntry.SetText(historyValue.Value.Response)
				})
				h.historyButton.Add(button)
				h.alreadyButtonList = append(h.alreadyButtonList, button)
			}
		}
	}
}
