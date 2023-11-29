package connector

import (
	"stc/internal/storage"
	"stc/internal/utils"
	"testing"
)

func TestConnectorRunRepo(t *testing.T) {
	stcDB := storage.NewInMemStcDB()
	stcDB.Migrate()

	expConnRun := ConnectorRun{ID: 1, ConnectorID: 1, Error: "", StartedAt: utils.NewDBTimeNow()}

	connRun, err := CreateConnectorRun(stcDB, 1)
	if err != nil {
		t.Fatal(err)
	}

	if expConnRun.ID != connRun.ID || expConnRun.ConnectorID != connRun.ConnectorID || expConnRun.Error != connRun.Error || connRun.StartedAt.T.Before(expConnRun.StartedAt.T) {
		t.Fatalf("CreateConnectorRun expected %#v, got %#v", expConnRun, connRun)
	}

	qcRun, err := QueryConnectorRunByID(stcDB, connRun.ID)
	if err != nil {
		t.Fatal(err)
	}

	if connRun.ID != qcRun.ID || connRun.ConnectorID != qcRun.ConnectorID || connRun.Error != qcRun.Error || qcRun.StartedAt.T.Before(connRun.StartedAt.T) {
		t.Fatalf("QueryConnectorRunByID expected %#v, got %#v", connRun, qcRun)
	}

	err = EndConnectorRun(stcDB, connRun.ID, "")
	if err != nil {
		t.Fatal(err)
	}

	qcRun, err = QueryConnectorRunByID(stcDB, connRun.ID)
	if err != nil {
		t.Fatal(err)
	}
	if connRun.ID != qcRun.ID || connRun.ConnectorID != qcRun.ConnectorID || qcRun.Error != "" || qcRun.StartedAt.T.Before(connRun.StartedAt.T) || qcRun.FinishedAt.T.Before(connRun.StartedAt.T) {
		t.Fatalf("QueryConnectorRunByID expected %#v, got %#v", connRun, qcRun)
	}
}
