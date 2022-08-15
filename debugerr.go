package notionapi

import (
	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

type debugErr struct {
	responseBody []byte
	remarshaled  []byte
	diff         gojsondiff.Diff
}

func (e debugErr) Error() string {
	res, _ := formatter.NewDeltaFormatter().Format(e.diff)
	return res
}
