package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, host string, port uint16, username, password, database, authDB string) (*mongo.Database, error) {
	var (
		mongoDBURL string
		isAuth     bool
	)
	if username == "" && password == "" {
		mongoDBURL = fmt.Sprintf("mongodb://%s:%d", host, port)
	} else {
		isAuth = true
		mongoDBURL = fmt.Sprintf("mongodb://%s:%s@%s:%d", username, password, host, port)
	}

	clientOptions := options.Client().ApplyURI(mongoDBURL)
	if isAuth {
		if authDB == "" {
			authDB = database
		}
		clientOptions.SetAuth(options.Credential{
			AuthSource: authDB,
			Username:   username,
			Password:   password,
		})
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb. %v", err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("Failed to ping mongoDB. %v", err)
	}

	return client.Database(database), nil
}
