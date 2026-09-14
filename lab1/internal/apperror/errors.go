package apperror

type Code string

const (
	CodeAmountRequired Code = "AMOUNT_REQUIRED"
	CodeAmountFormat   Code = "AMOUNT_FORMAT"
	CodeAmountPositive Code = "AMOUNT_POSITIVE"
	CodeRateRequired   Code = "RATE_REQUIRED"
	CodeRateFormat     Code = "RATE_FORMAT"
	CodeRateRange      Code = "RATE_RANGE"
	CodeTermRequired   Code = "TERM_REQUIRED"
	CodeTermFormat     Code = "TERM_FORMAT"
	CodeTermRange      Code = "TERM_RANGE"
	CodeCalculation    Code = "CALCULATION"
)

type Error struct {
	code    Code
	message string
}

func New(code Code, message string) *Error {
	return &Error{
		code:    code,
		message: message,
	}
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Code() Code {
	return e.code
}

var (
	ErrAmountRequired    = New(CodeAmountRequired, "Введите сумму кредита")
	ErrAmountFormat      = New(CodeAmountFormat, "Сумма должна быть числом, например 250000")
	ErrAmountPositive    = New(CodeAmountPositive, "Сумма кредита должна быть больше нуля")
	ErrRateRequired      = New(CodeRateRequired, "Введите процентную ставку")
	ErrRateFormat        = New(CodeRateFormat, "Ставка должна быть числом, например 14,5")
	ErrRateRange         = New(CodeRateRange, "Ставка должна быть от 0 до 100 процентов")
	ErrTermRequired      = New(CodeTermRequired, "Введите срок кредита")
	ErrTermFormat        = New(CodeTermFormat, "Срок должен быть целым числом")
	ErrTermRange         = New(CodeTermRange, "Срок кредита должен составлять от 1 до 1200 месяцев")
	ErrCalculation       = New(CodeCalculation, "Не удалось рассчитать кредит с указанными параметрами")
	ErrMoreThenMaxAmount = New(CodeCalculation, "Не удалось рассчитать кредит превышан лимит на получение кредита")
)
