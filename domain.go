package repository_mongo

import "go.mongodb.org/mongo-driver/v2/bson"

type MongoDomain interface {
	SetId(id string) string
	Bson() bson.M
	FromBson(bson.M) error
}
