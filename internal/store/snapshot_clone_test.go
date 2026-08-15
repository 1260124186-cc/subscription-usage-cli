package store

import (
	"testing"

	"github.com/1260124186-cc/subscription-usage-cli/internal/domain"
)

func TestSnapshotCloneDoesNotShareUsageEvents(t *testing.T) {
	original := Snapshot{Usage: []domain.UsageEvent{{ID: "evt-1", AccountID: "acme", Units: 4}}}
	clone := original.Clone()
	clone.Usage[0].Units = 99

	if got, want := original.Usage[0].Units, int64(4); got != want {
		t.Fatalf("original usage units = %d, want %d", got, want)
	}
}
