package gcp

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"franquel.in/bajidspotifyserver/bajid"
)

type Firestore struct {
	client    *firestore.Client
	projectId string
}

func NewFireStore(gcpProjectId string) (*Firestore, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, gcpProjectId)
	if err != nil {
		return nil, fmt.Errorf("failed to setup firestore client: %s", err)
	}

	return &Firestore{client: client, projectId: gcpProjectId}, nil
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

func (f *Firestore) getDocumentRef(userId bajid.UserId) *firestore.DocumentRef {
	return f.client.Collection("songs").Doc(string(userId))
}

func (f *Firestore) Read(userId bajid.UserId) (bajid.SongList, error) {
	doc := f.getDocumentRef(userId)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	docSnapshot, err := doc.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return make(bajid.SongList), nil
		}

		return nil, fmt.Errorf("failed to read doc: %s", err)
	}
	dataMap := docSnapshot.Data()

	res := make(bajid.SongList, len(dataMap))

	for key, value := range dataMap {
		res[bajid.Letter(key)] = bajid.SpotifyURI(value.(string))
	}

	return res, nil
}

func (f *Firestore) Write(userId bajid.UserId, songList bajid.SongList) error {
	doc := f.getDocumentRef(userId)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// Set either replaces an existing document or creates a new one
	_, err := doc.Set(ctx, songList)

	if err != nil {
		return fmt.Errorf("failed to write doc: %s", err)
	}

	return nil
}
