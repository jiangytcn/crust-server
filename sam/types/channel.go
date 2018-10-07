package types

import (
	"encoding/json"
	"fmt"
	"time"
)

type (
	Channel struct {
		ID    uint64          `json:"id" db:"id"`
		Name  string          `json:"name" db:"name"`
		Topic string          `json:"topic" db:"topic"`
		Type  ChannelType     `json:"type" db:"type"`
		Meta  json.RawMessage `json:"-" db:"meta"`

		CreatorID      uint64 `json:"creatorId" db:"rel_creator"`
		OrganisationID uint64 `json:"organisationId" db:"rel_organisation"`

		CreatedAt  time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt  *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
		ArchivedAt *time.Time `json:"archivedAt,omitempty" db:"archived_at"`
		DeletedAt  *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`

		LastMessageID uint64 `json:",omitempty" db:"rel_last_message"`

		Member  *ChannelMember `json:"-" db:"-"`
		Members []uint64       `json:"-" db:"-"`
	}

	ChannelFilter struct {
		Query string

		// Only return channels accessible by this user
		CurrentUserID uint64

		IncludeMembers bool
	}

	ChannelType string

	ChannelSet []*Channel
)

// Scope returns permissions group that for this type
func (r *Channel) Scope() string {
	return "channel"
}

// Resource returns a RBAC resource ID for this type
func (r *Channel) Resource() string {
	return fmt.Sprintf("%s:%d", r.Scope(), r.ID)
}

// Operation returns a RBAC resource-scoped role name for an operation
func (r *Channel) Operation(name string) string {
	return fmt.Sprintf("%s/%s", r.Resource(), name)
}

func (cc ChannelSet) Walk(w func(*Channel) error) (err error) {
	for i := range cc {
		if err = w(cc[i]); err != nil {
			return
		}
	}

	return
}

const (
	ChannelTypePublic  ChannelType = "public"
	ChannelTypePrivate             = "private"
	ChannelTypeGroup               = "group"
)
