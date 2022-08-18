package notionapi

import (
	"context"
	"testing"

	"github.com/psyark/notionapi/appsecret"
)

var client = NewDebugClient(appsecret.NotionAPITestKey)

func TestRetrieveDatabase(t *testing.T) {
	ctx := context.Background()
	_, err := client.RetrieveDatabase(ctx, "8b6685786cc647ecb614dbd9b3ee5113")
	if err != nil {
		t.Fatal(err)
	}
}
func TestRetrievePagePropertyItem(t *testing.T) {
	ctx := context.Background()
	page, err := client.RetrievePage(ctx, "7827e04dd13a4a1682744ec55bd85c56")
	if err != nil {
		t.Fatal(err)
	}

	for _, pv := range page.Properties {
		_, err := client.RetrievePagePropertyItem(ctx, "7827e04dd13a4a1682744ec55bd85c56", pv.ID)
		if err != nil {
			t.Fatal(err)
		}
	}
}
