package odinmock

import (
	"testing"
	"time"

	"github.com/poka-yoke/spaceflight/mcc/odin/odin"
)

type createInstanceCase struct {
	testCase
	name         string
	identifier   string
	instanceType string
	password     string
	user         string
	size         int64
}

var createInstanceCases = []createInstanceCase{
	// Creating simple instance
	{
		testCase: testCase{
			expected:      "test1.0.us-east-1.rds.amazonaws.com",
			expectedError: "",
		},
		name:         "Creating simple instance",
		identifier:   "test1",
		instanceType: "db.m1.small",
		user:         "master",
		password:     "master",
		size:         5,
	},
	// Fail because empty user
	{
		testCase: testCase{
			expected:      "",
			expectedError: "Specify Master User",
		},
		name:         "Fail because empty user",
		identifier:   "test1",
		instanceType: "db.m1.small",
		user:         "",
		password:     "master",
		size:         5,
	},
	// Fail because empty password
	{
		testCase: testCase{
			expected:      "",
			expectedError: "Specify Master User Password",
		},
		name:         "Fail because empty password",
		identifier:   "test1",
		instanceType: "db.m1.small",
		user:         "master",
		password:     "",
		size:         5,
	},
	// Fail because non-present size
	{
		testCase: testCase{
			expected:      "",
			expectedError: "Specify size between 5 and 6144",
		},
		name:         "Fail because non-present size",
		identifier:   "test1",
		instanceType: "db.m1.small",
		user:         "master",
		password:     "master",
	},
	// Fail because too small size
	{
		testCase: testCase{
			expected:      "",
			expectedError: "Specify size between 5 and 6144",
		},
		name:         "Fail because too small size",
		identifier:   "test1",
		instanceType: "db.m1.small",
		user:         "master",
		password:     "master",
	},
	// Fail because too big size
	{
		testCase: testCase{
			expected:      "",
			expectedError: "Specify size between 5 and 6144",
		},
		name:         "Fail because too big size",
		identifier:   "test1",
		instanceType: "db.m1.small",
		user:         "master",
		password:     "master",
		size:         6145,
	},
}

func TestCreateInstance(t *testing.T) {
	svc := newMockRDSClient()
	odin.Duration = time.Duration(0)
	for _, test := range createInstanceCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				params := odin.CreateParams{
					InstanceType: test.instanceType,
					User:         test.user,
					Password:     test.password,
					Size:         test.size,
				}
				actual, err := odin.CreateInstance(
					test.identifier,
					params,
					svc,
				)
				test.check(actual, err, t)
			},
		)
	}
}

type createInstanceCaseInput struct {
	name         string
	identifier   string
	instanceType string
	password     string
	user         string
	size         int64
}

type createInstanceCaseOutput struct {
	out string
	err error
}

var createInstanceCaseMap = map[string]struct {
	Input    createInstanceCaseInput
	Output   testCase
	Internal map[string]interface{}
}{
	"CreateExistingInstanceFail": {
		createInstanceCaseInput{
			name:         "CreateExistingInstanceFail",
			identifier:   "CreateExistingInstanceFail",
			instanceType: "db.m1.small",
			user:         "master",
			password:     "master",
			size:         5,
		},
		testCase{
			expected:      "",
			expectedError: "Instance CreateExistingInstanceFail already exists",
		},
		map[string]interface{}{
			"CreateDBInstance": nil,
		},
	},
}

func TestCreateInstanceWithMap(t *testing.T) {
	svc := newMockRDSClient()
	odin.Duration = time.Duration(0)
	for _, test := range createInstanceCaseMap {
		t.Run(
			test.Input.name,
			func(t *testing.T) {
				params := odin.CreateParams{
					InstanceType: test.Input.instanceType,
					User:         test.Input.user,
					Password:     test.Input.password,
					Size:         test.Input.size,
				}
				actual, err := odin.CreateInstance(
					test.Input.identifier,
					params,
					svc,
				)
				test.Output.check(actual, err, t)
			},
		)
	}
}
