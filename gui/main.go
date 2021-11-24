package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/crossoverJie/ptg/reflect"
	_ "github.com/crossoverJie/ptg/reflect"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
	"log"
	"strings"
)

func main() {
	app := app.New()
	window := app.NewWindow("Ptg gRPC client")
	window.Resize(fyne.NewSize(1000, 500))

	requestEntry := widget.NewMultiLineEntry()
	requestEntry.SetPlaceHolder("Input request json")
	requestEntry.Wrapping = fyne.TextWrapWord
	responseEntry := widget.NewMultiLineEntry()
	responseEntry.Wrapping = fyne.TextWrapWord
	reqLabel := widget.NewLabel("Request")
	targetInput := widget.NewEntry()
	targetInput.SetText("127.0.0.1:5000")
	targetInput.SetPlaceHolder("")
	processBar := widget.NewProgressBarInfinite()
	var parse *reflect.ParseReflect

	content := container.NewVBox()
	fileOpen := dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		if uri != nil {
			parse, err = reflect.NewParse(uri.URI().Path())
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			maps := parse.ServiceInfoMaps()
			fmt.Println(maps)
			serviceAccordion := widget.NewAccordion()
			for k, v := range maps {
				var methods []string
				for _, s := range v {
					methods = append(methods, k+"."+s)
				}
				serviceAccordion.Append(&widget.AccordionItem{
					Title: k,
					Detail: widget.NewRadioGroup(methods, func(s string) {
						service, method, err := reflect.ParseServiceMethod(s)
						if err != nil {
							dialog.ShowError(err, window)
						}
						json, err := parse.RequestJSON(service, method)
						if err != nil {
							dialog.ShowError(err, window)
							return
						}
						requestEntry.SetText(json)
						reqLabel.SetText("Request" + ":" + s)

					}),
					Open: false,
				})
			}
			content.Add(serviceAccordion)
		}
	}, window)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			fileOpen.Show()
		}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {

		}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)
	content.Add(toolbar)
	leftTool := container.New(layout.NewGridLayout(1), content)

	//
	rightTool := container.NewVBox()
	form := widget.NewForm(&widget.FormItem{
		Text:     "Target:",
		Widget:   targetInput,
		HintText: "Input target url",
	})
	rightTool.Add(form)
	rightTool.Add(reqLabel)
	rightTool.Add(requestEntry)
	rightTool.Add(widget.NewButtonWithIcon("RUN", theme.MediaPlayIcon(), func() {
		if requestEntry.Text == "" {
			dialog.ShowError(errors.New("request json can not nil"), window)
			return
		}
		processBar.Show()
		processBar.Start()
		serviceInfo := strings.Split(reqLabel.Text, ":")[1]
		service, method, err := reflect.ParseServiceMethod(serviceInfo)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		mds, err := parse.MethodDescriptor(service, method)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithInsecure())
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, targetInput.Text, opts...)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		stub := grpcdynamic.NewStub(conn)
		rpc, err := parse.InvokeRpc(ctx, stub, mds, requestEntry.Text)
		if err != nil {
			processBar.Stop()
			processBar.Hide()
			dialog.ShowError(err, window)
			return
		}
		processBar.Hide()
		marshalIndent, _ := json.MarshalIndent(rpc, "", "\t")
		responseEntry.SetText(string(marshalIndent))
	}))
	rightTool.Add(processBar)
	processBar.Hide()

	rightTool.Add(widget.NewLabel("Response:"))
	rightTool.Add(responseEntry)

	split := container.NewHSplit(leftTool, rightTool)

	window.SetContent(split)
	window.ShowAndRun()
}
