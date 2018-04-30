package mytrigger

import (
	"context"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	"github.com/Sirupsen/logrus"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	//"flag"
)

// log is the default package logger
var log = logger.GetLogger("trigger-jvanderl-tcmsub")

// tcmsubTrigger is a stub for your Trigger implementation
type tcmsubTrigger struct {
	metadata              *trigger.Metadata
	runner                action.Runner
	config                *trigger.Config
	destinationToActionId map[string]string
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &tcmsubFactory{metadata: md}
}

// tcmsubFactory Trigger factory
type tcmsubFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *tcmsubFactory) New(config *trigger.Config) trigger.Trigger {
	tcmsubTrigger := &tcmsubTrigger{metadata: t.metadata, config: config}
	return tcmsubTrigger
}

// Metadata implements trigger.Trigger.Metadata
func (t *tcmsubTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *tcmsubTrigger) Init(runner action.Runner) {
	t.runner = runner
}

// Start implements trigger.Trigger.Start
func (t *tcmsubTrigger) Start() error {

	// start the trigger
	consKey := t.config.GetSetting("consumerKey")
	consSec := t.config.GetSetting("consumerSecret")
	accTok := t.config.GetSetting("accessToken")
	accTokSec := t.config.GetSetting("accessTokenSecret")

	// Read Actions from trigger endpoints
	t.destinationToActionId = make(map[string]string)

	for _, handlerCfg := range t.config.Handlers {
		log.Debugf("handlers: [%s]", handlerCfg.ActionId)
		epdestination := handlerCfg.GetSetting("destinationname")
		log.Debugf("destination: [%s]", epdestination)
		t.destinationToActionId[epdestination] = handlerCfg.ActionId
	}

	anaconda.SetConsumerKey(consKey)
	anaconda.SetConsumerSecret(consSec)
	api := anaconda.NewTwitterApi(accTok, accTokSec)

	stream := api.PublicStreamFilter(url.Values{
		"track": []string{"TESTGOLANGTIBCO"},
	})

	defer stream.Stop()

	for v := range stream.C {
		t, ok := v.(anaconda.Tweet)
		if !ok {
			logrus.Warningf("received unexpected value of type %T", v)
			continue
		}

		if t.RetweetedStatus != nil {
			continue
		}

		_, err := api.Retweet(t.Id, false)
		if err != nil {
			logrus.Errorf("could not retweet %d: %v", t.Id, err)
			continue
		}
		logrus.Infof("retweeted %d", t.Id)

	}
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *tcmsubTrigger) Stop() error {
	// stop the trigger
	return nil
}

// RunAction starts a new Process Instance
func (t *tcmsubTrigger) RunAction(actionId string, payload string, destination string) {
	log.Debug("Starting new Process Instance")
	log.Debugf("Action Id: %s", actionId)
	log.Debugf("Payload: %s", payload)

	req := t.constructStartRequest(payload)

	startAttrs, _ := t.metadata.OutputsToAttrs(req.Data, false)

	action := action.Get(actionId)

	context := trigger.NewContext(context.Background(), startAttrs)

	_, replyData, err := t.runner.Run(context, action, actionId, nil)
	if err != nil {
		log.Error(err)
	}

	log.Debugf("Ran action: [%s]", actionId)
	log.Debugf("Reply data: [%v]", replyData)

}

func (t *tcmsubTrigger) constructStartRequest(message string) *StartRequest {

	//TODO how to handle reply to, reply feature
	req := &StartRequest{}
	data := make(map[string]interface{})
	data["message"] = message
	req.Data = data
	return req
}

// StartRequest describes a request for starting a ProcessInstance
type StartRequest struct {
	ProcessURI string                 `json:"flowUri"`
	Data       map[string]interface{} `json:"data"`
	ReplyTo    string                 `json:"replyTo"`
}

func convert(b []byte) string {
	n := len(b)
	return string(b[:n])
}