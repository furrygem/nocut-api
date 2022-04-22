package links

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/furrygem/nocut-api/pkg/logging"
)

type Service struct {
	storage Storage
	logger  *logging.Logger
	linkTTL time.Duration
}

func (s *Service) SendURLCheckResults(ucs *URLCheckException, rw http.ResponseWriter) error {
	resp, _ := json.Marshal(ucs)
	_, err := rw.Write(resp)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) runURLChecks(sourceURL string) *URLCheckException {
	ucs := URLCheckException{
		HostIsUp:   false,
		URLIsValid: false,
	}

	urlisvalid, err := URLIsValid(sourceURL)
	ucs.URLIsValid = urlisvalid
	if err != nil {
		s.logger.Warning("URL Host Is not valid. Assuming the host is not available.")
		ucs.HostIsUp = false
	}

	if ucs.URLIsValid {
		HostIsUp, err := URLHostIsUp(sourceURL)
		ucs.HostIsUp = HostIsUp
		if err != nil {
			s.logger.Warningf("URL Host is Down at %s. %v", sourceURL, err)
		}
	}
	return &ucs
}

func (s *Service) Create(ctx context.Context, dto CreateLinkDTO) (l Link, err error) {
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

func (s *Service) GetLinkByID(ctx context.Context, id string) (l Link, err error) {
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
