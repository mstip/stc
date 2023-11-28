package connector

import (
	"encoding/json"
	"stc/internal/storage"
	"stc/internal/utils"
	"time"
)

type ConnectorType int64

const (
	CONECTOR_TYPE_WEBTEXT ConnectorType = 1
)

type Trigger int64

const (
	TRIGGER_MANUAL Trigger = 1
)

func WebTextConnectorParams(url string) map[string]string {
	return map[string]string{"url": url}
}

type Connector struct {
	ID       int64             `db:"id"`
	BucketID int64             `db:"bucket_id"`
	ConnType ConnectorType     `db:"conn_type"`
	Trigger  Trigger           `db:"trigger"`
	Params   map[string]string `db:"params"`
}

func CreateConnector(stcDB *storage.StcDB, bucketID int64, connType ConnectorType, trigger Trigger, params map[string]string) (Connector, error) {
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return Connector{}, err
	}

	res, err := stcDB.Conn.Exec("INSERT INTO connectors(bucket_id, conn_type, trigger, params) VALUES(?,?,?,?)", bucketID, connType, trigger, jsonParams)
	if err != nil {
		return Connector{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Connector{}, err
	}

	return Connector{ID: id, BucketID: bucketID, ConnType: connType, Trigger: trigger, Params: params}, nil
}

type ConnectorRun struct {
	ID          int64     `db:"id"`
	ConnectorID int64     `db:"connector_id"`
	Error       string    `db:"error"`
	StartedAt   time.Time `db:"started_at"`
	FinishedAt  time.Time `db:"finished_at"`
}

func CreateConnectorRun(stcDB *storage.StcDB, connectorID int64) (ConnectorRun, error) {
	now := time.Now()
	res, err := stcDB.Conn.Exec("INSERT INTO connector_runs(connector_id, started_at) VALUES(?,?)", connectorID, utils.TimeToStr(now))
	if err != nil {
		return ConnectorRun{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return ConnectorRun{}, err
	}
	return ConnectorRun{ID: id, ConnectorID: connectorID, StartedAt: now}, nil
}
