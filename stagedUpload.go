package shopify

import (
	"context"
	"fmt"

	"github.com/r0busta/go-shopify-graphql-model/graph/model"
)

// StagedUploadInput represents staged upload request
type StagedUploadInput struct {
	FileName string `json:"filename"`
	MimeType string `json:"mimeType"`
	Resource string `json:"resource"`
}

// StagedMediaUploadTarget represents result of a staged upload
type StagedMediaUploadTarget struct {
	URL         string       `json:"url"`
	ResourceURL string       `json:"resourceUrl"`
	Parameters  []*Parameter `json:"parameters"`
}

// Parameter represents staged upload parameter
type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type stagedUploadCreateResult struct {
	StagedTargets []StagedMediaUploadTarget `json:"stagedTargets"`
	UserErrors    []model.UserError         `json:"userErrors,omitempty"`
}

type mutationStagedUploadsCreate struct {
	StagedUploadsCreateResult stagedUploadCreateResult `graphql:"stagedUploadsCreate(input: $input)" json:"stagedUploadsCreate"`
}

// StagedUploadService supports creation of staged upload targets
type StagedUploadService interface {
	CreateStagedUpload(ctx context.Context, uploads []StagedUploadInput) ([]StagedMediaUploadTarget, error)
}

// StagedUploadMutationClient defines the required mutation client operations.
type StagedUploadMutationClient interface {
	Mutate(ctx context.Context, m interface{}, variables map[string]interface{}) error
}

type StagedUploadOp struct {
	mutationClient StagedUploadMutationClient
}

func (s *StagedUploadOp) CreateStagedUpload(
	ctx context.Context, uploads []StagedUploadInput) ([]StagedMediaUploadTarget, error) {
	m := mutationStagedUploadsCreate{}

	err := s.mutationClient.Mutate(ctx, &m, map[string]interface{}{
		"input": uploads,
	})

	if err != nil {
		return nil, err
	}

	if len(m.StagedUploadsCreateResult.UserErrors) > 0 {
		return nil, fmt.Errorf("%+v", m.StagedUploadsCreateResult.UserErrors)
	}

	return m.StagedUploadsCreateResult.StagedTargets, nil
}
