package presenter

import (
	"context"
	"fmt"
	"micro/pkg/exception"
	"os"

	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/gogo/googleapis/google/rpc"

	"google.golang.org/grpc/status"

	"google.golang.org/grpc/codes"
)

// Error is a struct uses to holds the error information.
type Error struct {
	ctx              context.Context
	Code             codes.Code
	ErrorRPC         exception.ErrorRPCList
	ErrorTracingCode string
	Message          string
	MessageData      map[string]interface{}
	AppName          string
	debug            bool
}

// NewErrorPresenter is a constructor to initialize Error.
func NewErrorPresenter(ctx context.Context, code codes.Code, message string, errors exception.ErrorRPCList) *Error {
	return &Error{
		ctx:      ctx,
		Code:     code,
		ErrorRPC: errors,
		Message:  message,
		AppName:  os.Getenv("APP_GRPC_NAME"),
	}
}

// WithMessageData is a function to include message data in Error.
func (ero *Error) WithMessageData(data map[string]interface{}) *Error {
	ero.MessageData = data

	return ero
}

// WithDetails is a function to include details in Error.
func (ero *Error) WithDetails(errors exception.ErrorRPCList) *Error {
	ero.ErrorRPC = errors

	return ero
}

// WithErrorTracingCode is a function to include error tracing code in Error.
func (ero *Error) WithErrorTracingCode(code string) *Error {
	ero.ErrorTracingCode = code

	return ero
}

// WithDebug is a function to print Error and error details on console.
// Very useful on debugging :).
func (ero *Error) WithDebug() *Error {
	ero.debug = true

	return ero
}

// Build is a function uses to produce error response with details.
func (ero *Error) Build() (*status.Status, error) {
	errStatus := status.New(ero.Code, ero.Message)

	if ero.ErrorRPC != nil && ero.Code == codes.InvalidArgument {
		var errRPCFields []*rpc.BadRequest_FieldViolation

		for _, errorField := range ero.ErrorRPC.ToErrorRPCFieldList() {
			errRPCFields = append(errRPCFields, &rpc.BadRequest_FieldViolation{
				Field:       errorField.Field,
				Description: errorField.Description,
			})
		}

		errDetails := &rpc.BadRequest{
			FieldViolations: errRPCFields,
		}
		errStatus, _ = errStatus.WithDetails(errDetails)
	}

	if ero.ErrorRPC != nil && ero.Code == codes.ResourceExhausted {
		var errRPCFields []*rpc.QuotaFailure_Violation

		for _, errorField := range ero.ErrorRPC.ToErrorRPCFieldList() {
			errRPCFields = append(errRPCFields, &rpc.QuotaFailure_Violation{
				Subject:     errorField.Field,
				Description: errorField.Description,
			})
		}

		errDetails := &rpc.QuotaFailure{
			Violations: errRPCFields,
		}
		errStatus, _ = errStatus.WithDetails(errDetails)
	}

	for _, detail := range errStatus.Details() {
		switch t := detail.(type) {
		case *errdetails.BadRequest:
			fmt.Printf("%s - invalidated fields:\n", ero.AppName)
			for i, v := range t.GetFieldViolations() {
				fmt.Printf("%d) %s - %s\n", i+1, v.Field, v.Description)
			}

		case *errdetails.QuotaFailure:
			fmt.Printf("%s - quota exceeded:\n", ero.AppName)
			for _, v := range t.GetViolations() {
				fmt.Printf("- %s - %s\n", v.Subject, v.Description)
			}
		}
	}

	if ero.ErrorTracingCode != "" {
		info := rpc.ErrorInfo{
			Reason:   ero.ErrorTracingCode,
			Domain:   ero.AppName,
			Metadata: nil,
		}

		errStatus, _ = errStatus.WithDetails(&info)
	}

	return errStatus, nil
}

func (ero *Error) Error() error {
	errRPC, _ := ero.
		WithDetails(ero.ErrorRPC).
		WithDebug().
		Build()

	return errRPC.Err()
}
