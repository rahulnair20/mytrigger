package mytrigger

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

var jsonMetadata = getJSONMetadata()

func getJSONMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

const testConfig string = `{
  "name": "tcmsub",
  "settings": {
    "url": "<Your TCM URL Here>",
		"authkey": "<Your TCM Auth Key Here>",
		"clientid": "flogo-testsubscriber"
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
        "destinationname": "demo_tcm",
				"destinationmatch": "*",
				"messagename": "demo_tcm",
				"durable": "true",
				"durablename": "flogo_demo_tcm"
      }
    }
  ]
}`

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	log.Debugf("Ran Action: %v", uri)
	return 0, nil, nil
}

func (tr *TestRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {
	//log.Debugf("Ran Action: %v", act.Config().Id)
	return nil, nil
}

func TestEndpoint(t *testing.T) {
	log.Info("Testing Endpoint")
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	// New  factory
	f := &tcmsubFactory{}
	f.metadata = trigger.NewMetadata(jsonMetadata)
	tgr := f.New(&config)

	//runner := &TestRunner{}

	//tgr.Init(runner)

	tgr.Start()
	defer tgr.Stop()

	// just loop
	for {
	}
}
