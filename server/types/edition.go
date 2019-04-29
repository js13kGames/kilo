package types

type Edition struct {
	Id         uint16              `json:"id"`
	Title      string              `json:"title"`
	Timeframe  Timeframe           `json:"timeframe"`
	Categories EditionCategoryList `json:"-"`
	Judges     []*Person           `json:"-"`
}

type EditionCategory struct {
	Title       string        `json:"title"`
	Submissions []*Submission `json:"-"`
	bits        EditionCategoryMask
}

func (cat *EditionCategory) Mask() EditionCategoryMask { return cat.bits }

type EditionCategoryList []*EditionCategory
