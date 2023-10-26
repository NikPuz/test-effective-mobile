package getaway

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"test-zero-agency/internal/app/config"
	"test-zero-agency/internal/entity"
)

var ErrBadNameEnrichment = entity.NewLogicError(nil, "inappropriate name for enrichment", 400)

type getAgeResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   *int   `json:"age"`
}

type getGenderResponse struct {
	Count  int     `json:"count"`
	Name   string  `json:"name"`
	Gender *string `json:"gender"`
}

type getNationalityResponse struct {
	Count   int    `json:"count"`
	Name    string `json:"name"`
	Country []struct {
		CountryId   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type enrichment struct {
	enrichmentAgeDomain         string
	enrichmentGenderDomain      string
	enrichmentNationalityDomain string
	client                      *http.Client
}

func NewEnrichment(cfg *config.Config) entity.IPeopleGetaway {
	return &enrichment{
		enrichmentAgeDomain:         cfg.EnrichmentAgeDomain,
		enrichmentGenderDomain:      cfg.EnrichmentGenderDomain,
		enrichmentNationalityDomain: cfg.EnrichmentNationalityDomain,
		client:                      &http.Client{Timeout: cfg.EnrichmentTimeout},
	}
}

func (e *enrichment) GetAgeByName(ctx context.Context, name string) (int, error) {
	var ageResponse getAgeResponse

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://%s/?name=%s", e.enrichmentAgeDomain, name), nil)
	if err != nil {
		return 0, err
	}

	resp, err := e.client.Do(req.WithContext(ctx))
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		var errorResponse errorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return 0, errors.New(fmt.Sprintf("age request status code not OK: %d", resp.StatusCode))
		}
		return 0, errors.New(fmt.Sprintf("age request error: %s", errorResponse.Error))
	}

	err = json.NewDecoder(resp.Body).Decode(&ageResponse)
	if err != nil {
		return 0, err
	}
	if ageResponse.Age == nil {
		return 0, ErrBadNameEnrichment
	}

	return *ageResponse.Age, nil
}

func (e *enrichment) GetGenderByName(ctx context.Context, name string) (string, error) {
	var genderResponse getGenderResponse

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://%s/?name=%s", e.enrichmentGenderDomain, name), nil)
	if err != nil {
		return "", err
	}

	resp, err := e.client.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		var errorResponse errorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return "", errors.New(fmt.Sprintf("gender request status code not OK: %d", resp.StatusCode))
		}
		return "", errors.New(fmt.Sprintf("gender request error: %s", errorResponse.Error))
	}

	err = json.NewDecoder(resp.Body).Decode(&genderResponse)
	if err != nil {
		return "", err
	}
	if genderResponse.Gender == nil {
		return "", ErrBadNameEnrichment
	}

	return *genderResponse.Gender, nil
}

func (e *enrichment) GetNationalityByName(ctx context.Context, name string) (string, error) {
	var nationalityResponse getNationalityResponse

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://%s/?name=%s", e.enrichmentNationalityDomain, name), nil)
	if err != nil {
		return "", err
	}

	resp, err := e.client.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		var errorResponse errorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return "", errors.New(fmt.Sprintf("nationality request status code not OK: %d", resp.StatusCode))
		}
		return "", errors.New(fmt.Sprintf("nationality request error: %s", errorResponse.Error))
	}

	err = json.NewDecoder(resp.Body).Decode(&nationalityResponse)
	if err != nil {
		return "", err
	}
	if len(nationalityResponse.Country) == 0 {
		return "", ErrBadNameEnrichment
	}

	return nationalityResponse.Country[0].CountryId, nil
}
