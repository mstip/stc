package bucket

import (
	"database/sql"
	"stc/internal/storage"
	"stc/internal/utils"
)

type Bucket struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func CreateBucket(stcDB *storage.StcDB, name string) (Bucket, error) {
	res, err := stcDB.Conn.Exec("INSERT INTO buckets(name) VALUES(?)", name)
	if err != nil {
		return Bucket{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Bucket{}, err
	}
	return Bucket{ID: id, Name: name}, nil
}

type RawType int64

const (
	RAW_TYPE_TEXT   RawType = 1
	RAW_TYPE_BINARY RawType = 2
)

type BucketData struct {
	BucketID       int64        `db:"bucket_id"`
	ConnectorRunID int64        `db:"connector_run_id"`
	Raw            []byte       `db:"raw"`
	RType          RawType      `db:"raw_type"`
	Meta           utils.DBMap  `db:"meta"`
	CreatedAt      utils.DBTime `db:"created_at"`
}

func CreateBucketDataTx(tx *sql.Tx, data *BucketData) error {
	_, err := tx.Exec(
		"INSERT INTO bucket_data(bucket_id, connector_run_id, raw, raw_type, meta, created_at) VALUES(?,?,?,?,?,?)",
		data.BucketID, data.ConnectorRunID, data.Raw, data.RType, data.Meta, data.CreatedAt.String(),
	)
	if err != nil {
		return err
	}
	return err
}

func QueryBucketDataByBucketID(stcDB *storage.StcDB, bucketID int64) ([]BucketData, error) {
	var data []BucketData
	err := stcDB.Conn.Select(&data, "SELECT * FROM bucket_data WHERE bucket_id = ?", bucketID)
	if err != nil {
		return nil, err
	}
	return data, nil
}
