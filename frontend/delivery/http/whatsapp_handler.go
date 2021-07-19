package http

import (
	"github.com/cooljar/go-whatsapp-fiber/domain"
	"github.com/cooljar/go-whatsapp-fiber/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
	"os"
	"strconv"
)

type WhatsappHandler struct {
	WhatsappUsecase domain.WhatsappUsecase
	Validate *validator.Validate
}

func NewWhatsappHandler(whatsappUsecase domain.WhatsappUsecase, rPublic, rPrivate fiber.Router) {
	handler := &WhatsappHandler{
		WhatsappUsecase: whatsappUsecase,
		Validate: utils.NewValidator(),
	}

	rWa := rPublic.Group("/whatsapp")
	rWa.Post("/login", handler.Login)
	rWa.Get("/info", handler.GetInfo)
	rWa.Post("/send-text", handler.SendText)
	rWa.Post("/send-location", handler.SendLocation)
	rWa.Post("/send-document", handler.SendDocument)
	rWa.Post("/send-image", handler.SendImage)
	rWa.Post("/send-audio", handler.SendAudio)
	rWa.Post("/send-video", handler.SendVideo)
	rWa.Get("/groups/:jid", handler.Groups)
	rWa.Post("/logout", handler.Logout)
}

// Login func login whatsapp web.
// @Description Login to whatsapp web by scanning a QR Code.
// @Summary login whatsapp web
// @Tags Whatsapp
// @Accept mpfd
// @Produce png
// @Param reconnect formData int false "Reconnect, default to 50"
// @Param timeout formData int false "QR Scan timeout in second, default 20"
// @Param client_name_long formData string false "Long client name, default: Go Whatsapp REST Api Fiber"
// @Param client_name_short formData string false "Short client name, default: Go Whatsapp"
// @Param client_version_major formData int false "Whatsapp Client major version, default: 2"
// @Param client_version_minor formData int false "Whatsapp Client minor version, default: 2126"
// @Param client_version_build formData int false "Whatsapp Client build version, default: 11"
// @Success 200 {file} file "Description"
// @Failure 422 {object} []domain.HTTPErrorValidation
// @Failure 400 {object} domain.HTTPError
// @Failure 404 {object} domain.HTTPError
// @Failure 500 {object} domain.HTTPError
// @Router /v1/whatsapp/login [post]
func (w *WhatsappHandler) Login(c *fiber.Ctx) error {
	reconnect := c.FormValue("reconnect", "50")
	timeout := c.FormValue("timeout", "20")
	clientNameLong := c.FormValue("client_name_long", "Cooljar Whatsapp REST Api")
	clientNameShort := c.FormValue("client_name_short", "Cooljar Whatsapp")
	reqVersionClientMajor := c.FormValue("client_version_major", os.Getenv("WHATSAPP_CLIENT_VERSION_MAJOR"))
	reqVersionClientMinor := c.FormValue("client_version_minor", os.Getenv("WHATSAPP_CLIENT_VERSION_MINOR"))
	reqVersionClientBuild := c.FormValue("client_version_build", os.Getenv("WHATSAPP_CLIENT_VERSION_BUILD"))

	reqReconnect, err := strconv.Atoi(reconnect)
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusBadRequest, err)
	}

	reqTimeout, err := strconv.Atoi(timeout)
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusBadRequest, err)
	}

	reqVersionClientMajorInt, err := strconv.Atoi(reqVersionClientMajor)
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusBadRequest, err)
	}

	reqVersionClientMinorInt, err := strconv.Atoi(reqVersionClientMinor)
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusBadRequest, err)
	}

	reqVersionClientBuildInt, err := strconv.Atoi(reqVersionClientBuild)
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusBadRequest, err)
	}

	qrCodeStr, err := w.WhatsappUsecase.Login(reqVersionClientMajorInt, reqVersionClientMinorInt, reqVersionClientBuildInt, reqTimeout, reqReconnect, clientNameShort, clientNameLong)
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusInternalServerError, err)
	}

	qrCodePng, _ := qrcode.Encode(qrCodeStr, qrcode.Medium, 256)
	c.Set("content-type", "image/png")
	return c.Send(qrCodePng)
}

