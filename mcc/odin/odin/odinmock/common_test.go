// odinmock provides an alternative implementation of testing for odin
package odinmock

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/Devex/spaceflight/mcc/odin/odin"
)

type getLastSnapshotCase struct {
	testCase
	name       string
	identifier string
	snapshots  []*rds.DBSnapshot
}

var getLastSnapshotCases = []getLastSnapshotCase{
	// Get snapshot id by instance id
	{
		testCase: testCase{
			expected:      exampleSnapshot1,
			expectedError: "",
		},
		name:       "Get snapshot id by instance id",
		identifier: "production-rds",
		snapshots: []*rds.DBSnapshot{
			exampleSnapshot1,
		},
	},
	// Get non-existing snapshot id by instance id
	{
		testCase: testCase{
			expected:      nil,
			expectedError: "No snapshot found for develop instance",
		},
		name:       "Get non-existing snapshot id by instance id",
		identifier: "develop",
		snapshots:  []*rds.DBSnapshot{},
	},
}

func TestGetLastSnapshot(t *testing.T) {
	svc := newMockRDSClient()
	for _, test := range getLastSnapshotCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				svc.AddSnapshots(test.snapshots)
				actual, err := odin.GetLastSnapshot(
					test.identifier,
					svc,
				)
				test.check(actual, err, t)
			},
		)
	}
}
