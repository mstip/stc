package connector

import (
	"fmt"
	"stc/internal/bucket"
	"stc/internal/storage"
	"stc/internal/utils"
)

func StartConnectorRun(stcDB *storage.StcDB, conn Connector) (ConnectorRun, error) {
	connRun, err := CreateConnectorRun(stcDB, conn.ID)
	if err != nil {
		return ConnectorRun{}, err
	}

	switch conn.ConnType {
	case CONECTOR_TYPE_WEBTEXT:
		go func() {
			res, err := WebText(conn.Params.M["url"])
			if err != nil {
				EndConnectorRun(stcDB, connRun.ID, err.Error())
				return
			}
			tx, err := stcDB.Conn.Begin()
			if err != nil {
				EndConnectorRun(stcDB, connRun.ID, err.Error())
				return
			}

			for i := range res {
				err = bucket.CreateBucketDataTx(tx, &bucket.BucketData{
					BucketID:       conn.BucketID,
					ConnectorRunID: connRun.ID,
					Raw:            []byte(res[i].Text),
					RType:          bucket.RAW_TYPE_TEXT,
					Meta: utils.NewDBMap(map[string]string{
						"url":        res[i].URL,
						"path":       res[i].Path,
						"parentType": res[i].ParentType,
					}),
					CreatedAt: utils.NewDBTimeNow(),
				})
				if err != nil {
					tx.Rollback()
					EndConnectorRun(stcDB, connRun.ID, err.Error())
					return
				}
			}

			err = tx.Commit()
			if err != nil {
				EndConnectorRun(stcDB, connRun.ID, err.Error())
				return
			}
			EndConnectorRun(stcDB, connRun.ID, "")
		}()
	default:
		return connRun, fmt.Errorf("unkown connector type %d", conn.ConnType)
	}

	return connRun, nil
}
