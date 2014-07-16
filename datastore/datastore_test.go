package datastore

import "time"

func normalizeTime(t *time.Time) {
	*t = t.In(time.UTC)
}
