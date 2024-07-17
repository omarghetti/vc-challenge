package repo

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
)

var (
	client     *redis.Client
	mockClient redismock.ClientMock
	storer     *RedisStorer
)

func TestMain(m *testing.M) {
	client, mockClient = redismock.NewClientMock()
	storer = New(client)

	m.Run()
}

// Set a document successfully
func TestRedisStorer_SetNewDoc(t *testing.T) {
	ctx := context.Background()

	mockClient.ExpectGet("1").RedisNil()
	mockClient.ExpectSAdd("this", "1").SetVal(1)
	mockClient.ExpectSAdd("test", "1").SetVal(1)
	mockClient.ExpectSAdd("document", "1").SetVal(1)
	mockClient.ExpectSet("1", "This is a test document", 0).SetVal("OK")
	err := storer.SetNewDoc(ctx, "1", "This is a test document")

	if err != nil {
		t.Errorf("Failed to set document: %v", err)
	}

	if err := mockClient.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %v", err)
	}
}

// Get a document successfully
func TestRedisStorer_GetDoc(t *testing.T) {
	ctx := context.Background()

	// Test case 1: Get an existing document
	mockClient.ExpectGet("1").SetVal("This is a test document")
	doc, err := storer.GetDoc(ctx, "1")
	if err != nil {
		t.Errorf("Failed to get document: %v", err)
	}
	if doc != "This is a test document" {
		t.Errorf("Unexpected document content: %s", doc)
	}

	// Test case 2: Get a non-existing document
	mockClient.ExpectGet("3").RedisNil()
	_, err = storer.GetDoc(ctx, "3")
	if err == nil {
		t.Errorf("Expected an error when getting a non-existing document")
	}

	if err := mockClient.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %v", err)
	}
}

func TestRedisStorer_DelDoc(t *testing.T) {
	ctx := context.Background()

	// Test case 1: Delete an existing document
	mockClient.ExpectGet("1").SetVal("This is a test document")
	mockClient.ExpectSRem("this", "1").SetVal(1)
	mockClient.ExpectSRem("test", "1").SetVal(1)
	mockClient.ExpectSRem("document", "1").SetVal(1)
	mockClient.ExpectDel("1").SetVal(1)
	err := storer.DelDoc(ctx, "1")
	if err != nil {
		t.Errorf("Failed to delete document: %v", err)
	}

	// Test case 2: Delete a non-existing document
	mockClient.ExpectGet("3").RedisNil()
	err = storer.DelDoc(ctx, "3")
	if err == nil {
		t.Errorf("Expected an error when deleting a non-existing document")
	}

	if err := mockClient.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %v", err)
	}
}

func TestRedisStorer_SearchDocs(t *testing.T) {
	ctx := context.Background()

	// Test case 1: Search for documents containing all words in the query
	mockClient.ExpectSInter("this", "document").SetVal([]string{"1", "2"})
	docs, err := storer.SearchDocs(ctx, "this,document")

	if err != nil {
		t.Errorf("Failed to search documents: %v", err)
	}

	if len(docs) != 2 {
		t.Errorf("Unexpected number of documents found: %d", len(docs))
	}

	// Test case 2: Search for documents return no results
	mockClient.ExpectSInter("while", "walking").SetVal([]string{})
	docs, err = storer.SearchDocs(ctx, "while,walking")

	if err != nil {
		t.Errorf("Failed to search documents: %v", err)
	}

	if len(docs) != 0 {
		t.Errorf("Unexpected number of documents found: %d", len(docs))
	}

	if err := mockClient.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %v", err)
	}
}
