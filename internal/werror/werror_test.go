package werror_test

import (
	"errors"
	"testing"

	"github.com/brkcnr/getweatherapi/internal/werror"
)

func TestNew(t *testing.T) {
	err := werror.New("some error", true, 400)
	if err.Error() != "some error" {
		t.Errorf("expected message `some error`, got `%s`", err.Error())
	}
	if !err.(*werror.Error).Loggable {
		t.Errorf("expected loggable to be true, got false")
	}
	if err.Code() != 400 {
		t.Errorf("expected status code 400, got %d", err.Code())
	}
}

func TestSentinelErrors(t *testing.T) {
	if werror.ErrInvalidCity.Error() != "Invalid city name. Please try again." {
		t.Errorf("unexpected message for ErrInvalidCity")
	}
	if werror.ErrInvalidCity.Code() != 400 {
		t.Errorf("unexpected status code for ErrIncvalidCity")
	}
	if werror.ErrInvalidCity.(*werror.Error).Loggable {
		t.Errorf("expected ErrInvalidCity loggable to be false, got true")
	}
}

func TestWrapAndUnwrap(t *testing.T) {
	baseErr := errors.New("base error")
	err := werror.New("wrapped error", true, 500).Wrap(baseErr)
	if err.Unwrap() != baseErr {
		t.Errorf("expected unwrapped error to be base error, got %v", err.Unwrap())
	}
	expectedMsg := "base error, wrapped error"
	if err.Error() != expectedMsg {
		t.Errorf("expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestAddDataAndCleanData(t *testing.T) {
	err := werror.New("data error", true, 500)
	err.AddData("extra data")
	if err.(*werror.Error).Data != "extra data" {
		t.Errorf("expected data to be 'extra data', got %v", err.(*werror.Error).Data)
	}
	err.ClearData()
	if err.(*werror.Error).Data != nil {
		t.Errorf("expected data to be nil after ClearData, got %v", err.(*werror.Error).Data)
	}
}

func TestCode(t *testing.T) {
	err := werror.New("code error", true, 403)
	if err.Code() != 403 {
		t.Errorf("expected code 403, got %d", err.Code())
	}
}
