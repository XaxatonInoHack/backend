package info_user

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"xaxaton/internal/lib/converter"
)

type Handler struct {
	info info
}

func NewHandler(i info) *Handler {
	return &Handler{
		info: i,
	}
}

func (h *Handler) Handle(ctx *fiber.Ctx) error {
	var r request

	err := ctx.BodyParser(&r)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	feedback, score, res, err := h.info.GetInfo(ctx.Context(), r.UserID)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(
		fiber.Map{
			"user_id":        feedback.UserID,
			"criteria_score": score,
			"criteria_text": map[string]string{
				"Профессионализм":    res["Профессионализм"],
				"Командная работа":   res["Командная работа"],
				"Коммуникабельность": res["Коммуникабельность"],
				"Инициативность":     res["Инициативность"],
			},
			"overall_score":  converter.OverAllScore(score),
			"overall_resume": res["Резюме"],
		})
}
