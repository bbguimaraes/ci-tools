package modaltesting

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/slack-go/slack"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/ci-tools/pkg/slack/modals"
)

// ValidateBlockIds ensures that all the fields are present as block identifiers in the view
func ValidateBlockIds(t *testing.T, view slack.ModalViewRequest, fields ...string) {
	expected := sets.NewString(fields...)
	actual := sets.NewString()
	for _, block := range view.Blocks.BlockSet {
		// Slack's dynamic marshalling makes this really hard to extract
		blockId := reflect.ValueOf(block).Elem().FieldByName("BlockID").String()
		if blockId != "" {
			actual.Insert(blockId)
		}
	}
	if !expected.Equal(actual) {
		if missing := expected.Difference(actual).List(); len(missing) > 0 {
			t.Errorf("view is missing block IDs: %v", missing)
		}
		if extra := actual.Difference(expected).List(); len(extra) > 0 {
			t.Errorf("view has extra block IDs: %v", extra)
		}
	}
}

type ProcessTestCase struct {
	Name          string
	Callback      []byte
	ExpectedTitle string
	ExpectedBody  string
}

func ValidateParameterProcessing(t *testing.T, parameters modals.JiraIssueParameters, testCases []ProcessTestCase) {
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var callback slack.InteractionCallback
			if err := json.Unmarshal(testCase.Callback, &callback); err != nil {
				t.Errorf("%s: failed to unmarshal payload: %v", testCase.Name, err)
				return
			}
			title, body, err := parameters.Process(&callback)
			if diff := cmp.Diff(testCase.ExpectedTitle, title); diff != "" {
				t.Errorf("%s: got incorrect title: %v", testCase.Name, diff)
			}
			if diff := cmp.Diff(testCase.ExpectedBody, body); diff != "" {
				t.Errorf("%s: got incorrect body: %v", testCase.Name, diff)
			}
			if diff := cmp.Diff(nil, err); diff != "" {
				t.Errorf("%s: got incorrect error: %v", testCase.Name, diff)
			}
		})
	}
}
