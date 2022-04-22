package links

import (
	"time"
)

type Link struct {
	ID              string        `json:"id,omitempty" bson:"_id,omitempty"`
	Source          string        `json:"source" bson:"source"`
	Views           uint          `json:"views,omitempty" bson:"views"`
	ExpireAt        time.Time     `json:"expire_at,omitempty" bson:"expire_at"`
	CreatedAt       time.Time     `json:"created_at,omitempty" bson:"created_at"`
	TTLMilliseconds time.Duration `json:"ttl_milliseconds,omitempty" bson:"ttl"`
	Slug            string        `json:"slug,omitempty" bson:"slug,omitempty"`
}

type CreateLinkDTO struct {
	Source string `json:"source"`
}
