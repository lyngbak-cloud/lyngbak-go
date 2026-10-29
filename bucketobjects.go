package lyngbak

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// BucketObjectsService handles the Bucket Objects endpoints (OpenAPI tag
// "Bucket Objects").
type BucketObjectsService struct {
	client *Client
}

// List returns a page of objects in a bucket. opts may be nil. Use the
// returned NextCursor as opts.Cursor to fetch the next page while HasMore
// is true.
// Operation: listBucketObjects (GET /api/v1/buckets/{uuid}/objects/list)
func (s *BucketObjectsService) List(ctx context.Context, bucketUUID string, opts *BucketObjectListOptions) (*BucketObjectList, error) {
	query := url.Values{}
	if opts != nil {
		if opts.Prefix != "" {
			query.Set("prefix", opts.Prefix)
		}
		if opts.Limit > 0 {
			query.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Cursor != "" {
			query.Set("cursor", opts.Cursor)
		}
	}
	out := new(BucketObjectList)
	path := fmt.Sprintf("/api/v1/buckets/%s/objects/list", url.PathEscape(bucketUUID))
	if err := s.client.do(ctx, http.MethodGet, path, query, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// GenerateUploadURL generates a presigned upload URL for an object.
//
// The OpenAPI spec does not define the response schema for this endpoint, so
// the decoded JSON object is returned as a map.
// Operation: generateUploadUrl (POST /api/v1/buckets/{uuid}/objects/upload_url)
func (s *BucketObjectsService) GenerateUploadURL(ctx context.Context, bucketUUID string, req *GenerateUploadURLRequest) (map[string]any, error) {
	var out map[string]any
	path := fmt.Sprintf("/api/v1/buckets/%s/objects/upload_url", url.PathEscape(bucketUUID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GenerateDownloadURL generates a presigned download URL for an object.
//
// The OpenAPI spec does not define the response schema for this endpoint, so
// the decoded JSON object is returned as a map.
// Operation: generateDownloadUrl (POST /api/v1/buckets/{uuid}/objects/download_url)
func (s *BucketObjectsService) GenerateDownloadURL(ctx context.Context, bucketUUID, key string) (map[string]any, error) {
	body := struct {
		Key string `json:"key"`
	}{Key: key}
	var out map[string]any
	path := fmt.Sprintf("/api/v1/buckets/%s/objects/download_url", url.PathEscape(bucketUUID))
	if err := s.client.do(ctx, http.MethodPost, path, nil, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateFolder creates a folder in a bucket. folderPath should end with "/".
// Operation: createBucketFolder (POST /api/v1/buckets/{uuid}/objects/create_folder)
func (s *BucketObjectsService) CreateFolder(ctx context.Context, bucketUUID, folderPath string) error {
	body := struct {
		FolderPath string `json:"folder_path"`
	}{FolderPath: folderPath}
	path := fmt.Sprintf("/api/v1/buckets/%s/objects/create_folder", url.PathEscape(bucketUUID))
	return s.client.do(ctx, http.MethodPost, path, nil, body, nil)
}

// Delete deletes an object from a bucket by key.
// Operation: deleteBucketObject (DELETE /api/v1/buckets/{uuid}/objects/{key})
func (s *BucketObjectsService) Delete(ctx context.Context, bucketUUID, key string) error {
	path := fmt.Sprintf("/api/v1/buckets/%s/objects/%s", url.PathEscape(bucketUUID), url.PathEscape(key))
	return s.client.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// Rename renames (moves) an object within a bucket.
// Operation: renameBucketObject (PUT /api/v1/buckets/{uuid}/objects/{key}/rename)
func (s *BucketObjectsService) Rename(ctx context.Context, bucketUUID, key, newKey string) error {
	body := struct {
		NewKey string `json:"new_key"`
	}{NewKey: newKey}
	path := fmt.Sprintf("/api/v1/buckets/%s/objects/%s/rename", url.PathEscape(bucketUUID), url.PathEscape(key))
	return s.client.do(ctx, http.MethodPut, path, nil, body, nil)
}
