package field

import "errors"

const defaultError = "TelemetrySDK-GO/Span(Log).Error: "

type StringError string

func (e StringError) Error() string {
	return string(e)
}

const (
	NilPointerError = StringError("nil pointer")
	OverIndexError  = StringError("over index")
)

func GenerateSpecificError(e error) error {
	return errors.New(defaultError + e.Error())
}
