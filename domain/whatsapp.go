package domain

import (
	"encoding/json"
	"mime/multipart"
)

type WaSendTextForm struct {
	Msisdn      string `json:"msisdn" validate:"required"`
	Text        string `json:"text" validate:"required"`
	MsgQuotedID string `json:"msg_quoted_id"`
	MsgQuoted   string `json:"msg_quoted"`
}

type WaSendLocationForm struct {
	Msisdn      string  `json:"msisdn" validate:"required"`
	Latitude    float64 `json:"latitude" validate:"required,latitude"`
	Longitude   float64 `json:"longitude" validate:"required,longitude"`
	MsgQuotedID string  `json:"msg_quoted_id"`
	MsgQuoted   string  `json:"msg_quoted"`
}

type WaSendFileForm struct {
	Msisdn      string `json:"msisdn" validate:"required"`
	MsgQuotedID string `json:"msg_quoted_id"`
	MsgQuoted   string `json:"msg_quoted"`
	Message     string `json:"message"`
	//File        string `json:"file" validate:"required,file"`
	FileHeader  *multipart.FileHeader
}

type WaWebServer struct {
	Version struct {
		Major int
		Minor int
		Build int
	}
}

type WaWebClient struct {
	Version struct {
		Major int
		Minor int
		Build int
	}
}

type WaWeb struct {
	Server WaWebServer
	Client WaWebClient
}

type WhatsappWeb struct {
	Wa            WaWeb
	SessionJid    string
	SessionID     string
	SessionFile   string
	SessionStart  uint64
	ReconnectTime int
}

// FromJSON decode json to book struct
func (w *WhatsappWeb) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, w)
}

// ToJSON encode book struct to json
func (w *WhatsappWeb) ToJSON() []byte {
	str, _ := json.Marshal(w)
	return str
}

// WhatsappUsecase represent the whatsapp's use cases
type WhatsappUsecase interface {
	RestoreSession() error
	Login(vMajor, vMinor, vBuild, timeout, reconnect int, clientNameShort, clientNameLong string) (qrCode string, err error)
	GetInfo() (info WaWeb, err error)
	SendText(form WaSendTextForm) (msgId string, err error)
	SendLocation(form WaSendLocationForm) (msgId string, err error)
	SendFile(form WaSendFileForm, fileType string) (msgId string, err error)
	Logout() (err error)
	Groups(jid string) (g string, err error)
}
