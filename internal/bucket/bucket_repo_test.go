package bucket

import (
	"stc/internal/storage"
	"stc/internal/utils"
	"testing"
)

func TestCreateBucketData(t *testing.T) {
	stcDB := storage.NewInMemStcDB()
	stcDB.Migrate()

	data, err := QueryBucketDataByBucketID(stcDB, 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(data) != 0 {
		t.Fatalf("expected len 0 got %d", len(data))
	}

	tx, err := stcDB.Conn.Begin()
	if err != nil {
		t.Fatal(err)
	}
	if err = CreateBucketDataTx(tx, &BucketData{}); err != nil {
		t.Fatal(err)
	}

	if err = CreateBucketDataTx(tx, &BucketData{
		BucketID:       1,
		ConnectorRunID: 1,
		Raw:            []byte("wusssssssssssssssssarrrrrr"),
		RType:          RAW_TYPE_TEXT,
		CreatedAt:      utils.NewDBTimeNow(),
		Meta:           utils.NewDBMap(map[string]string{"source": "a test"}),
	}); err != nil {
		t.Fatal(err)
	}
	tx.Commit()

	data, err = QueryBucketDataByBucketID(stcDB, 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(data) != 1 {
		t.Fatalf("expected len 2 got %d", len(data))
	}

}
