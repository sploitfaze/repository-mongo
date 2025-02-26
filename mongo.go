package repository_mongo

import (
	"context"
	"github.com/sploitfaze/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoRepository[D MongoDomain, R MongoReader, U MongoUpdate] struct {
	collection *mongo.Collection
	reader     func() R
	updater    func() U
	domain     func() D
}

func (r *MongoRepository[D, R, U]) Create(ctx context.Context, domain D) (*D, error) {
	res, err := r.collection.InsertOne(ctx, domain)
	if err != nil {
		return nil, err
	}
	domain.SetId(res.InsertedID.(string))
	return &domain, nil
}

func (r *MongoRepository[D, R, U]) Read(ctx context.Context, opts ...repository.ReadOption[R]) ([]D, error) {
	reader := r.reader()

	for _, opt := range opts {
		if err := opt(reader); err != nil {
			return nil, err
		}
	}

	cursor, err := r.collection.Find(ctx, reader.Bson())
	if err != nil {
		return nil, err
	}

	var domains []D
	for cursor.Next(ctx) {
		decodedBson := bson.M{}
		domain := r.domain()

		if err := cursor.Decode(&decodedBson); err != nil {
			return nil, err
		}

		if err := domain.FromBson(decodedBson); err != nil {
			return nil, err
		}

		domains = append(domains, domain)
	}

	return domains, nil
}

func (r *MongoRepository[D, R, U]) Update(ctx context.Context, opts ...repository.UpdateOption[U]) error {
	reader := r.reader()
	updater := r.updater()

	for _, opt := range opts {
		if err := opt(updater); err != nil {
			return err
		}
	}

	for _, opts := range updater.ReaderOpts() {
		if err := opts(reader); err != nil {
			return err
		}
	}

	_, err := r.collection.UpdateOne(ctx, reader.Bson(), updater.Bson())
	return err
}

func (r *MongoRepository[D, R, U]) Delete(ctx context.Context, opts ...repository.ReadOption[R]) error {
	reader := r.reader()

	for _, opt := range opts {
		if err := opt(reader); err != nil {
			return err
		}
	}

	_, err := r.collection.DeleteOne(ctx, reader.Bson())
	return err
}

func NewRepositoryMongo[D MongoDomain, R MongoReader, U MongoUpdate](
	collection *mongo.Collection,
) *MongoRepository[D, R, U] {
	return &MongoRepository[D, R, U]{
		collection: collection,
	}
}
