package pgx_pq_example

import (
	"os"
	"testing"
)

func requirePostgresIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv("RUN_POSTGRES_INTEGRATION") != "1" {
		t.Skip("set RUN_POSTGRES_INTEGRATION=1 to run tests against PostgreSQL")
	}
	for _, name := range []string{"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB"} {
		if os.Getenv(name) == "" {
			t.Skipf("set %s to run PostgreSQL integration tests", name)
		}
	}
}
