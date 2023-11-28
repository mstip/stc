package int

import (
	"reflect"
	"stc/internal/bucket"
	"stc/internal/connector"
	"stc/internal/storage"
	"testing"
	"time"
)

// TestFEM is the integration test for fem
// 1. Create Bucket for page
// 2. Add webtext Connector to bucket with url configure to gather manual
// 3. Run connector fill bucket
// 4. Create Store, Type AI with Text Vector and meta data url
// 5. Add Pipeline to Store run manualy map webtext text to text and url to meta data url
// 6. Run Pipeline fille store
// 7. Create Query with user text search input, add sim search on store and configure chat comp
// 8. Expose Query as Rest and Web with low code web
func TestFEM(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping int tests in short mode")
	}

	expBuc := bucket.Bucket{ID: 1, Name: "fem_raise_bucket"}
	expConn := connector.Connector{ID: 1, ConnType: connector.CONECTOR_TYPE_WEBTEXT, BucketID: 1, Trigger: connector.TRIGGER_MANUAL, Params: connector.WebTextConnectorParams("https://raise-ai-mindsets.de/")}
	expConnRun := connector.ConnectorRun{ID: 1, ConnectorID: 1, Error: "", StartedAt: time.Now()}
	stcDB := storage.NewInMemStcDB()
	stcDB.Migrate()

	// 1. Create Bucket for page
	buc, err := bucket.CreateBucket(stcDB, expBuc.Name)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expBuc, buc) {
		t.Fatalf("CreateBucket expected %#v, got %#v", expBuc, buc)
	}

	// 2. Add webtext Connector to bucket with url configure to gather manual
	conn, err := connector.CreateConnector(stcDB, buc.ID, connector.CONECTOR_TYPE_WEBTEXT, connector.TRIGGER_MANUAL, connector.WebTextConnectorParams(expConn.Params["url"]))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expConn, conn) {
		t.Fatalf("CreateConnector expected %#v, got %#v", expConn, conn)
	}

	// 3. Run connector fill bucket
	connRun, err := connector.StartConnectorRun(stcDB, conn)
	if err != nil {
		t.Fatal(err)
	}

	if expConnRun.ID != connRun.ID || expConnRun.ConnectorID != connRun.ConnectorID || expConnRun.Error != connRun.Error || connRun.StartedAt.Before(expConnRun.StartedAt) {
		t.Fatalf("RunConnector expected %#v, got %#v", expConnRun, connRun)
	}

	// TODO: check in loop if connrun is done or died with an error timeout 5 min
}
