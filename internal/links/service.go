package links

import (
	"context"
	"fmt"
	"time"

	"github.com/furrygem/nocut-api/pkg/logging"
)

type Service struct {
	storage Storage
	logger  *logging.Logger
	linkTTL time.Duration
}

func (s *Service) Create(ctx context.Context, dto CreateLinkDTO) (l Link, err error) {
	l = Link{
		Source:    dto.Source,
		Views:     0,
		CreatedAt: time.Now(),
		ExpireAt:  time.Now().Add(s.linkTTL),
	}
	id, err := s.storage.Create(ctx, l)
	if err != nil {
		return l, fmt.Errorf("error creating link: %v", err)
	}
	l, err = s.storage.FindOne(ctx, id)
	return l, err
}

func (s *Service) GetLinkById(ctx context.Context, id string) (l Link, err error) {
	l, err = s.storage.FindOne(ctx, id)
	if err != nil {
		return l, fmt.Errorf("error getting link '%s': %v", id, err)
	}
	l.Slug, err = IdToUrl(id)
	if err != nil {
		return l, fmt.Errorf("error decoding id to slug '%s'. %v", id, err)
	}
	return l, err
}
