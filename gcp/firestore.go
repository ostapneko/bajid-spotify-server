package gcp

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"time"
)

type Firestore struct {
	client *firestore.Client
	projectId string
	ctx context.Context
}

func NewFireStore(gcpProjectId string) (*Firestore, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, gcpProjectId)
	if err != nil {
		return nil, fmt.Errorf("failed to setup firestore client: %s", err)
	}

	return &Firestore{client: client, projectId: gcpProjectId, ctx: ctx}, nil
}

/*Data Model
	$userId: serverside-generated
	$listName: chosen by user

	Collection: Songs
	  - Document: $userid-$listName
			- A: spotify-song-url-a
			- B: spotify-song-url-b
			...
			- Z: spotify-song-url-z
*/

func (f *Firestore) GetDocumentRef(userId string, listName string) (*firestore.DocumentRef) {
	docID := fmt.Sprintf("%s-%s", userId, listName)
	return f.client.Collection("songs").Doc(docID)
}

func (f *Firestore) ReadDocument(userId string, listName string) (map[string]string, error) {
	doc := f.GetDocumentRef(userId, listName)
	docSnapshot, err := doc.Get(f.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read doc: %s", err)
	}
	dataMap := docSnapshot.Data()

	mapString := make(map[string]string)

	for key, value := range dataMap {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)

		mapString[strKey] = strValue
	}
	return mapString, nil
}

func (f *Firestore) WriteDocument(userId string, listName string, songList map[string]string) (*firestore.WriteResult, error) {
	doc := f.GetDocumentRef(userId, listName)
	ctx, cancel := context.WithTimeout(f.ctx, 2*time.Second)
	defer cancel()
	// Set either replaces an existing document or creates a new one
	writeResult, err := doc.Set(ctx, songList)

	if err != nil {
		return nil, fmt.Errorf("failed to write doc: %s", err)
	}
	return writeResult, nil
}