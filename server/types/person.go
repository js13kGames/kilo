package types

import (
	"github.com/js13kgames/kilo/server/types/eid"
)

type ActorId = eid.ID

const (
	ActorPersonEidType = 1
)

type Person struct {
	Id          ActorId          `json:"id"`
	Email       string           `json:"-"`
	Name        PersonName       `json:"name"`
	Avatar      string           `json:"avatar,omitempty"`
	Bio         string           `json:"bio,omitempty"`
	Website     string           `json:"website,omitempty"`
	Identities  []PersonIdentity `json:"identities"`
	Submissions []SubmissionId   `json:"submissions"`
}

type PersonName struct {
	Given   string `json:"given,omitempty"`
	Family  string `json:"family,omitempty"`
	Display string `json:"display,omitempty"`
	Sort    string `json:"sort,omitempty"`
}

type PersonIdentity interface {
	Lookup(*ProviderMaps) ActorId
	Bind(*ProviderMaps, ActorId)
	Unbind(*ProviderMaps)
	Equals(PersonIdentity) bool
}
