package domain

import "fmt"

// UseCaseError is an error that can return a UseCase
type UseCaseError struct {
	Code   uint16
	Reason string
	Fix    string
}

func (caseErr UseCaseError) Error() string {
	return fmt.Sprintf("UseCaseError -> [%d] %s: %s", caseErr.Code, caseErr.Reason, caseErr.Fix)
}

var (
	MalformedRequestErr = UseCaseError{Code: 000, Reason: "Bad request", Fix: "See documentation and try again"}
	InternalErr         = UseCaseError{Code: 001, Reason: "Internal Error", Fix: ""}
	UpdateErr           = UseCaseError{Code: 002, Reason: "Error while updating your data", Fix: "Ensure you data is ok"}
)

var (
	SensorNotFoundErr = UseCaseError{Code: 100, Reason: "Sensor not found", Fix: "Create your sensor"}
)
