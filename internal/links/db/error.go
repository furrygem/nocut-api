package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type WriteError struct {
	WriteErrors mongo.WriteErrors
	error
}
