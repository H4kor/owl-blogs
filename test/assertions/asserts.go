package assertions

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func Assert(t *testing.T, condition bool, message string) {
	t.Helper()
	if !condition {
		t.Errorf(message)
	}
}

func AssertNot(t *testing.T, condition bool, message string) {
	t.Helper()
	if condition {
		t.Errorf(message)
	}
}

func AssertContains(t *testing.T, containing string, search string) {
	t.Helper()
	if !strings.Contains(containing, search) {
		t.Errorf("Expected '%s' to contain '%s'", containing, search)
	}
}

func AssertArrayContains[T comparable](t *testing.T, list []T, search T) {
	t.Helper()
	for _, item := range list {
		if item == search {
			return
		}
	}
	t.Errorf("Expected '%v' to be in '%v'", search, list)
}

func AssertNotContains(t *testing.T, containing string, search string) {
	t.Helper()
	if strings.Contains(containing, search) {
		t.Errorf("Expected '%s' to not contain '%s'", containing, search)
	}
}

func AssertNoError(t *testing.T, err error, message string) {
	t.Helper()
	if err != nil {
		t.Errorf(message+": %s", err.Error())
	}
}

func AssertError(t *testing.T, err error, message string) {
	t.Helper()
	if err == nil {
		t.Errorf(message)
	}
}

func AssertLen[T any](t *testing.T, list []T, expected int) {
	t.Helper()
	if len(list) != expected {
		t.Errorf("Expected list to have length %d, got %d", expected, len(list))
	}
}

func AssertMapLen[T any, S comparable](t *testing.T, list map[S]T, expected int) {
	t.Helper()
	if len(list) != expected {
		t.Errorf("Expected list to have length %d, got %d", expected, len(list))
	}
}

func AssertEqual[T comparable](t *testing.T, actual T, expected T) {
	t.Helper()
	if actual != expected {
		t.Errorf("Expected '%v', got '%v'", expected, actual)
	}
}

func AssertNotEqual[T comparable](t *testing.T, actual T, expected T) {
	t.Helper()
	if actual == expected {
		t.Errorf("Expected '%v' to not be '%v'", expected, actual)
	}
}

func AssertStatus(t *testing.T, rr *httptest.ResponseRecorder, expStatus int) {
	if status := rr.Code; status != expStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expStatus)
		return
	}
}

func AssertLessThan(t *testing.T, actual int, expected int) {
	t.Helper()
	if actual >= expected {
		t.Errorf("Expected '%d' to be less than '%d'", actual, expected)
	}
}

func AssertGreaterThan(t *testing.T, actual int, expected int) {
	t.Helper()
	if actual <= expected {
		t.Errorf("Expected '%d' to be greater than '%d'", actual, expected)
	}
}
