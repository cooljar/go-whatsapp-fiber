package domain

import "errors"

var (
	ErrAlreadyConnected           = errors.New("already connected")
	ErrAlreadyLoggedIn            = errors.New("already logged in")
	ErrInvalidSession             = errors.New("invalid session")
	ErrLoginInProgress            = errors.New("login or restore already running")
	ErrNotConnected               = errors.New("not connected")
	ErrInvalidWsData              = errors.New("received invalid data")
	ErrInvalidWsState             = errors.New("can't handle binary data when not logged in")
	ErrConnectionTimeout          = errors.New("connection timed out")
	ErrMissingMessageTag          = errors.New("no messageTag specified or to short")
	ErrInvalidHmac                = errors.New("invalid hmac")
	ErrInvalidServerResponse      = errors.New("invalid response received from server")
	ErrServerRespondedWith404     = errors.New("server responded with status 404")
	ErrMediaDownloadFailedWith404 = errors.New("download failed with status code 404")
	ErrMediaDownloadFailedWith410 = errors.New("download failed with status code 410")
	ErrInvalidWebsocket           = errors.New("invalid websocket")
	ErrPhoneNotConnected          = errors.New("something when wrong while trying to ping, please check phone connectivity")

	ErrOptionsNotProvided         = errors.New("new conn options not provided")
)