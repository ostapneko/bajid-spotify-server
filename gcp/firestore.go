package gcp

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"

	"franquel.in/bajidspotifyserver/bajid"
)

type Firestore struct {
	client    *firestore.Client
	projectId string
	ctx       context.Context
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

Collection: Songs
  - Document: $userid
		- A: spotify-song-url-a
		- B: spotify-song-url-b
		...
		- Z: spotify-song-url-z
*/

func (f *Firestore) GetDocumentRef(userId bajid.UserId) *firestore.DocumentRef {
	return f.client.Collection("songs").Doc(string(userId))
}

func (f *Firestore) ReadDocument(userId bajid.UserId) (bajid.LetterToSong, error) {
	doc := f.GetDocumentRef(userId)
	docSnapshot, err := doc.Get(f.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read doc: %s", err)
	}
	dataMap := docSnapshot.Data()

	res := make(bajid.LetterToSong, len(dataMap))

	for key, value := range dataMap {
		res[bajid.Letter(key)] = value.(bajid.SpotifyURI)
	}

	return res, nil
}

func (f *Firestore) WriteDocument(userId bajid.UserId, songList bajid.LetterToSong) (*firestore.WriteResult, error) {
	doc := f.GetDocumentRef(userId)
	ctx, cancel := context.WithTimeout(f.ctx, 2*time.Second)
	defer cancel()
	// Set either replaces an existing document or creates a new one
	writeResult, err := doc.Set(ctx, songList)

	if err != nil {
		return nil, fmt.Errorf("failed to write doc: %s", err)
	}
	return writeResult, nil
}
