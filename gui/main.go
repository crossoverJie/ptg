package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/crossoverJie/ptg/reflect"
	_ "github.com/crossoverJie/ptg/reflect"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
	"image/color"
	"net/url"
	"strings"
)

func main() {
	app := app.New()
	window := app.NewWindow("PTG gRPC client")
	window.Resize(fyne.NewSize(1000, 500))

	requestEntry := widget.NewMultiLineEntry()
	requestEntry.SetPlaceHolder("Input request json")
	requestEntry.Wrapping = fyne.TextWrapWord
	responseEntry := widget.NewMultiLineEntry()
	responseEntry.Wrapping = fyne.TextWrapWord
	reqLabel := widget.NewLabel("Request:")
	targetInput := widget.NewEntry()
	targetInput.SetText("127.0.0.1:6001")
	targetInput.SetPlaceHolder("")
	processBar := widget.NewProgressBarInfinite()
	processBar.Hide()
	serviceAccordionRemove := false
	serviceAccordion := widget.NewAccordion()

	content := container.NewVBox()
	fileOpen := dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		if uri != nil {
			parseAdapter, exit, err := RegisterReflect(uri.URI().Path())
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if exit {
				dialog.ShowError(errors.New("proto file already exists"), window)
				return
			}

			maps := parseAdapter.Parse().ServiceInfoMaps()
			if serviceAccordionRemove {
				content.Add(serviceAccordion)
				serviceAccordionRemove = false
			}
			for k, v := range maps {
				var methods []string
				for _, s := range v {
					methods = append(methods, k+"."+s+"-"+fmt.Sprint(parseAdapter.Index()))
				}
				serviceAccordion.Append(&widget.AccordionItem{
					Title: k,
					Detail: widget.NewRadioGroup(methods, func(s string) {
						if s == "" {
							return
						}
						methodInfo := strings.Split(s, "-")
						service, method, err := reflect.ParseServiceMethod(methodInfo[0])
						if err != nil {
							dialog.ShowError(err, window)
							return
						}
						json, err := GetParseAdapter(methodInfo[1]).Parse().RequestJSON(service, method)
						if err != nil {
							dialog.ShowError(err, window)
							return
						}
						requestEntry.SetText(json)
						reqLabel.SetText(s)

					}),
					Open: false,
				})

			}
		}
	}, window)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			fileOpen.Show()
		}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			dialog.ShowInformation("Notice", "coming soon", window)
		}),
		widget.NewToolbarAction(theme.DeleteIcon(), func() {
			ClearReflect()
			content.Remove(serviceAccordion)
			serviceAccordionRemove = true
			serviceAccordion.Items = nil
			dialog.ShowInformation("Notice", "all proto files have been reset", window)
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			w := fyne.CurrentApp().NewWindow("Help")
			u, _ := url.Parse("https://github.com/crossoverJie/ptg")
			w.SetContent(container.New(layout.NewCenterLayout(), widget.NewHyperlink("help?", u)))
			w.Resize(fyne.NewSize(130, 100))
			w.SetFixedSize(true)
			w.Show()
		}),
	)
	content.Add(toolbar)
	content.Add(serviceAccordion)
	leftTool := container.New(layout.NewGridLayout(1), content)

	// Right
	form := widget.NewForm(&widget.FormItem{
		Text:     "Target:",
		Widget:   targetInput,
		HintText: "Input target url",
	})

	requestContainer := container.New(layout.NewGridLayoutWithColumns(1))
	requestContainer.Add(requestEntry)
	requestButton := widget.NewButtonWithIcon("RUN", theme.MediaPlayIcon(), func() {
		if requestEntry.Text == "" {
			dialog.ShowError(errors.New("request json can not nil"), window)
			return
		}
		if reqLabel.Text == "" {
			dialog.ShowError(errors.New("proto can not nil"), window)
			return
		}
		methodInfo := strings.Split(reqLabel.Text, "-")
		service, method, err := reflect.ParseServiceMethod(methodInfo[0])
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		index := methodInfo[1]
		parse := GetParseAdapter(index).Parse()
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
		processBar.Show()
		rpc, err := parse.InvokeRpc(ctx, stub, mds, requestEntry.Text)
		if err != nil {
			processBar.Hide()
			dialog.ShowError(err, window)
			return
		}
		processBar.Hide()
		marshalIndent, _ := json.MarshalIndent(rpc, "", "\t")
		responseEntry.SetText(string(marshalIndent))
	})
	bottomBox := container.NewVBox(canvas.NewLine(color.Black), requestButton)
	bottomBox.Add(canvas.NewLine(color.Black))
	bottomBox.Add(processBar)
	requestPanel := container.NewBorder(form, bottomBox, nil, nil)
	requestPanel.Add(requestContainer)

	responseContainer := container.New(layout.NewGridLayoutWithColumns(1))
	responseContainer.Add(responseEntry)
	responseLabel := widget.NewLabel("Response:")
	responsePanel := container.NewBorder(responseLabel, nil, nil, nil)
	responsePanel.Add(responseContainer)

	rightTool := container.NewGridWithColumns(1,
		requestPanel, responsePanel)
	split := container.NewHSplit(leftTool, rightTool)

	window.SetContent(split)
	window.ShowAndRun()
}

func RequestErr(window fyne.Window, pb *widget.ProgressBarInfinite, err error) {
	pb.Hide()
	dialog.ShowError(err, window)
}
