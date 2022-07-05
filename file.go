package shopify

import (
	"context"
	"fmt"

	"github.com/r0busta/go-shopify-graphql-model/graph/model"
)

// FileCreateInput file creation params
type FileCreateInput struct {
	OriginalSource string `json:"originalSource"`
	ContentType    string `json:"contentType"`
}

// FileError errors wrapper
type FileError struct {
	Code    string `json:"code"`
	Details string `json:"details"`
	Message string `json:"message"`
}

// File uploaded file
type File struct {
	Alt        string      `json:"alt"`
	CreatedAt  string      `json:"createdAt"`
	FileStatus string      `json:"fileStatus"`
	FileErrors []FileError `json:"fileErrors"`
}

type fileCreateResult struct {
	Files      []File            `json:"files"`
	UserErrors []model.UserError `json:"userErrors,omitempty"`
}

type mutationFileCreate struct {
	FileCreateResult fileCreateResult `graphql:"fileCreate(files: $files)" json:"fileCreate"`
}

// FileCreateService supports creation of staged upload targets
type FileCreateService interface {
	FileCreate(ctx context.Context, files []FileCreateInput) ([]File, error)
}

// FileCreateMutationClient defines the required mutation client operations.
type FileCreateMutationClient interface {
	Mutate(ctx context.Context, m interface{}, variables map[string]interface{}) error
}

type FileCreateOp struct {
	mutationClient FileCreateMutationClient
}

func (f *FileCreateOp) FileCreate(ctx context.Context, files []FileCreateInput) ([]File, error) {
	m := mutationFileCreate{}

	err := f.mutationClient.Mutate(ctx, &m, map[string]interface{}{
		"files": files,
	})
	if err != nil {
		return nil, err
	}

	if len(m.FileCreateResult.UserErrors) > 0 {
		return nil, fmt.Errorf("%+v", m.FileCreateResult.UserErrors)
	}

	return m.FileCreateResult.Files, nil
}
