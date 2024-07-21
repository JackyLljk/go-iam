package store

import (
	"context"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	v1 "j-iam/pkg/model/v1"
)

type UserStore interface {
	Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error
}