// GetInfo func for get info metadata.
// @Summary get info metadata
// @Description Get info metadata.
// @Tags Info
// @Produce json
// @Success 200 {object} domain.JSONResult{data=domain.WaWeb,message=string} "Description"
// @Failure 422 {object} []domain.HTTPErrorValidation
// @Failure 400 {object} domain.HTTPError
// @Failure 404 {object} domain.HTTPError
// @Failure 500 {object} domain.HTTPError
// @Router /v1/whatsapp/info [get]
func (w *WhatsappHandler) GetInfo(c *fiber.Ctx) error {
	info, err := w.WhatsappUsecase.GetInfo()
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(domain.JSONResult{
		Data: info,
		Message: "Success",
	})
}

// SendText func for send text.
// @Summary send text message
// @Description Send text message.
// @Tags Messaging
// @Accept mpfd
// @Produce json
// @Param msisdn formData string true "Destination number. eg: 6281255423 or group_creator-timstamp_created -> 6281271471566-1619679643 for group"
// @Param text formData string true "Message text"
// @Param msg_quoted_id formData string false "Message Quoted ID"
// @Param msg_quoted formData string false "Message Quoted"
// @Success 200 {object} domain.JSONResult{data=object,message=string} "Description"
// @Failure 422 {object} []domain.HTTPErrorValidation
// @Failure 400 {object} domain.HTTPError
// @Failure 404 {object} domain.HTTPError
// @Failure 500 {object} domain.HTTPError
// @Router /v1/whatsapp/send-text [post]
func (w *WhatsappHandler) SendText(c *fiber.Ctx) error {
	// Instantiate new Book struct
	var form domain.WaSendTextForm
	form.Msisdn = c.FormValue("msisdn")
	form.Text = c.FormValue("text")
	form.MsgQuotedID = c.FormValue("msg_quoted_id")
	form.MsgQuoted = c.FormValue("msg_quoted")

	// Validate form input
	err := w.Validate.Struct(&form)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ValidatorErrors(err))
	}

	msgId, err := w.WhatsappUsecase.SendText(form)
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusBadRequest, err)
	}

	return c.JSON(domain.JSONResult{
		Data: map[string]string{"message_id": msgId},
		Message: "Success",
	})
}

// SendLocation func for send location.
// @Summary send location message
// @Description Send location message.
// @Tags Messaging
// @Accept mpfd
// @Produce json
// @Param msisdn formData string true "Destination number. eg: 6281255423 or group_creator-timstamp_created -> 6281271471566-1619679643 for group"
// @Param latitude formData number false "Latitude. eg: -5.3836767"
// @Param longitude formData number false "Longitude. eg: 105.2937439"
// @Param msg_quoted_id formData string false "Message Quoted ID"
// @Param msg_quoted formData string false "Message Quoted"
// @Success 200 {object} domain.JSONResult{data=object,message=string} "Description"
// @Failure 422 {object} []domain.HTTPErrorValidation
// @Failure 400 {object} domain.HTTPError
// @Failure 404 {object} domain.HTTPError
// @Failure 500 {object} domain.HTTPError
// @Router /v1/whatsapp/send-location [post]
func (w *WhatsappHandler) SendLocation(c *fiber.Ctx) error {
	var err error

	var form domain.WaSendLocationForm
	form.Msisdn = c.FormValue("msisdn")
	form.MsgQuotedID = c.FormValue("msg_quoted_id")
	form.MsgQuoted = c.FormValue("msg_quoted")

	form.Latitude, err = strconv.ParseFloat(c.FormValue("latitude"), 64)
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusInternalServerError, err)
	}

	form.Longitude, err = strconv.ParseFloat(c.FormValue("longitude"), 64)
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusInternalServerError, err)
	}

	// Validate form input
	err = w.Validate.Struct(&form)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ValidatorErrors(err))
	}

	msgId, err := w.WhatsappUsecase.SendLocation(form)

	return c.JSON(domain.JSONResult{
		Data: map[string]string{"message_id": msgId},
		Message: "Success",
	})
}

