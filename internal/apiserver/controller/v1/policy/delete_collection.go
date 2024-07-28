package policy

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"

	"j-iam/internal/pkg/middleware"
	"j-iam/pkg/log"
)

// DeleteCollection delete policies by policy names.
func (p *PolicyController) DeleteCollection(c *gin.Context) {
	log.L(c).Info("batch delete policy function called.")

	if err := p.srv.Policies().DeleteCollection(c, c.GetString(middleware.UsernameKey),
		c.QueryArray("name"), metav1.DeleteOptions{}); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
