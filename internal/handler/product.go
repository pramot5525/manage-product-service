package handler

import (
	"encoding/json"
	"prm-product/internal/model"
	"prm-product/internal/service"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	service  service.ProductService
	validate *validator.Validate
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req model.ProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Successful: false,
			ErrorCode:  "invalid_request_body",
		})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Successful: false,
			ErrorCode:  "validation_error",
		})
	}

	id, err := h.service.Create(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Successful: false,
			ErrorCode:  "failed_to_create_product",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.APIResponse{
		Successful: true,
		Data: &model.Product{
			ID:          id,
			Name:        req.Name,
			Description: req.Description,
			SalePrice:   req.SalePrice,
			Price:       req.Price,
		},
	})
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Successful: false,
			ErrorCode:  "id_required",
		})
	}

	var raw map[string]json.RawMessage
	if err := c.BodyParser(&raw); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Successful: false,
			ErrorCode:  "invalid_request_body",
		})
	}

	updateData, err := h.buildPatchUpdateData(raw)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Successful: false,
			ErrorCode:  err.Error(),
		})
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse{
			Successful: false,
			ErrorCode:  "invalid_id",
		})
	}

	if err := h.service.Update(idInt, updateData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse{
			Successful: false,
			ErrorCode:  "failed_to_update_product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.APIResponse{
		Successful: true,
	})
}

func (h *ProductHandler) buildPatchUpdateData(raw map[string]json.RawMessage) (map[string]interface{}, error) {
	allowedFields := map[string]bool{
		"name":        true,
		"description": true,
		"sale_price":  true,
		"price":       true,
	}

	updateData := map[string]interface{}{}

	for key := range raw {
		if !allowedFields[key] {
			return nil, fiber.NewError(fiber.StatusBadRequest, "validation_error")
		}
	}

	if value, ok := raw["name"]; ok {
		if string(value) == "null" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "validation_error")
		}

		var name string
		if err := json.Unmarshal(value, &name); err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "validation_error")
		}
		if err := h.validate.Var(name, "required"); err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "validation_error")
		}

		updateData["name"] = name
	}

	if value, ok := raw["description"]; ok {
		if string(value) == "null" {
			updateData["description"] = nil
		} else {
			var description string
			if err := json.Unmarshal(value, &description); err != nil {
				return nil, fiber.NewError(fiber.StatusBadRequest, "validation_error")
			}
			updateData["description"] = description
		}
	}

	if value, ok := raw["sale_price"]; ok {
		if string(value) == "null" {
			updateData["sale_price"] = nil
		} else {
			var salePrice float64
			if err := json.Unmarshal(value, &salePrice); err != nil {
				return nil, fiber.NewError(fiber.StatusBadRequest, "validation_error")
			}
			if err := h.validate.Var(salePrice, "gte=0"); err != nil {
				return nil, fiber.NewError(fiber.StatusBadRequest, "validation_error")
			}

			updateData["sale_price"] = salePrice
		}
	}

	if value, ok := raw["price"]; ok {
		if string(value) == "null" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "validation_error")
		}

		var price float64
		if err := json.Unmarshal(value, &price); err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "validation_error")
		}
		if err := h.validate.Var(price, "gte=0"); err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "validation_error")
		}

		updateData["price"] = price
	}

	if len(updateData) == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "validation_error")
	}

	return updateData, nil
}
