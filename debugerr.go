package notionapi

import (
	"fmt"

	"github.com/yudai/gojsondiff"
)

type debugErr struct {
	responseBody []byte
	remarshaled  []byte
	diff         gojsondiff.Diff
}

func (e debugErr) Error() string {
	// res, _ := formatter.NewDeltaFormatter().Format(e.diff)
	return fmt.Sprintf("validation failed. \nwant: %v\ngot: %v", string(e.responseBody), string(e.remarshaled))
}
