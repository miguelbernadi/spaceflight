package odinmock

import (
	"testing"

	"github.com/go-test/deep"
)

type testCase struct {
	expected      interface{}
	expectedError string
}

func (tc *testCase) expectingError(err error) bool {
	return tc.expectedError != "" && err.Error() == tc.expectedError
}

func (tc *testCase) check(actual interface{}, err error, t *testing.T) {
	switch {
	case err != nil && !tc.expectingError(err):
		t.Errorf(
			"Unexpected error: %v",
			err,
		)
	case err != nil && tc.expectingError(err):
	case err == nil && tc.expectedError != "":
		t.Errorf(
			"Expected error: %v missing",
			tc.expectedError,
		)
	case err == nil:
		if diff := deep.Equal(
			actual,
			tc.expected,
		); diff != nil {
			t.Errorf(
				"Unexpected output: %s",
				diff,
			)
		}
	}
}
