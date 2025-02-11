package v1

import (
	"context"

	v1 "j-iam/internal/pkg/model/apiserver/v1"

	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/errors"

	"j-iam/internal/apiserver/store"
	"j-iam/internal/pkg/code"
)

// SecretService defines functions used to handle secret request.
type SecretService interface {
	Create(ctx context.Context, secret *v1.Secret, opts metav1.CreateOptions) error
	Update(ctx context.Context, secret *v1.Secret, opts metav1.UpdateOptions) error
	Delete(ctx context.Context, username, secretID string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, username string, secretIDs []string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, username, secretID string, opts metav1.GetOptions) (*v1.Secret, error)
	List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.SecretList, error)
}

type secretService struct {
	store store.Factory
}

var _ SecretService = (*secretService)(nil)

func newSecrets(srv *service) *secretService {
	return &secretService{store: srv.store}
}

func (s *secretService) Create(ctx context.Context, secret *v1.Secret, opts metav1.CreateOptions) error {
	if err := s.store.Secrets().Create(ctx, secret, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (s *secretService) Update(ctx context.Context, secret *v1.Secret, opts metav1.UpdateOptions) error {
	// Save changed fields.
	if err := s.store.Secrets().Update(ctx, secret, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (s *secretService) Delete(ctx context.Context, username, secretID string, opts metav1.DeleteOptions) error {
	if err := s.store.Secrets().Delete(ctx, username, secretID, opts); err != nil {
		return err
	}

	return nil
}

func (s *secretService) DeleteCollection(
	ctx context.Context,
	username string,
	secretIDs []string,
	opts metav1.DeleteOptions,
) error {
	if err := s.store.Secrets().DeleteCollection(ctx, username, secretIDs, opts); err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}

	return nil
}

func (s *secretService) Get(
	ctx context.Context,
	username, secretID string,
	opts metav1.GetOptions,
) (*v1.Secret, error) {
	secret, err := s.store.Secrets().Get(ctx, username, secretID, opts)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (s *secretService) List(ctx context.Context, username string, opts metav1.ListOptions) (*v1.SecretList, error) {
	secrets, err := s.store.Secrets().List(ctx, username, opts)
	if err != nil {
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}

	return secrets, nil
}
