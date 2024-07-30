package apiserver

import (
	"context"
	"encoding/json"
	v1 "j-iam/internal/pkg/proto/apiserver/v1"

	"github.com/AlekSi/pointer"
	"github.com/avast/retry-go"
	"github.com/marmotedu/errors"
	"github.com/ory/ladon"

	"j-iam/pkg/log"
)

type policies struct {
	cli v1.CacheClient
}

func newPolicies(ds *datastore) *policies {
	return &policies{ds.client}
}

// List returns all the authorization policies.
func (p *policies) List() (map[string][]*ladon.DefaultPolicy, error) {
	pols := make(map[string][]*ladon.DefaultPolicy)

	log.Info("Loading policies")

	req := &v1.ListPoliciesRequest{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(-1),
	}

	var resp *v1.ListPoliciesResponse
	err := retry.Do(
		func() error {
			var listErr error
			resp, listErr = p.cli.ListPolicies(context.Background(), req)
			if listErr != nil {
				return listErr
			}

			return nil
		}, retry.Attempts(3),
	)
	if err != nil {
		return nil, errors.Wrap(err, "list policies failed")
	}

	log.Infof("Policies found (%d total)[username:name]:", len(resp.Items))

	for _, v := range resp.Items {
		log.Infof(" - %s:%s", v.Username, v.Name)

		var policy ladon.DefaultPolicy

		if err := json.Unmarshal([]byte(v.PolicyShadow), &policy); err != nil {
			log.Warnf("failed to load policy for %s, error: %s", v.Name, err.Error())

			continue
		}

		pols[v.Username] = append(pols[v.Username], &policy)
	}

	return pols, nil
}
