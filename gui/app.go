package main

import (
	"github.com/flopp/go-findfont"
	"os"
	"strings"
)

const (
	AppName                  = "PTG gRPC client"
	AppWeight                = 1000
	AppHeight                = 800
	HelpUrl                  = "https://github.com/crossoverJie/ptg"
	SearchFormText           = "SearchResult"
	SearchFormPlaceHolder    = "keyword"
	TargetFormText           = "Target:"
	TargetFormHintText       = "Input target url"
	RequestEntryPlaceHolder  = "Input request json"
	MetaDataAccordion        = "metadata"
	MetaDataInputPlaceHolder = "Input metadata json"
	TargetInputText          = "127.0.0.1:6001"
	RequestButtonText        = "RUN"
	ResponseLabelText        = "Response:"
)

type App struct {
	AppName               string
	AppWidth, AppHeight   float32
	HelpUrl               string
	SearchFormText        string
	SearchFormPlaceHolder string
	RightRequest          *RightRequest
	RightResponse         *RightResponse
}

type RightRequest struct {
	TargetFormText, TargetFormHintText               string
	RequestEntryPlaceHolder, TargetInputText         string
	MetaDataAccordionTitle, MetaDataInputPlaceHolder string
	RequestButtonText                                string
}

type RightResponse struct {
	ResponseLabelText string
}

func InitApp() *App {
	// init font
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		//楷体:simkai.ttf
		//黑体:simhei.ttf
		if strings.Contains(path, "阿里汉仪智能黑体") || strings.Contains(path, "simkai.ttf") || strings.Contains(path, "msyhl.ttc") {
			os.Setenv("FYNE_FONT", path)
			break
		}
	}
	return &App{
		AppName:               AppName,
		AppWidth:              AppWeight,
		AppHeight:             AppHeight,
		HelpUrl:               HelpUrl,
		SearchFormText:        SearchFormText,
		SearchFormPlaceHolder: SearchFormPlaceHolder,
		RightRequest: &RightRequest{
			TargetFormText:           TargetFormText,
			TargetFormHintText:       TargetFormHintText,
			RequestEntryPlaceHolder:  RequestEntryPlaceHolder,
			TargetInputText:          TargetInputText,
			MetaDataAccordionTitle:   MetaDataAccordion,
			MetaDataInputPlaceHolder: MetaDataInputPlaceHolder,
			RequestButtonText:        RequestButtonText,
		},
		RightResponse: &RightResponse{ResponseLabelText: ResponseLabelText},
	}
}
