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

// FileService file operations
type FileService interface {
	FileCreate(ctx context.Context, files []FileCreateInput) ([]File, error)
	GetFileStatusByName(ctx context.Context, fileName string) (*FileStatus, error)
}

// GQLClient defines the required mutation client operations.
type GQLClient interface {
	Mutate(ctx context.Context, m interface{}, variables map[string]interface{}) error
	QueryString(ctx context.Context, q string, variables map[string]interface{}, v interface{}) error
}

// FileOp wraps file GQL API
type FileOp struct {
	gqlClient GQLClient
}

// FileCreate encapsulates file creation API, returns limited set of parameters
func (f *FileOp) FileCreate(ctx context.Context, files []FileCreateInput) ([]File, error) {
	m := mutationFileCreate{}

	err := f.gqlClient.Mutate(ctx, &m, map[string]interface{}{
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

// FileStatus represents file uploads processing status
type FileStatus struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Status string `json:"fileStatus"`
}

type fileQueryResult struct {
	Files struct {
		Edges []struct {
			Node FileStatus `json:"node"`
		} `json:"edges"`
	} `json:"files"`
}

const fileByNameQuery = `
	{
	 	files (first: 1, query: "filename:%s") {
			edges {
				node {
					fileStatus
					... on GenericFile {
						id url
					}
				}
			}
		}
	}
`

// GetFileStatusByName returns file upload processing status, assumes the file name is UNIQUE!
// For non-unique file name the first occurrence will be returned
func (f *FileOp) GetFileStatusByName(ctx context.Context, fileName string) (*FileStatus, error) {
	var result fileQueryResult

	err := f.gqlClient.QueryString(ctx, fmt.Sprintf(fileByNameQuery, fileName), nil, &result)
	if err != nil {
		return nil, err
	}

	if len(result.Files.Edges) > 0 {
		return &result.Files.Edges[0].Node, nil
	}

	return nil, nil
}
