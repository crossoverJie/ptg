package call

import (
	"context"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/crossoverJie/ptg/reflect"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"io"
)

type (
	Call struct {
		parse                       *reflect.ParseReflect
		responseEntry, requestEntry *widget.Entry
		processBar                  *widget.ProgressBarInfinite
		mds                         *desc.MethodDescriptor
		stub                        grpcdynamic.Stub
		errorHandle                 errorHandle
	}
	errorHandle func(window fyne.Window, processBar *widget.ProgressBarInfinite, err error)
)

func NewCallBuilder() *Call {
	return &Call{}
}
func (c *Call) Parse(parse *reflect.ParseReflect) *Call {
	c.parse = parse
	return c
}
func (c *Call) ResponseEntry(responseEntry *widget.Entry) *Call {
	c.responseEntry = responseEntry
	return c
}
func (c *Call) Mds(mds *desc.MethodDescriptor) *Call {
	c.mds = mds
	return c
}
func (c *Call) Stub(stub grpcdynamic.Stub) *Call {
	c.stub = stub
	return c
}
func (c *Call) RequestEntry(requestEntry *widget.Entry) *Call {
	c.requestEntry = requestEntry
	return c
}
func (c *Call) ProcessBar(processBar *widget.ProgressBarInfinite) *Call {
	c.processBar = processBar
	return c
}
func (c *Call) ErrorHandle(e errorHandle) *Call {
	c.errorHandle = e
	return c
}

func (c *Call) Run(ctx context.Context) (string, error) {
	c.responseEntry.SetText("")
	mds := c.mds
	parse := c.parse
	responseEntry := c.responseEntry
	stub := c.stub
	data := c.requestEntry.Text

	if !mds.IsClientStreaming() && !mds.IsServerStreaming() {
		// unary RPC
		rpc, err := parse.InvokeRpc(ctx, stub, mds, data)
		if err != nil {
			return "", err
		}
		viewCallBack := &unaryCallBack{}
		return viewCallBack.ViewCallBack(rpc, responseEntry), nil
	}

	if !mds.IsClientStreaming() && mds.IsServerStreaming() {
		// server stream
		rpc, err := parse.InvokeServerStreamRpc(ctx, stub, mds, data)
		if err != nil {
			return "", err
		}
		viewCallBack := &serverStreamCallBack{
			ch:            make(chan proto.Message, 10),
			responseEntry: responseEntry,
		}
		viewCallBack.receive()
		for {
			msg, err := rpc.RecvMsg()
			fmt.Printf("stream call %s \n", msg)
			if err == io.EOF {
				close(viewCallBack.ch)
			}
			if err != nil {
				return "", nil
			}
			viewCallBack.ViewCallBack(msg, responseEntry)
		}

	}

	if mds.IsClientStreaming() && !mds.IsServerStreaming() {
		// client stream
		rpc, err := parse.InvokeClientStreamRpc(ctx, stub, mds)
		if err != nil {
			return "", err
		}
		w := fyne.CurrentApp().NewWindow("Client stream call")
		w.Resize(fyne.NewSize(300, 100))
		request := widget.NewMultiLineEntry()
		request.SetText(data)
		var totalRequest string
		requestButton := widget.NewButtonWithIcon("Push", theme.MediaPlayIcon(), func() {
			messages, err := reflect.CreatePayloadsFromJSON(mds, request.Text)
			if err != nil {
				c.errorHandle(w, c.processBar, err)
				return
			}
			rpc.SendMsg(messages[0])
			totalRequest += request.Text + "\n"
		})
		finishButton := widget.NewButtonWithIcon("Finish", theme.ConfirmIcon(), func() {
			receive, err := rpc.CloseAndReceive()
			if err != nil {
				c.errorHandle(w, c.processBar, err)
				return
			}
			marshalIndent, _ := json.MarshalIndent(receive, "", "\t")
			c.responseEntry.SetText(string(marshalIndent))
			w.Close()
			c.requestEntry.SetText(totalRequest)

		})
		w.SetContent(container.NewVBox(request, requestButton, finishButton))
		w.CenterOnScreen()
		w.Show()

	}
	if mds.IsClientStreaming() && mds.IsServerStreaming() {
		// bidi stream
		rpc, err := parse.InvokeBidiStreamRpc(ctx, stub, mds)
		if err != nil {
			return "", err
		}
		w := fyne.CurrentApp().NewWindow("Bidi stream call")
		w.Resize(fyne.NewSize(300, 100))
		request := widget.NewMultiLineEntry()
		request.SetText(data)
		var totalRequest string

		streamCallBack := bidiStreamCallBack{
			ch:            make(chan proto.Message, 10),
			responseEntry: responseEntry,
		}
		streamCallBack.receive()

		requestButton := widget.NewButtonWithIcon("Push", theme.MediaPlayIcon(), func() {
			messages, err := reflect.CreatePayloadsFromJSON(mds, request.Text)
			if err != nil {
				c.errorHandle(w, c.processBar, err)
				return
			}
			rpc.SendMsg(messages[0])
			totalRequest += request.Text + "\n"

			receive, _ := rpc.RecvMsg()
			streamCallBack.ViewCallBack(receive)
		})
		finishButton := widget.NewButtonWithIcon("Finish", theme.ConfirmIcon(), func() {
			err := rpc.CloseSend()
			if err != nil {
				c.errorHandle(w, c.processBar, err)
				return
			}
			w.Close()
			c.requestEntry.SetText(totalRequest)

		})
		w.SetContent(container.NewVBox(request, requestButton, finishButton))
		w.CenterOnScreen()
		w.Show()
	}

	return "", nil
}

type unaryCallBack struct{}

func (u unaryCallBack) ViewCallBack(message proto.Message, responseEntry *widget.Entry) string {
	marshalIndent, _ := json.MarshalIndent(message, "", "\t")
	responseEntry.SetText(string(marshalIndent))
	return string(marshalIndent)
}

type serverStreamCallBack struct {
	ch            chan proto.Message
	responseEntry *widget.Entry
}

func (s *serverStreamCallBack) ViewCallBack(message proto.Message, responseEntry *widget.Entry) string {
	s.ch <- message
	return ""
}

func (s *serverStreamCallBack) receive() {
	go func() {
		for {
			select {
			case message, ok := <-s.ch:
				if !ok {
					fmt.Println("break")
					return
				}
				marshalIndent, _ := json.MarshalIndent(message, "", "\t")
				s.responseEntry.SetText(s.responseEntry.Text + string(marshalIndent) + "\n")
			}
		}
	}()

}

type clientStreamCallBack struct{}

func (u clientStreamCallBack) ViewCallBack(message proto.Message, responseEntry *widget.Entry) string {
	marshalIndent, _ := json.MarshalIndent(message, "", "\t")
	responseEntry.SetText(string(marshalIndent))
	return string(marshalIndent)
}

type bidiStreamCallBack struct {
	ch            chan proto.Message
	responseEntry *widget.Entry
}

func (b bidiStreamCallBack) ViewCallBack(message proto.Message) string {
	b.ch <- message
	return ""
}

func (b *bidiStreamCallBack) receive() {
	go func() {
		for {
			select {
			case message, ok := <-b.ch:
				if !ok {
					fmt.Println("break")
					return
				}
				marshalIndent, _ := json.MarshalIndent(message, "", "\t")
				b.responseEntry.SetText(b.responseEntry.Text + string(marshalIndent) + "\n")
			}
		}
	}()

}
