package notionapi

import (
	"context"
	"testing"

	"github.com/psyark/notionapi/appsecret"
)

var client = NewDebugClient(appsecret.NotionAPITestKey)

func TestRetrievePage(t *testing.T) {
	ctx := context.Background()
	_, err := client.RetrievePage(ctx, "22a5412dd0ab4167930cb644d11fffea")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRetrieveDatabase(t *testing.T) {
	ctx := context.Background()
	_, err := client.RetrieveDatabase(ctx, "8b6685786cc647ecb614dbd9b3ee5113")
	if err != nil {
		t.Fatal(err)
	}
}
