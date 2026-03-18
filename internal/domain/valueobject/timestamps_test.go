package valueobject

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCreatedAt_JSONSerialization(t *testing.T) {
	t.Run("CreatedAt with valid time serializes to JSON correctly", func(t *testing.T) {
		// Truncate to seconds since RFC3339 loses nanosecond precision
		now := time.Now().UTC().Truncate(time.Second)
		createdAt := NewCreatedAt(now)

		data, err := json.Marshal(createdAt)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		var parsed CreatedAt
		err = json.Unmarshal(data, &parsed)
		if err != nil {
			t.Errorf("expected no error on unmarshal, got: %v", err)
		}

		if !parsed.Time().Equal(now) {
			t.Errorf("time mismatch after JSON round-trip: got %v, want %v", parsed.Time(), now)
		}
	})

	t.Run("CreatedAt handles zero time", func(t *testing.T) {
		var createdAt CreatedAt
		data, err := json.Marshal(createdAt)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		var parsed CreatedAt
		err = json.Unmarshal(data, &parsed)
		if err != nil {
			t.Errorf("expected no error on unmarshal, got: %v", err)
		}
	})
}

func TestUpdatedAt_JSONSerialization(t *testing.T) {
	t.Run("UpdatedAt with valid time serializes to JSON correctly", func(t *testing.T) {
		// Truncate to seconds since RFC3339 loses nanosecond precision
		now := time.Now().UTC().Truncate(time.Second)
		updatedAt := NewUpdatedAt(now)

		data, err := json.Marshal(updatedAt)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		var parsed UpdatedAt
		err = json.Unmarshal(data, &parsed)
		if err != nil {
			t.Errorf("expected no error on unmarshal, got: %v", err)
		}

		if !parsed.Time().Equal(now) {
			t.Errorf("time mismatch after JSON round-trip: got %v, want %v", parsed.Time(), now)
		}
	})
}

func TestPublishedAt_JSONSerialization(t *testing.T) {
	t.Run("PublishedAt with valid time serializes to JSON correctly", func(t *testing.T) {
		// Truncate to seconds since RFC3339 loses nanosecond precision
		now := time.Now().UTC().Truncate(time.Second)
		publishedAt := NewPublishedAt(&now)

		data, err := json.Marshal(publishedAt)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		var parsed PublishedAt
		err = json.Unmarshal(data, &parsed)
		if err != nil {
			t.Errorf("expected no error on unmarshal, got: %v", err)
		}

		if parsed.Time() == nil || !parsed.Time().Equal(now) {
			t.Errorf("time mismatch after JSON round-trip: got %v, want %v", parsed.Time(), now)
		}
	})

	t.Run("PublishedAt handles nil time", func(t *testing.T) {
		publishedAt := NewPublishedAt(nil)

		data, err := json.Marshal(publishedAt)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		var parsed PublishedAt
		err = json.Unmarshal(data, &parsed)
		if err != nil {
			t.Errorf("expected no error on unmarshal, got: %v", err)
		}

		if parsed.Time() != nil {
			t.Error("expected nil time after JSON round-trip")
		}
	})
}
