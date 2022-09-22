package codegen

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func ExampleParseObjectDoc() {
	sections, err := ParseObjectDoc("https://developers.notion.com/reference/page")
	if err != nil {
		panic(err)
	}

	for _, section := range sections {
		if section.Heading != nil {
			fmt.Printf("%v\n", section.Heading.text)
		}

		for _, element := range section.Elements {
			switch element := element.(type) {
			case *BlockCodeElement, *BlockCalloutElement:
				data, _ := json.MarshalIndent(element, "    ", "  ")
				fmt.Printf("    %v\n", string(data))
			case *BlockParametersElement:
				data, _ := json.MarshalIndent(element.MapSlice(), "    ", "  ")
				fmt.Printf("    %v\n", string(data))
			case *ParagraphElement:
				content := strings.ReplaceAll(element.content, "\n", " ")
				if content != "" {
					fmt.Printf("    %v\n", content)
				}
			default:
				panic(reflect.TypeOf(element))
			}
		}
	}

	// Output:
	// The Page object contains the [property values](ref:property-value-object) of a single Notion page.  All pages have a [Parent](ref:parent-object). If the parent is a [database](ref:database), the property values conform to the schema laid out database's [properties](ref:property-object). Otherwise, the only property value is the `title`.  Page content is available as [blocks](ref:block). The content can be read using [retrieve block children](ref:get-block-children) and appended using [append block children](ref:patch-block-children).
	//     {
	//       "type": "info",
	//       "title": "",
	//       "body": "Properties marked with an * are available to integrations with any capabilities. Other properties require read content capabilities in order to be returned from the Notion API. For more information on integration capabilities, see the [capabilities guide](ref:capabilities)."
	//     }
	// ## All pages
	//     [
	//       {
	//         "Description": "Always `\"page\"`.",
	//         "Example value": "`\"page\"`",
	//         "Property": "`object`*",
	//         "Type": "`string`"
	//       },
	//       {
	//         "Description": "Unique identifier of the page.",
	//         "Example value": "`\"45ee8d13-687b-47ce-a5ca-6e2e45548c4b\"`",
	//         "Property": "`id`*",
	//         "Type": "`string` (UUIDv4)"
	//       },
	//       {
	//         "Description": "Date and time when this page was created. Formatted as an [ISO 8601 date time](https://en.wikipedia.org/wiki/ISO_8601) string.",
	//         "Example value": "`\"2020-03-17T19:10:04.968Z\"`",
	//         "Property": "`created_time`",
	//         "Type": "`string` ([ISO 8601 date and time](https://en.wikipedia.org/wiki/ISO_8601))"
	//       },
	//       {
	//         "Description": "User who created the page.",
	//         "Example value": "`{\"object\": \"user\",\"id\": \"45ee8d13-687b-47ce-a5ca-6e2e45548c4b\"}`",
	//         "Property": "`created_by`",
	//         "Type": "[Partial User](ref:user)"
	//       },
	//       {
	//         "Description": "Date and time when this page was updated. Formatted as an [ISO 8601 date time](https://en.wikipedia.org/wiki/ISO_8601) string.",
	//         "Example value": "`\"2020-03-17T19:10:04.968Z\"`",
	//         "Property": "`last_edited_time`",
	//         "Type": "`string` ([ISO 8601 date and time](https://en.wikipedia.org/wiki/ISO_8601))"
	//       },
	//       {
	//         "Description": "User who last edited the page.",
	//         "Example value": "`{\"object\": \"user\",\"id\": \"45ee8d13-687b-47ce-a5ca-6e2e45548c4b\"}`",
	//         "Property": "`last_edited_by`",
	//         "Type": "[Partial User](ref:user)"
	//       },
	//       {
	//         "Description": "The archived status of the page.",
	//         "Example value": "`false`",
	//         "Property": "`archived`",
	//         "Type": "`boolean`"
	//       },
	//       {
	//         "Description": "Page icon.",
	//         "Example value": "",
	//         "Property": "`icon`",
	//         "Type": "[File Object](ref:file-object) (only `type` of `\"external\"` is supported currently) or [Emoji object](ref:emoji-object)"
	//       },
	//       {
	//         "Description": "Page cover image.",
	//         "Example value": "",
	//         "Property": "`cover`",
	//         "Type": "[File object](ref:file-object) (only `type` of `\"external\"` is supported currently)"
	//       },
	//       {
	//         "Description": "Property values of this page. As of version `2022-06-28`, `properties` only contains the ID of the property; in prior versions `properties` contained the values as well.\n\nIf `parent.type` is `\"page_id\"` or `\"workspace\"`, then the only valid key is `title`.\n\nIf `parent.type` is `\"database_id\"`, then the keys and values of this field are determined by the [`properties`](https://developers.notion.com/reference/property-object)  of the [database](ref:database) this page belongs to.\n\n`key` string\nName of a property as it appears in Notion.\n\n`value` object\nSee [Property value object](https://developers.notion.com/reference/property-value-object).",
	//         "Example value": "`{ \"id\": \"A%40Hk\" }`",
	//         "Property": "`properties`",
	//         "Type": "`object`"
	//       },
	//       {
	//         "Description": "Information about the page's parent. See [Parent object](ref:parent-object).",
	//         "Example value": "`{ \"type\": \"database_id\", \"database_id\": \"d9824bdc-8445-4327-be8b-5b47500af6ce\" }`",
	//         "Property": "`parent`",
	//         "Type": "`object`"
	//       },
	//       {
	//         "Description": "The URL of the Notion page.",
	//         "Example value": "`\"https://www.notion.so/Avocado-d093f1d200464ce78b36e58a3f0d8043\"`",
	//         "Property": "`url`",
	//         "Type": "`string`"
	//       }
	//     ]
}
