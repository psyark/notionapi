package mapping

import (
	"reflect"
	"strings"
)

type tagInfo struct {
	name  string
	icon  bool
	cover bool
}

func parseTag(field reflect.StructField) *tagInfo {
	if tagStr, ok := field.Tag.Lookup("notion"); ok {
		return &tagInfo{
			name:  strings.SplitN(tagStr, ",", 2)[0],
			icon:  strings.Contains(tagStr, ",icon"),
			cover: strings.Contains(tagStr, ",cover"),
		}
	}
	return nil
}
