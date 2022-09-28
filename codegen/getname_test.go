package codegen

import (
	"testing"
)

func TestGetName(t *testing.T) {

	hoge := map[string]string{
		"any":                                  "Any",
		"avatar_url":                           "AvatarURL",
		"block_id":                             "BlockID",
		"Bookmark block Data":                  "BookmarkBlockData",
		"boolean":                              "Boolean",
		"bot":                                  "Bot",
		"Bulleted list item block Data":        "BulletedListItemBlockData",
		"bulleted_list_item":                   "BulletedListItem",
		"divider":                              "Divider",
		"does_not_equal":                       "DoesNotEqual",
		"Dual_property_relation_configuration": "DualPropertyRelationConfiguration",
		"Multi-select filter condition":        "MultiSelectFilterCondition",
		"none":                                 "None",
		"number":                               "Number",
		"option_ids":                           "OptionIds",
		"options":                              "Options",
		"or":                                   "Or",
		"pdf":                                  "PDF",
		"PDF block Data":                       "PDFBlockData",
		"to_do":                                "ToDo",
		"url":                                  "URL",
	}
	for input, want := range hoge {
		got := getName(input)
		if want != got {
			t.Errorf("%q: %q,\n", want, got)
		}
	}
}
