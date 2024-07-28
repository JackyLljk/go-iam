package store

import (
	pb "j-iam/internal/pkg/proto/apiserver/v1"
)

// SecretStore defines the secret storage interface.
type SecretStore interface {
	List() (map[string]*pb.SecretInfo, error)
}

// List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.SecretList, error)
