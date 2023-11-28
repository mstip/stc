package connector

import (
	"fmt"
	"stc/internal/storage"
)

func StartConnectorRun(stcDB *storage.StcDB, conn Connector) (ConnectorRun, error) {
	connRun, err := CreateConnectorRun(stcDB, conn.ID)
	if err != nil {
		return ConnectorRun{}, err
	}

	switch conn.ConnType {
	case CONECTOR_TYPE_WEBTEXT:
		go func() {
			// TODO: continue here
			WebText(conn.Params["url"])
		}()
	default:
		return connRun, fmt.Errorf("unkown connector type %d", conn.ConnType)
	}

	return connRun, nil
}
