package view

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/ekuzm/sdlc-lab1/internal/apperror"
	"github.com/ekuzm/sdlc-lab1/internal/controller"
	"github.com/ekuzm/sdlc-lab1/internal/model"
)

type Mortgage struct {
	app        fyne.App
	window     fyne.Window
	controller *controller.Mortgage

	amountValue      *widget.Label
	rateValue        *widget.Label
	termValue        *widget.Label
	monthlyValue     *widget.Label
	totalValue       *widget.Label
	overpaymentValue *widget.Label
	factorValue      *widget.Label
}

func NewMortgage(application fyne.App, mortgageController *controller.Mortgage) *Mortgage {
	application.Settings().SetTheme(newModernTheme())

	view := &Mortgage{
		app:        application,
		window:     application.NewWindow("Кредит на жильё"),
		controller: mortgageController,
	}
	view.build()

	return view
}

func (v *Mortgage) ShowAndRun() {
	v.window.ShowAndRun()
}

func (v *Mortgage) ModelChanged(state model.State) {
	if !state.HasCalculation {
		v.renderEmptyState()
		return
	}

	v.amountValue.SetText(formatMoney(state.Input.Amount))
	v.rateValue.SetText(formatPercent(state.Input.AnnualRate))
	v.termValue.SetText(formatTerm(state.Input.TermMonths))
	v.monthlyValue.SetText(formatMoney(state.Result.MonthlyPayment))
	v.totalValue.SetText(formatMoney(state.Result.TotalPayment))
	v.overpaymentValue.SetText(formatMoney(state.Result.Overpayment))
	v.factorValue.SetText(formatFactor(state.Result.OverpaymentFactor))
}

func (v *Mortgage) build() {
	v.window.Resize(fyne.NewSize(960, 650))
	v.window.CenterOnScreen()

	v.amountValue = valueLabel()
	v.rateValue = valueLabel()
	v.termValue = valueLabel()
	v.monthlyValue = heroValueLabel()
	v.totalValue = valueLabel()
	v.overpaymentValue = valueLabel()
	v.factorValue = valueLabel()

	title := widget.NewLabelWithStyle("Кредит на жильё", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	title.TextStyle = fyne.TextStyle{Bold: true}
	subtitle := widget.NewLabel("Рассчитайте платёжи спокойно и прозрачно")
	subtitle.Importance = widget.LowImportance

	header := container.NewVBox(title, subtitle)

	inputCard := surfaceCard(
		"Параметры кредита",
		"Последние введённые данные сохраняются",
		container.NewGridWithColumns(3,
			metric("Сумма", v.amountValue),
			metric("Ставка", v.rateValue),
			metric("Срок", v.termValue),
		),
		false,
	)

	monthlyCard := surfaceCard(
		"Ежемесячный платёж",
		"Аннуитетный платёж",
		container.NewCenter(v.monthlyValue),
		true,
	)

	results := container.NewGridWithColumns(3,
		metricCard("Всего банку", "Сумма выплат за весь срок", v.totalValue),
		metricCard("Переплата", "Проценты сверх суммы кредита", v.overpaymentValue),
		metricCard("Итоговый коэффициент", "Во сколько раз больше суммы кредита", v.factorValue),
	)

	calculateButton := widget.NewButton("Ввести данные", v.showInputWindow)
	calculateButton.Importance = widget.HighImportance

	actions := container.NewBorder(nil, nil, nil, calculateButton)

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		inputCard,
		monthlyCard,
		results,
		actions,
	)

	v.window.SetContent(container.NewPadded(container.NewVScroll(content)))
	v.renderEmptyState()
}

func (v *Mortgage) renderEmptyState() {
	const empty = "—"
	v.amountValue.SetText(empty)
	v.rateValue.SetText(empty)
	v.termValue.SetText(empty)
	v.monthlyValue.SetText("Введите параметры")
	v.totalValue.SetText(empty)
	v.overpaymentValue.SetText(empty)
	v.factorValue.SetText(empty)
}

