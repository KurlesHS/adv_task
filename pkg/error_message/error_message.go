package error_message

type ErrorType int

const (
	Db ErrorType = iota
	BadRequest
	NotFound
	Internal
)

type ErrorMessage struct {
	msg   string
	eType ErrorType
}

func NewError(eType ErrorType, msg string) ErrorMessage {
	return ErrorMessage{
		msg:   msg,
		eType: eType,
	}
}

func (e ErrorMessage) Error() string {
	return e.msg
}

func (e ErrorMessage) Type() ErrorType {
	return e.eType
}
