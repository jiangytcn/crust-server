package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type (
	User struct {
		ID       uint64   `json:"userID,string" db:"id"`
		Username string   `json:"username" db:"username"`
		Email    string   `json:"email" db:"email"`
		Name     string   `json:"name" db:"name"`
		Handle   string   `json:"handle" db:"handle"`
		Kind     UserKind `json:"kind" db:"kind"`
		SatosaID string   `json:"-" db:"satosa_id"`

		Meta *UserMeta `json:"meta" db:"meta"`

		OrganisationID uint64 `json:"organisationID,string" db:"rel_organisation"`
		RelatedUserID  uint64 `json:"relatedUserID,string" db:"rel_user_id"`
		User           *User  `json:"user" db:"-"`

		Password []byte `json:"-" db:"password"`

		CreatedAt   time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt   *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
		SuspendedAt *time.Time `json:"suspendedAt,omitempty" db:"suspended_at"`
		DeletedAt   *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`

		Roles []*Role `json:"roles,omitempty" db:"-"`
	}

	UserMeta struct {
		Avatar string `json:"avatar,omitempty"`
	}

	UserFilter struct {
		Query    string
		Email    string
		Username string
		OrderBy  string
	}

	UserKind string
)

const (
	NormalUser UserKind = ""
	BotUser             = "bot"
)

func (u *User) Valid() bool {
	return u.ID > 0 && u.SuspendedAt == nil && u.DeletedAt == nil
}

func (u *User) Identity() uint64 {
	return u.ID
}

func (u *User) ValidatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(password)) == nil
}

func (u *User) GeneratePassword(password string) error {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = pwd
	return nil
}

func (mm *UserMeta) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return errors.Errorf("User.Meta must be a string, got %T instead", value)
	}
	return json.Unmarshal([]byte(str), mm)
}

func (mm *UserMeta) Value() (driver.Value, error) {
	return json.Marshal(mm)
}
