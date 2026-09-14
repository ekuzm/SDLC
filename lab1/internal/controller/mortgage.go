package controller

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ekuzm/sdlc-lab1/internal/apperror"
	"github.com/ekuzm/sdlc-lab1/internal/model"
)

const maxAmount = 1_000_000

type MortgageModel interface {
	Calculate(model.MortgageInput) error
	State() model.State
}

type MortgageForm struct {
	Amount     string
	AnnualRate string
	TermMonths string
}

type Mortgage struct {
	model MortgageModel
}

func NewMortgage(mortgageModel MortgageModel) *Mortgage {
	return &Mortgage{model: mortgageModel}
}

func (c *Mortgage) LastInput() model.MortgageInput {
	return c.model.State().Input
}

func (c *Mortgage) Calculate(form MortgageForm) error {
	amount, err := parseFloat(form.Amount, apperror.ErrAmountRequired, apperror.ErrAmountFormat)
	if err != nil {
		return fmt.Errorf("parse amount: %w", err)
	}

	rate, err := parseFloat(form.AnnualRate, apperror.ErrRateRequired, apperror.ErrRateFormat)
	if err != nil {
		return fmt.Errorf("parse annual rate: %w", err)
	}

	termMonths, err := parseInt(form.TermMonths)
	if err != nil {
		return fmt.Errorf("parse term: %w", err)
	}

	input := model.MortgageInput{
		Amount:     amount,
		AnnualRate: rate,
		TermMonths: termMonths,
	}
	if err = c.model.Calculate(input); err != nil {
		return fmt.Errorf("calculate mortgage: %w", err)
	}

	return nil
}

func parseFloat(value string, requiredErr, formatErr error) (float64, error) {
	value = normalizeNumber(value)
	if value == "" {
		return 0, requiredErr
	}

	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, formatErr
	}

	if result > maxAmount {
		return 0, apperror.ErrMoreThenMaxAmount
	}

	return result, nil
}

func parseInt(value string) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, apperror.ErrTermRequired
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, apperror.ErrTermFormat
	}

	return result, nil
}

func normalizeNumber(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, "\u00a0", "")
	value = strings.ReplaceAll(value, ",", ".")

	return value
}
