package db

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

// IsDup checks if MongoDB error has error code 11000 (duplicate)
func IsDup(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}
