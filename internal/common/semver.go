package common

import (
	"encoding/json"
	"github.com/Masterminds/semver/v3"
)

type SemverConstraints struct {
	original string
	*semver.Constraints
}

func NewSemverConstraints(c string) (*SemverConstraints, error) {
	constraints, err := semver.NewConstraint(c)
	if err != nil {
		return nil, err
	}

	return &SemverConstraints{
		original:    c,
		Constraints: constraints,
	}, nil
}

func (c *SemverConstraints) String() string {
	return c.original
}

// UnmarshalJSON implements JSON.Unmarshaler interface.
func (c *SemverConstraints) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	temp, err := semver.NewConstraint(s)
	if err != nil {
		return err
	}
	c.original = s
	c.Constraints = temp
	return nil
}

// MarshalJSON implements JSON.Marshaler interface.
func (c *SemverConstraints) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}