// SendImage func for send image.
// @Summary send image message
// @Description Send image message.
// @Tags Messaging
// @Accept x-www-form-urlencoded
// @Produce json
// @Param msisdn formData string true "Destination number. eg: 6281255423 or group_creator-timstamp_created -> 6281271471566-1619679643 for group"
// @Param image_file formData file true "Image file"
// @Param msg_quoted_id formData string false "Message Quoted ID"
// @Param msg_quoted formData string false "Message Quoted"
// @Param message formData string false "Message to include"
// @Success 200 {object} domain.JSONResult{data=object,message=string} "Description"
// @Failure 422 {object} []domain.HTTPErrorValidation
// @Failure 400 {object} domain.HTTPError
// @Failure 404 {object} domain.HTTPError
// @Failure 500 {object} domain.HTTPError
// @Router /v1/whatsapp/send-image [post]
func (w *WhatsappHandler) SendImage(c *fiber.Ctx) error {
	var err error

	var form domain.WaSendFileForm
	form.Msisdn = c.FormValue("msisdn")
	form.MsgQuotedID = c.FormValue("msg_quoted_id")
	form.MsgQuoted = c.FormValue("msg_quoted")
	form.Message = c.FormValue("message")

	form.FileHeader, err = c.FormFile("image_file")
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusBadRequest, err)
	}
	//form.File = fmt.Sprintf("./%s", form.FileHeader.Filename)

	// Validate form input
	err = w.Validate.Struct(&form)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ValidatorErrors(err))
	}

	msgId, err := w.WhatsappUsecase.SendFile(form, "image")
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(domain.JSONResult{
		Data: map[string]string{"message_id": msgId},
		Message: "Success",
	})
}

// SendAudio func for send audio.
// @Summary send audio message
// @Description Send audio message.
// @Tags Messaging
// @Accept x-www-form-urlencoded
// @Produce json
// @Param msisdn formData string true "Destination number. eg: 6281255423 or group_creator-timstamp_created -> 6281271471566-1619679643 for group"
// @Param audio_file formData file true "Audio file"
// @Param msg_quoted_id formData string false "Message Quoted ID"
// @Param msg_quoted formData string false "Message Quoted"
// @Param message formData string false "Message to include"
// @Success 200 {object} domain.JSONResult{data=object,message=string} "Description"
// @Failure 422 {object} []domain.HTTPErrorValidation
// @Failure 400 {object} domain.HTTPError
// @Failure 404 {object} domain.HTTPError
// @Failure 500 {object} domain.HTTPError
// @Router /v1/whatsapp/send-audio [post]
func (w *WhatsappHandler) SendAudio(c *fiber.Ctx) error {
	var err error

	var form domain.WaSendFileForm
	form.Msisdn = c.FormValue("msisdn")
	form.MsgQuotedID = c.FormValue("msg_quoted_id")
	form.MsgQuoted = c.FormValue("msg_quoted")
	form.Message = c.FormValue("message")

	form.FileHeader, err = c.FormFile("audio_file")
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusBadRequest, err)
	}
	//form.File = fmt.Sprintf("./%s", form.FileHeader.Filename)

	// Validate form input
	err = w.Validate.Struct(&form)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ValidatorErrors(err))
	}

	msgId, err := w.WhatsappUsecase.SendFile(form, "audio")
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(domain.JSONResult{
		Data: map[string]string{"message_id": msgId},
		Message: "Success",
	})
}

// SendVideo func for send video.
// @Summary send video message
// @Description Send video message.
// @Tags Messaging
// @Accept x-www-form-urlencoded
// @Produce json
// @Param msisdn formData string true "Destination number. eg: 6281255423 or group_creator-timstamp_created -> 6281271471566-1619679643 for group"
// @Param video_file formData file true "Video file"
// @Param msg_quoted_id formData string false "Message Quoted ID"
// @Param msg_quoted formData string false "Message Quoted"
// @Param message formData string false "Message to include"
// @Success 200 {object} domain.JSONResult{data=object,message=string} "Description"
// @Failure 422 {object} []domain.HTTPErrorValidation
// @Failure 400 {object} domain.HTTPError
// @Failure 404 {object} domain.HTTPError
// @Failure 500 {object} domain.HTTPError
// @Router /v1/whatsapp/send-video [post]
func (w *WhatsappHandler) SendVideo(c *fiber.Ctx) error {
	var err error

	var form domain.WaSendFileForm
	form.Msisdn = c.FormValue("msisdn")
	form.MsgQuotedID = c.FormValue("msg_quoted_id")
	form.MsgQuoted = c.FormValue("msg_quoted")
	form.Message = c.FormValue("message")

	form.FileHeader, err = c.FormFile("video_file")
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusBadRequest, err)
	}
	//form.File = fmt.Sprintf("./%s", form.FileHeader.Filename)

	// Validate form input
	err = w.Validate.Struct(&form)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ValidatorErrors(err))
	}

	msgId, err := w.WhatsappUsecase.SendFile(form, "video")
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(domain.JSONResult{
		Data: map[string]string{"message_id": msgId},
		Message: "Success",
	})
}

