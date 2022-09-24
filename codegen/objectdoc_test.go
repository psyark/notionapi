package codegen

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestStripMarkdown(t *testing.T) {
	cases := map[string]string{
		"_An external file is any URL that isn't hosted by Notion_.": "An external file is any URL that isn't hosted by Notion.",
		"last_edited_time":                             "last_edited_time",
		"database_id database_id":                      "database_id database_id",
		"[property values](ref:property-value-object)": "property values",
		"**bolded** *italicized*":                      "bolded italicized",
		"*he\n*who\n*me":                               "*he\n*who\n*me",
	}

	for input, expected := range cases {
		if output := stripMarkdown(input); output != expected {
			t.Errorf("mismatch: \nwant: %v\ngot : %v", expected, output)
		}
	}
}

func ExampleParseObjectDoc() {
	sections, err := ParseObjectDoc("https://developers.notion.com/reference/page")
	if err != nil {
		panic(err)
	}

	for _, section := range sections {
		if section.Heading != nil {
			fmt.Printf("%v\n", section.Heading.Text)
		}

		for _, element := range section.Elements {
			switch element := element.(type) {
			case *BlockCodeElement, *BlockCalloutElement, *BlockParametersElement:
				data, _ := json.MarshalIndent(element, "    ", "  ")
				fmt.Printf("    %v\n", string(data))
			case *ParagraphElement:
				content := strings.ReplaceAll(element.Content, "\n", " ")
				if content != "" {
					fmt.Printf("    %v\n", content)
				}
			default:
				panic(reflect.TypeOf(element))
			}
		}
	}

	// Output:
	// The Page object contains the property values of a single Notion page.  All pages have a Parent. If the parent is a database, the property values conform to the schema laid out database's properties. Otherwise, the only property value is the title.  Page content is available as blocks. The content can be read using retrieve block children and appended using append block children.
	//     {
	//       "type": "info",
	//       "title": "",
	//       "body": "Properties marked with an * are available to integrations with any capabilities. Other properties require read content capabilities in order to be returned from the Notion API. For more information on integration capabilities, see the [capabilities guide](ref:capabilities)."
	//     }
	// All pages
	//     [
	//       {
	//         "Property": "object*",
	//         "Type": "string",
	//         "Description": "Always \"page\".",
	//         "Example value": "\"page\""
	//       },
	//       {
	//         "Property": "id*",
	//         "Type": "string (UUIDv4)",
	//         "Description": "Unique identifier of the page.",
	//         "Example value": "\"45ee8d13-687b-47ce-a5ca-6e2e45548c4b\""
	//       },
	//       {
	//         "Property": "created_time",
	//         "Type": "string (ISO 8601 date and time)",
	//         "Description": "Date and time when this page was created. Formatted as an ISO 8601 date time string.",
	//         "Example value": "\"2020-03-17T19:10:04.968Z\""
	//       },
	//       {
	//         "Property": "created_by",
	//         "Type": "Partial User",
	//         "Description": "User who created the page.",
	//         "Example value": "{\"object\": \"user\",\"id\": \"45ee8d13-687b-47ce-a5ca-6e2e45548c4b\"}"
	//       },
	//       {
	//         "Property": "last_edited_time",
	//         "Type": "string (ISO 8601 date and time)",
	//         "Description": "Date and time when this page was updated. Formatted as an ISO 8601 date time string.",
	//         "Example value": "\"2020-03-17T19:10:04.968Z\""
	//       },
	//       {
	//         "Property": "last_edited_by",
	//         "Type": "Partial User",
	//         "Description": "User who last edited the page.",
	//         "Example value": "{\"object\": \"user\",\"id\": \"45ee8d13-687b-47ce-a5ca-6e2e45548c4b\"}"
	//       },
	//       {
	//         "Property": "archived",
	//         "Type": "boolean",
	//         "Description": "The archived status of the page.",
	//         "Example value": "false"
	//       },
	//       {
	//         "Property": "icon",
	//         "Type": "File Object (only type of \"external\" is supported currently) or Emoji object",
	//         "Description": "Page icon.",
	//         "Example value": ""
	//       },
	//       {
	//         "Property": "cover",
	//         "Type": "File object (only type of \"external\" is supported currently)",
	//         "Description": "Page cover image.",
	//         "Example value": ""
	//       },
	//       {
	//         "Property": "properties",
	//         "Type": "object",
	//         "Description": "Property values of this page. As of version 2022-06-28, properties only contains the ID of the property; in prior versions properties contained the values as well.\n\nIf parent.type is \"page_id\" or \"workspace\", then the only valid key is title.\n\nIf parent.type is \"database_id\", then the keys and values of this field are determined by the properties  of the database this page belongs to.\n\nkey string\nName of a property as it appears in Notion.\n\nvalue object\nSee Property value object.",
	//         "Example value": "{ \"id\": \"A%40Hk\" }"
	//       },
	//       {
	//         "Property": "parent",
	//         "Type": "object",
	//         "Description": "Information about the page's parent. See Parent object.",
	//         "Example value": "{ \"type\": \"database_id\", \"database_id\": \"d9824bdc-8445-4327-be8b-5b47500af6ce\" }"
	//       },
	//       {
	//         "Property": "url",
	//         "Type": "string",
	//         "Description": "The URL of the Notion page.",
	//         "Example value": "\"https://www.notion.so/Avocado-d093f1d200464ce78b36e58a3f0d8043\""
	//       }
	//     ]
}
