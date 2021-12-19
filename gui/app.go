package main

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