// SendDocument func for send document.
// @Summary send document message
// @Description Send document message.
// @Tags Messaging
// @Accept x-www-form-urlencoded
// @Produce json
// @Param msisdn formData string true "Destination number. eg: 6281255423 or group_creator-timstamp_created -> 6281271471566-1619679643 for group"
// @Param document_file formData file true "Document file"
// @Param msg_quoted_id formData string false "Message Quoted ID"
// @Param msg_quoted formData string false "Message Quoted"
// @Param message formData string false "Message to include"
// @Success 200 {object} domain.JSONResult{data=object,message=string} "Description"
// @Failure 422 {object} []domain.HTTPErrorValidation
// @Failure 400 {object} domain.HTTPError
// @Failure 404 {object} domain.HTTPError
// @Failure 500 {object} domain.HTTPError
// @Router /v1/whatsapp/send-document [post]
func (w *WhatsappHandler) SendDocument(c *fiber.Ctx) error {
	var err error

	var form domain.WaSendFileForm
	form.Msisdn = c.FormValue("msisdn")
	form.MsgQuotedID = c.FormValue("msg_quoted_id")
	form.MsgQuoted = c.FormValue("msg_quoted")
	form.Message = c.FormValue("message")

	form.FileHeader, err = c.FormFile("document_file")
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusBadRequest, err)
	}
	//form.File = fmt.Sprintf("./%s", form.FileHeader.Filename)

	// Validate form input
	err = w.Validate.Struct(&form)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ValidatorErrors(err))
	}

	msgId, err := w.WhatsappUsecase.SendFile(form, "document")
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(domain.JSONResult{
		Data: map[string]string{"message_id": msgId},
		Message: "Success",
	})
}

// Groups func for get group metadata.
// @Summary get group metadata
// @Description Get group metadata by phone number.
// @Tags Info
// @Produce json
// @Param jid path string true "JID. eg: group_creator-timstamp_created -> 6281271471566-1619679643"
// @Success 200 {object} domain.JSONResult{data=domain.WaGroup,message=string} "Description"
// @Failure 422 {object} []domain.HTTPErrorValidation
// @Failure 400 {object} domain.HTTPError
// @Failure 404 {object} domain.HTTPError
// @Failure 500 {object} domain.HTTPError
// @Router /v1/whatsapp/groups/{jid} [get]
func (w *WhatsappHandler) Groups(c *fiber.Ctx) error {
	jid := c.Params("jid")
	if len(jid) == 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "invalid jid"})
	}

	groupMd, err := w.WhatsappUsecase.Groups(jid)
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusBadRequest, err)
	}

	var modelGroup domain.WaGroup
	err = modelGroup.FromJSON([]byte(groupMd))
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusBadRequest, err)
	}

	return c.JSON(domain.JSONResult{
		Data: modelGroup,
		Message: "Success",
	})
}

// Logout func logout whatsapp web.
// @Description Logout from whatsapp web.
// @Summary logout whatsapp web
// @Tags Whatsapp
// @Produce json
// @Success 200 {object} domain.JSONResult{data=string,message=string} "Description"
// @Failure 422 {object} []domain.HTTPErrorValidation
// @Failure 400 {object} domain.HTTPError
// @Failure 404 {object} domain.HTTPError
// @Failure 500 {object} domain.HTTPError
// @Router /v1/whatsapp/logout [post]
func (w *WhatsappHandler) Logout(c *fiber.Ctx) error {
	err := w.WhatsappUsecase.Logout()
	if err != nil {
		return domain.NewHttpError(c, fiber.StatusInternalServerError, err)
	}

	return c.JSON(domain.JSONResult{
		Data: "Logout",
		Message: "Success",
	})
}
