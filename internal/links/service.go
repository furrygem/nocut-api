package links

import (
	"context"
	"fmt"
	"time"

	"github.com/furrygem/nocut-api/pkg/logging"
)

// Service object
type Service struct {
	storage         Storage
	blackListedURls []string
	logger          *logging.Logger
	linkTTL         time.Duration
}

// CheckLinkService for CheckLinkHandler
func (s *Service) CheckLinkService(url string) *URLCheckException {
	uce := s.runURLChecks(url)
	return uce
}

func (s *Service) runURLChecks(sourceURL string) *URLCheckException {
	uc := URLCheckException{
		HostIsUp:   false,
		URLIsValid: false,
	}

	if StringInSlice(sourceURL, s.blackListedURls) {
		return &uc
	}

	urlisvalid, err := URLIsValid(sourceURL, s.blackListedURls)
	uc.URLIsValid = urlisvalid
	if err != nil {
		s.logger.Warning("URL Host Is not valid. Assuming the host is not available.")
		uc.HostIsUp = false
	}

	if uc.URLIsValid {
		HostIsUp, err := URLHostIsUp(sourceURL)
		uc.HostIsUp = HostIsUp
		if err != nil {
			s.logger.Warningf("URL Host is Down at %s. %v", sourceURL, err)
		}
	}
	return &uc
}

// CreateService for CreateHandler
func (s *Service) CreateService(ctx context.Context, dto CreateLinkDTO) (l Link, err error) {
	l = Link{
		Source:    dto.Source,
		Views:     0,
		CreatedAt: time.Now(),
		ExpireAt:  time.Now().Add(s.linkTTL),
	}
	summary := s.runURLChecks(dto.Source)
	if !(summary.HostIsUp && summary.URLIsValid) {
		return l, summary
	}
	id, isDup, err := s.storage.Create(ctx, l)
	if isDup {
		s.logger.Warnf("Duplicate link %s. %v", dto.Source, err)
		l, err = s.storage.FindOneBySource(ctx, dto.Source)
		l.Slug, err = IDToURL(l.ID)
		return l, err
	}
	if err != nil {
		s.logger.Errorf("error creating link: %v", err)
		s.logger.Errorf("%s", err.Error())
		return l, fmt.Errorf("error creating link: %v", err)
	}
	l, err = s.storage.FindOne(ctx, id)
	l.Slug, err = IDToURL(id)
	return l, err
}

// GetLinkByIDService for GetLinkByIDHandler
func (s *Service) GetLinkByIDService(ctx context.Context, id string) (l Link, err error) {
	l, err = s.storage.FindOne(ctx, id)
	if err != nil {
		return l, fmt.Errorf("error getting link '%s': %v", id, err)
	}
	l.Slug, err = IDToURL(id)
	if err != nil {
		return l, fmt.Errorf("error decoding id to slug '%s'. %v", id, err)
	}
	return l, err
}
