package models

type KeySplitError string

func (e KeySplitError) Error() string {
	return "error splitting key"
}

const (
	PkSplitError KeySplitError = "partition-key-split-error"
	SkSplitError KeySplitError = "search-key-split-error"
)
