package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/apito-io/buffers/protobuff"
	"google.golang.org/api/option"
)

func GetFirestoreDriver(engine *protobuff.DriverCredentials) (*FireStoreDriver, error) {

	// Sets your Google Cloud Platform project ID.
	projectID := engine.ProjectID

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsJSON([]byte(engine.FirebaseProjectCredentialJson)))
	if err != nil {
		return nil, err
	}
	// Close client when done with
	// defer client.Close()

	return &FireStoreDriver{Db: client}, nil
}