func (v *Mortgage) showInputWindow() {
	input := v.controller.LastInput()
	inputWindow := v.app.NewWindow("Параметры кредита")
	inputWindow.Resize(fyne.NewSize(560, 440))
	inputWindow.CenterOnScreen()

	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("Например, 250 000")
	rateEntry := widget.NewEntry()
	rateEntry.SetPlaceHolder("Например, 14,5")
	termEntry := widget.NewEntry()
	termEntry.SetPlaceHolder("Например, 240")

	if input.Amount > 0 {
		amountEntry.SetText(formatInputFloat(input.Amount))
		rateEntry.SetText(formatInputFloat(input.AnnualRate))
		termEntry.SetText(strconv.Itoa(input.TermMonths))
	}

	form := widget.NewForm(
		widget.NewFormItem("Сумма кредита, BYN", amountEntry),
		widget.NewFormItem("Годовая ставка, %", rateEntry),
		widget.NewFormItem("Срок кредита, месяцев", termEntry),
	)

	intro := widget.NewLabel("Укажите условия кредита. Дробную ставку можно вводить через запятую или точку.")
	intro.Wrapping = fyne.TextWrapWord
	intro.Importance = widget.LowImportance

	cancelButton := widget.NewButton("Отмена", inputWindow.Close)
	calculateButton := widget.NewButton("Рассчитать", func() {
		err := v.controller.Calculate(controller.MortgageForm{
			Amount:     amountEntry.Text,
			AnnualRate: rateEntry.Text,
			TermMonths: termEntry.Text,
		})
		if err != nil {
			dialog.ShowError(userFacingError(err), inputWindow)
			return
		}

		inputWindow.Close()
	})
	calculateButton.Importance = widget.HighImportance

	actions := container.NewHBox(layout.NewSpacer(), cancelButton, calculateButton)
	content := container.NewVBox(
		widget.NewLabelWithStyle("Новый расчёт", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		intro,
		widget.NewSeparator(),
		form,
		layout.NewSpacer(),
		actions,
	)

	inputWindow.SetContent(container.NewPadded(content))
	inputWindow.Show()
	inputWindow.RequestFocus()
}

func metric(title string, value *widget.Label) fyne.CanvasObject {
	label := widget.NewLabel(title)
	label.Importance = widget.LowImportance

	return container.NewVBox(
		container.NewCenter(label),
		value,
	)
}

func metricCard(title, subtitle string, value *widget.Label) fyne.CanvasObject {
	return surfaceCard(title, subtitle, container.NewCenter(value), false)
}

func surfaceCard(title, subtitle string, content fyne.CanvasObject, emphasized bool) fyne.CanvasObject {
	backgroundColor := surfaceColor
	borderColor := surfaceBorderColor
	if emphasized {
		backgroundColor = raisedSurfaceColor
		borderColor = accentBorderColor
	}

	background := canvas.NewRectangle(backgroundColor)
	background.CornerRadius = 10
	background.StrokeColor = borderColor
	background.StrokeWidth = 1

	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	subtitleLabel := widget.NewLabel(subtitle)
	subtitleLabel.Importance = widget.LowImportance

	body := container.NewVBox(titleLabel, subtitleLabel, content)
	return container.NewStack(background, container.NewPadded(body))
}

func valueLabel() *widget.Label {
	return widget.NewLabelWithStyle("—", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
}

func heroValueLabel() *widget.Label {
	label := widget.NewLabelWithStyle("Введите параметры", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	label.Importance = widget.HighImportance

	return label
}

func formatTerm(months int) string {
	return fmt.Sprintf("%d %s", months, plural(months, "месяц", "месяца", "месяцев"))
}

func plural(number int, one, few, many string) string {
	lastTwo := number % 100
	if lastTwo >= 11 && lastTwo <= 14 {
		return many
	}

	switch number % 10 {
	case 1:
		return one
	case 2, 3, 4:
		return few
	default:
		return many
	}
}

func formatMoney(value float64) string {
	return groupThousands(fmt.Sprintf("%.2f", value)) + " BYN"
}

func formatPercent(value float64) string {
	return strings.ReplaceAll(formatInputFloat(value), ".", ",") + "%"
}

func formatFactor(value float64) string {
	return strings.ReplaceAll(fmt.Sprintf("%.2f", value), ".", ",") + " раза"
}

func formatInputFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func groupThousands(value string) string {
	parts := strings.SplitN(value, ".", 2)
	integer := parts[0]

	for index := len(integer) - 3; index > 0; index -= 3 {
		integer = integer[:index] + " " + integer[index:]
	}

	if len(parts) == 1 {
		return integer
	}

	return integer + "," + parts[1]
}

func userFacingError(err error) error {
	knownErrors := []error{
		apperror.ErrAmountRequired,
		apperror.ErrAmountFormat,
		apperror.ErrAmountPositive,
		apperror.ErrRateRequired,
		apperror.ErrRateFormat,
		apperror.ErrRateRange,
		apperror.ErrTermRequired,
		apperror.ErrTermFormat,
		apperror.ErrTermRange,
		apperror.ErrCalculation,
	}

	for _, knownErr := range knownErrors {
		if errors.Is(err, knownErr) {
			return knownErr
		}
	}

	return apperror.ErrCalculation
}
