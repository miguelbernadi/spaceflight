package odinmock

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/poka-yoke/spaceflight/mcc/odin/odin"
)

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
	"CreateSimpleInstanceSuccess": {
		createInstanceCaseInput{
			name:         "CreateSimpleInstanceSuccess",
			identifier:   "CreateSimpleInstanceSuccess",
			instanceType: "db.m1.small",
			user:         "master",
			password:     "master",
			size:         5,
		},
		testCase{
			expected:      "CreateSimpleInstanceSuccess.0.us-east-1.rds.amazonaws.com",
			expectedError: "",
		},
		map[string]interface{}{
			"CreateDBInstance": &rds.CreateDBInstanceOutput{
				DBInstance: &rds.DBInstance{
					AllocatedStorage:     aws.Int64(5),
					AvailabilityZone:     aws.String("us-east-1"),
					DBInstanceIdentifier: aws.String("CreateSimpleInstanceSuccess"),
					DBInstanceStatus:     aws.String("creating"),
					Engine:               aws.String(""),
				},
			},
			"DescribeDBInstances": &rds.DescribeDBInstancesOutput{
				DBInstances: []*rds.DBInstance{
					{
						DBInstanceStatus: aws.String("creating"),
						Endpoint: &rds.Endpoint{
							Address: aws.String("CreateSimpleInstanceSuccess.0.us-east-1.rds.amazonaws.com"),
							Port:    aws.Int64(5432),
						},
					},
				},
			},
		},
	},
	"CreateEmptyUserFail": {
		createInstanceCaseInput{
			name:         "CreateEmptyUserFail",
			identifier:   "CreateEmptyUserFail",
			instanceType: "db.m1.small",
			user:         "",
			password:     "master",
			size:         5,
		},
		testCase{
			expected:      "",
			expectedError: "Specify Master User",
		},
		map[string]interface{}{
			"CreateDBInstance": nil,
		},
	},
	"CreateEmptyPasswordFail": {
		createInstanceCaseInput{
			name:         "CreateEmptyPasswordFail",
			identifier:   "CreateEmptyPasswordFail",
			instanceType: "db.m1.small",
			user:         "master",
			password:     "",
			size:         5,
		},
		testCase{
			expected:      "",
			expectedError: "Specify Master User Password",
		},
		map[string]interface{}{
			"CreateDBInstance": nil,
		},
	},
	"CreateAbsentSizeFail": {
		createInstanceCaseInput{
			name:         "CreateAbsentSizeFail",
			identifier:   "CreateAbsentSizeFail",
			instanceType: "db.m1.small",
			user:         "master",
			password:     "master",
		},
		testCase{
			expected:      "",
			expectedError: "Specify size between 5 and 6144",
		},
		map[string]interface{}{
			"CreateDBInstance": nil,
		},
	},
	"CreateTooSmallFail": {
		createInstanceCaseInput{
			name:         "CreateTooSmallFail",
			identifier:   "CreateTooSmallFail",
			instanceType: "db.m1.small",
			user:         "master",
			password:     "master",
		},
		testCase{
			expected:      "",
			expectedError: "Specify size between 5 and 6144",
		},
		map[string]interface{}{
			"CreateDBInstance": nil,
		},
	},
	"CreateTooBigFail": {
		createInstanceCaseInput{
			name:         "CreateTooBigFail",
			identifier:   "CreateTooBigFail",
			instanceType: "db.m1.small",
			user:         "master",
			password:     "master",
		},
		testCase{
			expected:      "",
			expectedError: "Specify size between 5 and 6144",
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
