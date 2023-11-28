package bucket

import (
	"stc/internal/storage"
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
