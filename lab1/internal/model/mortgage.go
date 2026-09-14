package model

import (
	"math"
	"sync"

	"github.com/ekuzm/sdlc-lab1/internal/apperror"
)

type MortgageInput struct {
	Amount     float64
	AnnualRate float64
	TermMonths int
}

type MortgageResult struct {
	MonthlyPayment     float64
	TotalPayment       float64
	Overpayment        float64
	OverpaymentFactor  float64
	OverpaymentPercent float64
}

type State struct {
	Input          MortgageInput
	Result         MortgageResult
	HasCalculation bool
}

type Observer interface {
	ModelChanged(State)
}

type Mortgage struct {
	mu        sync.RWMutex
	state     State
	observers []Observer
}

func NewMortgage() *Mortgage {
	return &Mortgage{}
}

func (m *Mortgage) Subscribe(observer Observer) {
	m.mu.Lock()
	m.observers = append(m.observers, observer)
	state := m.state
	m.mu.Unlock()

	observer.ModelChanged(state)
}

func (m *Mortgage) State() State {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.state
}

func (m *Mortgage) Calculate(input MortgageInput) error {
	if err := validate(input); err != nil {
		return err
	}

	result, err := calculate(input)
	if err != nil {
		return err
	}

	m.mu.Lock()
	m.state = State{
		Input:          input,
		Result:         result,
		HasCalculation: true,
	}
	state := m.state
	observers := append([]Observer(nil), m.observers...)
	m.mu.Unlock()

	for _, observer := range observers {
		observer.ModelChanged(state)
	}

	return nil
}

func validate(input MortgageInput) error {
	if input.Amount <= 0 || math.IsNaN(input.Amount) || math.IsInf(input.Amount, 0) {
		return apperror.ErrAmountPositive
	}

	if input.AnnualRate < 0 || input.AnnualRate > 100 || math.IsNaN(input.AnnualRate) || math.IsInf(input.AnnualRate, 0) {
		return apperror.ErrRateRange
	}

	if input.TermMonths < 1 || input.TermMonths > 1200 {
		return apperror.ErrTermRange
	}

	return nil
}

func calculate(input MortgageInput) (MortgageResult, error) {
	months := float64(input.TermMonths)
	monthlyRate := input.AnnualRate / 100 / 12

	monthlyPayment := input.Amount / months
	if monthlyRate > 0 {
		growth := math.Pow(1+monthlyRate, months)
		monthlyPayment = input.Amount * monthlyRate * growth / (growth - 1)
	}

	totalPayment := monthlyPayment * months
	overpayment := totalPayment - input.Amount
	result := MortgageResult{
		MonthlyPayment:     monthlyPayment,
		TotalPayment:       totalPayment,
		Overpayment:        overpayment,
		OverpaymentFactor:  totalPayment / input.Amount,
		OverpaymentPercent: overpayment / input.Amount * 100,
	}

	if math.IsNaN(result.MonthlyPayment) || math.IsInf(result.MonthlyPayment, 0) {
		return MortgageResult{}, apperror.ErrCalculation
	}

	return result, nil
}
