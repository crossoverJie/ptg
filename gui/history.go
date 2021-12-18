package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/crossoverJie/ptg/gui/io"
	"strings"
)

type (
	History struct {
		lruCache          *LruCache
		writeSearchChan   chan struct{}
		searchChan        chan []*HistoryValue
		historyButton     *fyne.Container
		alreadyButtonList []*widget.Button
		targetInput       *widget.Entry
		requestEntry      *widget.Entry
		metadataEntry     *widget.Entry
		responseEntry     *widget.Entry
	}

	HistoryValue struct {
		Id         int     `json:"id"`
		Value      *io.Log `json:"value"`
		MethodInfo string  `json:"method_info"`
	}
)

func NewHistory(size int, historyButton *fyne.Container, targetInput, requestEntry, metadataEntry, responseEntry *widget.Entry) *History {
	h := &History{
		lruCache:        NewLruList(size),
		writeSearchChan: make(chan struct{}, size),
		searchChan:      make(chan []*HistoryValue, size),
		historyButton:   historyButton,
		targetInput:     targetInput,
		requestEntry:    requestEntry,
		metadataEntry:   metadataEntry,
		responseEntry:   responseEntry,
	}
	go h.viewHistory()
	go h.ViewSearch()
	return h

}

func (h *History) Put(k int, v *HistoryValue) {
	h.lruCache.Put(k, v)
	h.writeSearchChan <- struct{}{}
}

func (h *History) viewHistory() {
	for {
		select {
		case <-h.writeSearchChan:

			// Reset view.
			for _, button := range h.alreadyButtonList {
				h.historyButton.Remove(button)
			}
			h.alreadyButtonList = make([]*widget.Button, 0)

			// Draw view.
			for _, v := range h.lruCache.List() {
				//index := i
				historyValue := v.(*HistoryValue)
				h.drawHistoryButton(historyValue)
			}
		}
	}
}

func (h *History) SearchResult(kw string) []*HistoryValue {
	var result []*HistoryValue
	for _, v := range h.lruCache.List() {
		historyValue := v.(*HistoryValue)
		if kw == "" {
			result = append(result, historyValue)
			continue
		}
		if strings.Contains(strings.ToLower(historyValue.MethodInfo), kw) {
			result = append(result, historyValue)
			continue
		}
		if strings.Contains(strings.ToLower(historyValue.Value.Target), kw) {
			result = append(result, historyValue)
			continue
		}
		if strings.Contains(strings.ToLower(historyValue.Value.Request), kw) {
			result = append(result, historyValue)
			continue
		}
		if strings.Contains(strings.ToLower(historyValue.Value.Response), kw) {
			result = append(result, historyValue)
			continue
		}
		if strings.Contains(strings.ToLower(historyValue.Value.Metadata), kw) {
			result = append(result, historyValue)
			continue
		}

	}
	h.searchChan <- result

	return result
}

func (h *History) ViewSearch() {
	for {
		select {
		case searchList := <-h.searchChan:
			// Reset view.
			for _, button := range h.alreadyButtonList {
				h.historyButton.Remove(button)
			}
			h.alreadyButtonList = make([]*widget.Button, 0)
			for _, v := range searchList {
				historyValue := v
				h.drawHistoryButton(historyValue)
			}
		}
	}

}

func (h *History) drawHistoryButton(historyValue *HistoryValue) {
	button := widget.NewButtonWithIcon(historyValue.MethodInfo, theme.HistoryIcon(), func() {
		fmt.Println("Search tapped", historyValue.Id)
		h.lruCache.Get(historyValue.Id)
		h.targetInput.SetText(historyValue.Value.Target)
		h.requestEntry.SetText(historyValue.Value.Request)
		h.metadataEntry.SetText(historyValue.Value.Metadata)
		h.responseEntry.SetText(historyValue.Value.Response)
	})
	h.historyButton.Add(button)
	h.alreadyButtonList = append(h.alreadyButtonList, button)
}
