package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/crossoverJie/ptg/client"
	"github.com/crossoverJie/ptg/meta"
	"log"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Form Widget")

	entry := widget.NewEntry()
	textArea := widget.NewMultiLineEntry()
	meta.NewResult()
	body := `{"order_id":27669403038384131}`
	grpcClient := client.NewGrpcClient(meta.NewMeta("127.0.0.1:5000", "",
		"", body,
		client.Grpc, "/Users/chenjie/Documents/dev/easi/easi-order-service/proto/order/v1/order.proto",
		"order.v1.OrderService.Detail", 1, 1, nil, nil))

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			//{Text: "Entry", Widget: entry},
			{Text: "json", Widget: textArea}},
		OnSubmit: func() { // optional, handle form submission
			log.Println("Form submitted:", entry.Text)
			log.Println("multiline:", textArea.Text)
			request, err := grpcClient.Request()
			log.Println("grpcClient:", request, err)

			myWindow.Close()
		},
	}
	myWindow.SetContent(form)
	myWindow.ShowAndRun()
}
