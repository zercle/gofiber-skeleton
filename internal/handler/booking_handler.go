package handler

import (
	"strconv"
	"time"

	"gofiber-skeleton/internal/domain"
	"gofiber-skeleton/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type BookingHandler struct {
	uc usecase.BookingUsecase
}

func NewBookingHandler(uc usecase.BookingUsecase) *BookingHandler {
	return &BookingHandler{uc: uc}
}

func RegisterRoutes(app *fiber.App, uc usecase.BookingUsecase) {
	h := NewBookingHandler(uc)
	app.Post("/bookings", h.CreateBooking)
	app.Get("/bookings/:id", h.GetBooking)
	app.Put("/bookings/:id", h.UpdateBooking)
	app.Delete("/bookings/:id", h.DeleteBooking)
	app.Get("/users/:user_id/bookings", h.ListBookingsByUser)
}

type createBookingRequest struct {
	UserID    uint   `json:"user_id"`
	ItemID    uint   `json:"item_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	var req createBookingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request payload"})
	}
	start, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid start_date format"})
	}
	end, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid end_date format"})
	}
	booking := &domain.Booking{
		UserID:    req.UserID,
		ItemID:    req.ItemID,
		StartDate: start,
		EndDate:   end,
	}
	if err := h.uc.CreateBooking(booking); err != nil {
		switch err {
		case usecase.ErrInvalidBooking, usecase.ErrBookingDateConflict:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create booking"})
		}
	}
	return c.Status(fiber.StatusCreated).JSON(booking)
}

func (h *BookingHandler) GetBooking(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid booking id"})
	}
	booking, err := h.uc.GetBooking(uint(id))
	if err != nil {
		if err == usecase.ErrBookingNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "booking not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve booking"})
	}
	return c.JSON(booking)
}

type updateBookingRequest struct {
	UserID    uint   `json:"user_id"`
	ItemID    uint   `json:"item_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (h *BookingHandler) UpdateBooking(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid booking id"})
	}
	var req updateBookingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request payload"})
	}
	start, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid start_date format"})
	}
	end, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid end_date format"})
	}
	booking := &domain.Booking{
		ID:        uint(id),
		UserID:    req.UserID,
		ItemID:    req.ItemID,
		StartDate: start,
		EndDate:   end,
	}
	if err := h.uc.UpdateBooking(booking); err != nil {
		switch err {
		case usecase.ErrInvalidBooking, usecase.ErrBookingDateConflict:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		case usecase.ErrBookingNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "booking not found"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update booking"})
		}
	}
	return c.JSON(booking)
}

func (h *BookingHandler) DeleteBooking(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid booking id"})
	}
	if err := h.uc.DeleteBooking(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete booking"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *BookingHandler) ListBookingsByUser(c *fiber.Ctx) error {
	userParam := c.Params("user_id")
	userID, err := strconv.ParseUint(userParam, 10, 64)
	if err != nil || userID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}
	bookings, err := h.uc.ListBookingsByUser(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list bookings"})
	}
	return c.JSON(bookings)
}