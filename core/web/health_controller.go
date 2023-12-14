package web

import (
	"cmp"
	"net/http"
	"slices"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

type HealthController struct {
	App chainlink.Application
}

const (
	HealthStatusPassing = "passing"
	HealthStatusFailing = "failing"
)

// NOTE: We only implement the k8s readiness check, *not* the liveness check. Liveness checks are only recommended in cases
// where the app doesn't crash itself on panic, and if implemented incorrectly can cause cascading failures.
// See the following for more information:
// - https://srcco.de/posts/kubernetes-liveness-probes-are-dangerous.html
func (hc *HealthController) Readyz(c *gin.Context) {
	status := http.StatusOK

	checker := hc.App.GetHealthChecker()

	ready, errors := checker.IsReady()

	if !ready {
		status = http.StatusServiceUnavailable
	}

	c.Status(status)

	if _, ok := c.GetQuery("full"); !ok {
		return
	}

	checks := make([]presenters.Check, 0, len(errors))

	for name, err := range errors {
		status := HealthStatusPassing
		var output string

		if err != nil {
			status = HealthStatusFailing
			output = err.Error()
		}

		checks = append(checks, presenters.Check{
			JAID:   presenters.NewJAID(name),
			Name:   name,
			Status: status,
			Output: output,
		})
	}

	// return a json description of all the checks
	jsonAPIResponse(c, checks, "checks")
}

func (hc *HealthController) Health(c *gin.Context) {
	status := http.StatusOK

	checker := hc.App.GetHealthChecker()

	healthy, errors := checker.IsHealthy()

	if !healthy {
		status = http.StatusServiceUnavailable
	}

	c.Status(status)

	checks := make([]presenters.Check, 0, len(errors))

	for name, err := range errors {
		status := HealthStatusPassing
		var output string

		if err != nil {
			status = HealthStatusFailing
			output = err.Error()
		}

		checks = append(checks, presenters.Check{
			JAID:   presenters.NewJAID(name),
			Name:   name,
			Status: status,
			Output: output,
		})
	}

	if testing.Testing() {
		slices.SortFunc(checks, func(a, b presenters.Check) int {
			return cmp.Compare(a.Name, b.Name)
		})
	}

	// return a json description of all the checks
	jsonAPIResponseWithStatus(c, checks, "checks", status)
}
