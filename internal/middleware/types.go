package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type Locals struct {
	CustomerID int
	Nik        string
	Email      string
	FullName   string
}

func GetLocals(c *fiber.Ctx) *Locals {
	var l = Locals{}
	CustomerID, ok := c.Locals("customer_id").(int)
	if ok {
		l.CustomerID = CustomerID
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get user_id from locals")
	}

	nik, ok := c.Locals("nik").(string)
	if ok {
		l.Nik = nik
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get nik from locals")
	}

	email, ok := c.Locals("email").(string)
	if ok {
		l.Email = email
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get email from locals")
	}

	fullName, ok := c.Locals("full_name").(string)
	if ok {
		l.FullName = fullName
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get full_name from locals")
	}

	return &l
}

func (l *Locals) GetCustomerID() int {
	return l.CustomerID
}

func (l *Locals) GetNik() string {
	return l.Nik
}

func (l *Locals) GetEmail() string {
	return l.Email
}

func (l *Locals) GetFullName() string {
	return l.FullName
}
