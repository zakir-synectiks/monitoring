package sqlstore

import (
	"github.com/xformation/synectiks-monitoring/pkg/bus"
	m "github.com/xformation/synectiks-monitoring/pkg/models"
)

func init() {
	bus.AddHandler("sql", GetDBHealthQuery)
}

func GetDBHealthQuery(query *m.GetDBHealthQuery) error {
	return x.Ping()
}
