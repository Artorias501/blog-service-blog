package valueobject

import (
	"encoding/json"
	"time"
)

// CreatedAt represents the creation timestamp of an entity
type CreatedAt struct {
	value time.Time
}

// NewCreatedAt creates a new CreatedAt with the given time
func NewCreatedAt(value time.Time) CreatedAt {
	return CreatedAt{value: value}
}

// NowCreatedAt creates a new CreatedAt with current time
func NowCreatedAt() CreatedAt {
	return CreatedAt{value: time.Now().UTC()}
}

// Time returns the underlying time.Time value
func (c CreatedAt) Time() time.Time {
	return c.value
}

// MarshalJSON implements json.Marshaler
func (c CreatedAt) MarshalJSON() ([]byte, error) {
	if c.value.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(c.value.Format(time.RFC3339))
}

// UnmarshalJSON implements json.Unmarshaler
func (c *CreatedAt) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == nil {
		c.value = time.Time{}
		return nil
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return err
	}
	c.value = t
	return nil
}

// UpdatedAt represents the last update timestamp of an entity
type UpdatedAt struct {
	value time.Time
}

// NewUpdatedAt creates a new UpdatedAt with the given time
func NewUpdatedAt(value time.Time) UpdatedAt {
	return UpdatedAt{value: value}
}

// NowUpdatedAt creates a new UpdatedAt with current time
func NowUpdatedAt() UpdatedAt {
	return UpdatedAt{value: time.Now().UTC()}
}

// Time returns the underlying time.Time value
func (u UpdatedAt) Time() time.Time {
	return u.value
}

// MarshalJSON implements json.Marshaler
func (u UpdatedAt) MarshalJSON() ([]byte, error) {
	if u.value.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(u.value.Format(time.RFC3339))
}

// UnmarshalJSON implements json.Unmarshaler
func (u *UpdatedAt) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == nil {
		u.value = time.Time{}
		return nil
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return err
	}
	u.value = t
	return nil
}

// After reports whether the time instant u is after other
func (u UpdatedAt) After(other UpdatedAt) bool {
	return u.value.After(other.value)
}

// PublishedAt represents the publication timestamp of a post
type PublishedAt struct {
	value *time.Time
}

// NewPublishedAt creates a new PublishedAt with the given time pointer
func NewPublishedAt(value *time.Time) PublishedAt {
	return PublishedAt{value: value}
}

// Time returns the underlying time.Time pointer
func (p PublishedAt) Time() *time.Time {
	return p.value
}

// MarshalJSON implements json.Marshaler
func (p PublishedAt) MarshalJSON() ([]byte, error) {
	if p.value == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(p.value.Format(time.RFC3339))
}

// UnmarshalJSON implements json.Unmarshaler
func (p *PublishedAt) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == nil {
		p.value = nil
		return nil
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return err
	}
	p.value = &t
	return nil
}
