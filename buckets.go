package lyngbak

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// BucketsService handles the Buckets endpoints (OpenAPI tag "Buckets").
type BucketsService struct {
	client *Client
}

// List returns all buckets in the workspace.
// Operation: listBuckets (GET /api/v1/buckets)
func (s *BucketsService) List(ctx context.Context) ([]Bucket, error) {
	var out []Bucket
	if err := s.client.do(ctx, http.MethodGet, "/api/v1/buckets", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Create creates a new bucket.
// Operation: createBucket (POST /api/v1/buckets)
func (s *BucketsService) Create(ctx context.Context, req *CreateBucketRequest) (*Bucket, error) {
	out := new(Bucket)
	if err := s.client.do(ctx, http.MethodPost, "/api/v1/buckets", nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get returns a single bucket by UUID.
// Operation: getBucket (GET /api/v1/buckets/{uuid})
func (s *BucketsService) Get(ctx context.Context, uuid string) (*Bucket, error) {
	out := new(Bucket)
	path := fmt.Sprintf("/api/v1/buckets/%s", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodGet, path, nil, nil, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Update updates a bucket.
// Operation: updateBucket (PUT /api/v1/buckets/{uuid})
func (s *BucketsService) Update(ctx context.Context, uuid string, req *UpdateBucketRequest) (*Bucket, error) {
	out := new(Bucket)
	path := fmt.Sprintf("/api/v1/buckets/%s", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodPut, path, nil, req, out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete deletes a bucket.
// Operation: deleteBucket (DELETE /api/v1/buckets/{uuid})
func (s *BucketsService) Delete(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("/api/v1/buckets/%s", url.PathEscape(uuid))
	return s.client.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// UpdatePolicy sets the bucket's S3-compatible policy. policy is a JSON
// string.
// Operation: updateBucketPolicy (PUT /api/v1/buckets/{uuid}/policy)
func (s *BucketsService) UpdatePolicy(ctx context.Context, uuid, policy string) (*Bucket, error) {
	body := struct {
		Policy string `json:"policy"`
	}{Policy: policy}
	out := new(Bucket)
	path := fmt.Sprintf("/api/v1/buckets/%s/policy", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodPut, path, nil, body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeletePolicy removes the bucket's S3-compatible policy.
// Operation: deleteBucketPolicy (DELETE /api/v1/buckets/{uuid}/policy)
func (s *BucketsService) DeletePolicy(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("/api/v1/buckets/%s/policy", url.PathEscape(uuid))
	return s.client.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// UpdateACL sets the bucket's S3-compatible ACL. acl is an XML string.
// Operation: updateBucketAcl (PUT /api/v1/buckets/{uuid}/acl)
func (s *BucketsService) UpdateACL(ctx context.Context, uuid, acl string) (*Bucket, error) {
	body := struct {
		ACL string `json:"acl"`
	}{ACL: acl}
	out := new(Bucket)
	path := fmt.Sprintf("/api/v1/buckets/%s/acl", url.PathEscape(uuid))
	if err := s.client.do(ctx, http.MethodPut, path, nil, body, out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteACL removes the bucket's S3-compatible ACL.
// Operation: deleteBucketAcl (DELETE /api/v1/buckets/{uuid}/acl)
func (s *BucketsService) DeleteACL(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("/api/v1/buckets/%s/acl", url.PathEscape(uuid))
	return s.client.do(ctx, http.MethodDelete, path, nil, nil, nil)
}
