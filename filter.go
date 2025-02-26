package repository_mongo

import "go.mongodb.org/mongo-driver/v2/bson"

type MongoReader interface {
	Bson() bson.M
}
