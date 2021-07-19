package http

type DebugHandler struct {

}

/*
func NewDebugHandler(rPublic, rPrivate fiber.Router) {
	handler := &DebugHandler{}

	rSys := rPublic.Group("/debug")
	rSys.Get("/", handler.Index)
}

// Index func for debugging.
// @Summary debugging
// @Description debugging.
// @Tags Debugging
// @Produce json
// @Success 200 {object} object "Description"
// @Failure 422 {object} []domain.HTTPErrorValidation
// @Failure 400 {object} domain.HTTPError
// @Failure 404 {object} domain.HTTPError
// @Failure 500 {object} domain.HTTPError
// @Router /v1/debug [get]
func (d *DebugHandler) Index(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"data": ""})
}
*/