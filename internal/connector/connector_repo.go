package connector

import (
	"stc/internal/storage"
	"stc/internal/utils"
)

type ConnectorType int64

const (
	CONECTOR_TYPE_WEBTEXT ConnectorType = 1
)

type Trigger int64

const (
	TRIGGER_MANUAL Trigger = 1
)

func WebTextConnectorParams(url string) utils.DBMap {
	return utils.NewDBMap(map[string]string{"url": url})
}

type Connector struct {
	ID       int64         `db:"id"`
	BucketID int64         `db:"bucket_id"`
	ConnType ConnectorType `db:"conn_type"`
	Trigger  Trigger       `db:"trigger"`
	Params   utils.DBMap   `db:"params"`
}

func CreateConnector(stcDB *storage.StcDB, bucketID int64, connType ConnectorType, trigger Trigger, params utils.DBMap) (Connector, error) {
	res, err := stcDB.Conn.Exec("INSERT INTO connectors(bucket_id, conn_type, trigger, params) VALUES(?,?,?,?)", bucketID, connType, trigger, params)
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
	ID          int64        `db:"id"`
	ConnectorID int64        `db:"connector_id"`
	Error       string       `db:"error"`
	StartedAt   utils.DBTime `db:"started_at"`
	FinishedAt  utils.DBTime `db:"finished_at"`
}

func CreateConnectorRun(stcDB *storage.StcDB, connectorID int64) (ConnectorRun, error) {
	now := utils.NewDBTimeNow()
	res, err := stcDB.Conn.Exec("INSERT INTO connector_runs(connector_id, started_at, error) VALUES(?,?,?)", connectorID, now.String(), "")
	if err != nil {
		return ConnectorRun{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return ConnectorRun{}, err
	}
	return ConnectorRun{ID: id, ConnectorID: connectorID, StartedAt: now}, nil
}

func EndConnectorRun(stcDB *storage.StcDB, connRunID int64, errStr string) error {
	now := utils.NewDBTimeNow()
	_, err := stcDB.Conn.Exec("UPDATE connector_runs SET finished_at = ?, error = ? WHERE id = ?", now.String(), errStr, connRunID)
	return err
}

func QueryConnectorRunByID(stcDB *storage.StcDB, connRunID int64) (ConnectorRun, error) {
	var connRunn ConnectorRun
	err := stcDB.Conn.Get(&connRunn, "SELECT id,connector_id, error, started_at, finished_at FROM connector_runs WHERE id = ?", connRunID)
	if err != nil {
		return ConnectorRun{}, err
	}
	return connRunn, nil
}
