package types

type SubmissionId uint64

type Submission struct {
	Id          SubmissionId        `json:"id"`
	Edition     uint16              `json:"edition"`
	Categories  EditionCategoryMask `json:"categories"`
	Title       string              `json:"title"`
	Slug        string              `json:"slug"`
	Description string              `json:"description"`
	Repository  string              `json:"repository"`
	Owner       ActorId             `json:"owner"`
}

// EditionCategoryMask holds the bitmask of the categories the Submission is attached to.
// This being a single byte, at most 8 categories can exist per edition which gets
// enforced by the respective endpoints.
type EditionCategoryMask byte

func (c *EditionCategoryMask) Add(m EditionCategoryMask)      { *c |= m }
func (c *EditionCategoryMask) Remove(m EditionCategoryMask)   { *c &^= m }
func (c *EditionCategoryMask) Has(m EditionCategoryMask) bool { return *c&m != 0 }
