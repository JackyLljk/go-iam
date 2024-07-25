package authorization

import (
	"context"

	"github.com/ory/ladon"

	"j-iam/pkg/log"
)

// AuditLogger outputs and cache information about granting or rejecting policies.
type AuditLogger struct {
	client AuthorizationInterface
}

// NewAuditLogger creates a AuditLogger with default parameters.
func NewAuditLogger(client AuthorizationInterface) *AuditLogger {
	return &AuditLogger{
		client: client,
	}
}

// LogRejectedAccessRequest write rejected subject access to log.
func (a *AuditLogger) LogRejectedAccessRequest(ctx context.Context, r *ladon.Request, p ladon.Policies, d ladon.Policies) {
	a.client.LogRejectedAccessRequest(r, p, d)
	log.Debug("subject access review rejected", log.Any("request", r), log.Any("deciders", d))
}

// LogGrantedAccessRequest write granted subject access to log.
func (a *AuditLogger) LogGrantedAccessRequest(ctx context.Context, r *ladon.Request, p ladon.Policies, d ladon.Policies) {
	a.client.LogGrantedAccessRequest(r, p, d)
	log.Debug("subject access review granted", log.Any("request", r), log.Any("deciders", d))
}
