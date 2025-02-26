package repository_mongo

import (
	"github.com/sploitfaze/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoUpdate interface {
	Bson() bson.M
	ReaderOpts() []repository.ReadOption[MongoReader]
}
