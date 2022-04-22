package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/furrygem/nocut-api/internal/links"
	"github.com/furrygem/nocut-api/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

// func (d *db)

// Create creates a new link in the database. Returns (id, isDup, err)
func (d *db) Create(ctx context.Context, link links.Link) (string, bool, error) {
	d.logger.Debug("create link")
	r, err := d.collection.InsertOne(ctx, link)
	if err != nil {
		if IsDup(err) {
			return "", true, fmt.Errorf("Failed to insert link. %v", err)
		}
		return "", false, fmt.Errorf("Failed to insert link. %v", err)
	}
	d.logger.Debug("Convert InsertedID to ObjectID")
	oid, ok := r.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), false, nil
	}
	d.logger.Trace(link)
	d.logger.Debug(err.Error())
	return "", false, fmt.Errorf("Failed to convert ObjectID to hex. oid: '%s'", oid)
}

// func (d *db) FindOne(ctx context.Context, id string) (l links.Link, err error) {
// 	oid, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return l, fmt.Errorf("Failed to convert hex to ObjectID. hex: '%s'", id)
// 	}
// 	filter := bson.M{"_id": oid}

// 	result := d.collection.FindOne(ctx, filter)
// 	if err := result.Err(); err != nil {
// 		// TODO 404
// 		return l, fmt.Errorf("Failed to find link by id '%s'. %v", id, err)
// 	}

// 	if err = result.Decode(&l); err != nil {
// 		return l, fmt.Errorf("Failed to decode link '%s' from DB. %v", id, err)
// 	}
// 	return l, nil
// }

func (d *db) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("Failed to convert hex to ObjecId. hex: '%s'", id)
	}
	filter := bson.M{"_id": oid}
	r, err := d.collection.DeleteOne(ctx, filter)

	if err != nil {
		return fmt.Errorf("Failed to delete link by id '%s'. %v", id, err)
	}

	if r.DeletedCount == 0 {
		return fmt.Errorf("not found")
	}
	d.logger.Tracef("Deleted %d documents", r.DeletedCount)
	return nil
}

// func (d *db) FindOne(ctx context.Context, id string) (links.Link, error) {
// 	oid, err := primitive.ObjectIDFromHex(id)
// 	l := links.Link{}

// 	if err != nil {
// 		return l, fmt.Errorf("Failed to convert hex to ObjectID. hex: '%s'", id)
// 	}
// 	filter := bson.M{"_id": oid}
// 	update := bson.D{{
// 		"$inc", bson.D{{"views", 1}},
// 	}}

// 	r := d.collection.FindOneAndUpdate(ctx, filter, update)
// 	if err := r.Err(); err != nil {
// 		return l, fmt.Errorf("Failed to find and update link by id: %s. %v", id, err)
// 	}
// 	if err := r.Decode(&l); err != nil {
// 		return l, fmt.Errorf("Failed to decode link '%s'. %v", id, err)
// 	}

// 	return l, nil

// }

func (d *db) FindOne(ctx context.Context, id string) (links.Link, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	l := links.Link{}

	if err != nil {
		return l, fmt.Errorf("Failed to convert hex to ObjectID. hex: '%s'", id)
	}
	filter := bson.D{{"_id", oid}}
	project := bson.D{
		{"source", 1},
		{"views", 1},
		{"created_at", 1},
		{"expire_at", 1},
		{"ttl", bson.D{
			{"$subtract", bson.A{"$expire_at", primitive.NewDateTimeFromTime(time.Now())}},
		},
		},
	}
	pipeline := mongo.Pipeline{
		{{"$match", filter}},
		{{"$limit", 1}},
		{{"$project", project}},
	}
	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if cursor.RemainingBatchLength() == 0 {
		return l, fmt.Errorf("Failed to fine link by id '%s'. Batch Length is 0", id)
	}
	cursor.Next(ctx)
	err = cursor.Decode(&l)
	if err != nil {
		d.logger.Errorf("Failed to deocode link by id '%s'. %s", id, err.Error())
		return l, fmt.Errorf("Failed to decode link '%s'. %v", id, err)
	}
	d.logger.Debugf("%v", l)
	d.IncrViews(ctx, id)
	return l, err
}

func (d *db) FindOneBySource(ctx context.Context, source string) (links.Link, error) {
	var l links.Link
	filter := bson.D{
		{"source", source},
	}
	result := d.collection.FindOne(ctx, filter)
	d.logger.Debugf("Source: %s", source)
	if err := result.Err(); err != nil {
		return l, fmt.Errorf("Failed to find link by source '%s'. %v", source, err)
	}
	err := result.Decode(&l)
	if err != nil {
		return l, fmt.Errorf("Failed to decode link by source '%s'. %v", source, err)
	}
	return l, nil
}

func (d *db) IncrViews(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("Failed to convert hex to ObjectID. hex: '%s'", id)
	}
	_, err = d.collection.UpdateOne(
		ctx,
		bson.M{
			"_id": oid,
		},
		bson.D{
			{"$inc", bson.D{{"views", 1}}},
		}, options.Update())

	if err != nil {
		return fmt.Errorf("Failed to increment views count. '%v'", err)
	}

	return nil
}

func (d *db) CreateIndices(ctx context.Context) ([]string, error) {
	expireAtIndex := mongo.IndexModel{
		Keys:    bson.M{"expire_at": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	uniqueSourceIndex := mongo.IndexModel{
		Keys:    bson.M{"source": 1},
		Options: options.Index().SetUnique(true),
	}
	models := []mongo.IndexModel{}
	models = append(models, expireAtIndex)
	models = append(models, uniqueSourceIndex)
	names, err := d.collection.Indexes().CreateMany(ctx, models)
	if err != nil {
		return nil, fmt.Errorf("Failed to create index. %v", err)
	}
	return names, err
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) links.Storage {
	d := &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
	names, err := d.CreateIndices(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, name := range names {
		log.Printf("Cretead index: %s\n", name)
	}
	return d
}
