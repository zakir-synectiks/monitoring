package alerting

import (
	"sync"
	"time"

	"github.com/xformation/synectiks-monitoring/pkg/bus"
	"github.com/xformation/synectiks-monitoring/pkg/log"
	"github.com/xformation/synectiks-monitoring/pkg/metrics"
	m "github.com/xformation/synectiks-monitoring/pkg/models"
)

type RuleReader interface {
	Fetch() []*Rule
}

type DefaultRuleReader struct {
	sync.RWMutex
	//serverID       string
	serverPosition int
	clusterSize    int
	log            log.Logger
}

func NewRuleReader() *DefaultRuleReader {
	ruleReader := &DefaultRuleReader{
		log: log.New("alerting.ruleReader"),
	}

	go ruleReader.initReader()
	return ruleReader
}

func (arr *DefaultRuleReader) initReader() {
	heartbeat := time.NewTicker(time.Second * 10)

	for range heartbeat.C {
		arr.heartbeat()
	}
}

func (arr *DefaultRuleReader) Fetch() []*Rule {
	cmd := &m.GetAllAlertsQuery{}

	if err := bus.Dispatch(cmd); err != nil {
		arr.log.Error("Could not load alerts", "error", err)
		return []*Rule{}
	}

	res := make([]*Rule, 0)
	for _, ruleDef := range cmd.Result {
		if model, err := NewRuleFromDBAlert(ruleDef); err != nil {
			arr.log.Error("Could not build alert model for rule", "ruleId", ruleDef.Id, "error", err)
		} else {
			res = append(res, model)
		}
	}

	metrics.M_Alerting_Active_Alerts.Set(float64(len(res)))
	return res
}

func (arr *DefaultRuleReader) heartbeat() {
	arr.clusterSize = 1
	arr.serverPosition = 1
}
