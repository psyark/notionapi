package notionapi

// Code generated by notion.codegen; DO NOT EDIT.

// https://developers.notion.com/reference/property-item-object

// A property_item object describes the identifier, type, and value of a page property. It's returned from the Retrieve a page property item
type PropertyItem struct {
	// Each page property item object contains the following keys. In addition, it will contain a key corresponding with the value of type. The value is an object containing type-specific data. The type-specific data are described in the sections below.
	Object         string                   `json:"object"`                           // Always "property_item".
	ID             string                   `json:"id"`                               // Underlying identifier for the property. This identifier is guaranteed to remain constant when the property name changes. It may be a UUID, but is often a short random string.  The id may be used in place of name when creating or updating pages.
	Type           string                   `json:"type"`                             // Type of the property. Possible values are "rich_text", "number", "select", "multi_select", "date", "formula", "relation", "rollup", "title", "people", "files", "checkbox", "url", "email", "phone_number", "created_time", "created_by", "last_edited_time", and "last_edited_by".
	Title          *RichText                `json:"title" specific:"type"`            // Title property value objects contain an array of rich text objects within the title property.
	RichText       *RichText                `json:"rich_text" specific:"type"`        // Rich Text property value objects contain an array of rich text objects within the rich_text property.
	Number         float64                  `json:"number" specific:"type"`           // Number property value objects contain a number within the number property.
	Select         *SelectOption            `json:"select" specific:"type"`           // Select property value objects contain the following data within the select property:
	Status         *StatusOption            `json:"status" specific:"type"`           // undocumented
	MultiSelect    []SelectOption           `json:"multi_select" specific:"type"`     // Multi-select property value objects contain an array of multi-select option values within the multi_select property.
	Date           *DateValue               `json:"date" specific:"type"`             // Date property value objects contain the following data within the date property:
	Formula        *FormulaPropertyItemData `json:"formula" specific:"type"`          // Formula property value objects represent the result of evaluating a formula described in the  database's properties. These objects contain a type key and a key corresponding with the value of type. The value is an object containing type-specific data. The type-specific data are described in the sections below.
	Relation       *PageReference           `json:"relation" specific:"type"`         // Relation property value objects contain an array of relation property items with page references within the relation property. A page reference is an object with an id property which is a string value (UUIDv4) corresponding to a page ID in another database.
	Rollup         *RollupPropertyItemData  `json:"rollup" specific:"type"`           // Rollup property value objects represent the result of evaluating a rollup described in the  database's properties. The property is returned as a list object of type property_item with a list of relation items used to computed the rollup under results.   A rollup property item is also returned under the property_type key that describes the rollup aggregation and computed result.   In order to avoid timeouts, if the rollup has a with a large number of aggregations or properties the endpoint returns a next_cursor value that is used to determinate the aggregation value so far for the subset of relations that have been paginated through.   Once has_more is false, then the final rollup value is returned.  See the Pagination documentation for more information on pagination in the Notion API.   Computing the values of following aggregations are not supported. Instead the endpoint returns a list of property_item objects for the rollup: * show_unique (Show unique values) * unique (Count unique values) * median(Median)
	People         *User                    `json:"people" specific:"type"`           // People property value objects contain an array of user objects within the people property.
	Files          []File                   `json:"files" specific:"type"`            // File property value objects contain an array of file references within the files property. A file reference is an object with a File Object and name property, with a string value corresponding to a filename of the original file upload (i.e. "Whole_Earth_Catalog.jpg").
	Checkbox       bool                     `json:"checkbox" specific:"type"`         // Checkbox property value objects contain a boolean within the checkbox property.
	URL            string                   `json:"url" specific:"type"`              // URL property value objects contain a non-empty string within the url property. The string describes a web address (i.e. "http://worrydream.com/EarlyHistoryOfSmalltalk/").
	Email          string                   `json:"email" specific:"type"`            // Email property value objects contain a string within the email property. The string describes an email address (i.e. "hello@example.org").
	PhoneNumber    string                   `json:"phone_number" specific:"type"`     // Phone number property value objects contain a string within the phone_number property. No structure is enforced.
	CreatedTime    string                   `json:"created_time" specific:"type"`     // Created time property value objects contain a string within the created_time property. The string contains the date and time when this page was created. It is formatted as an ISO 8601 date time string (i.e. "2020-03-17T19:10:04.968Z").
	CreatedBy      *User                    `json:"created_by" specific:"type"`       // Created by property value objects contain a user object within the created_by property. The user object describes the user who created this page.
	LastEditedTime string                   `json:"last_edited_time" specific:"type"` // Last edited time property value objects contain a string within the last_edited_time property. The string contains the date and time when this page was last updated. It is formatted as an ISO 8601 date time string (i.e. "2020-03-17T19:10:04.968Z").
	LastEditedBy   *User                    `json:"last_edited_by" specific:"type"`   // Last edited by property value objects contain a user object within the last_edited_by property. The user object describes the user who last updated this page.
}

func (c *PropertyItem) propertyItemOrPropertyItemPagination() {}
func (p PropertyItem) MarshalJSON() ([]byte, error) {
	type Alias PropertyItem
	return marshalByType(Alias(p), p.Type)
}

// The title, rich_text, relation and people property items of are returned as a paginated list object of individual property_item objects in the results. An abridged set of the the properties found in the list object are found below, see the Pagination documentation for additional information.
type PaginatedPropertyItem struct {
	ID       string                  `json:"id"`       // Underlying identifier for the property. This identifier is guaranteed to remain constant when the property name changes. It may be a UUID, but is often a short random string.  The id may be used in place of name when creating or updating pages.
	Type     string                  `json:"type"`     // Always "property_item".
	NextURL  *string                 `json:"next_url"` // The URL the user can request to get the next page of results.
	Title    struct{}                `json:"title" specific:"type"`
	RichText struct{}                `json:"rich_text" specific:"type"`
	Relation struct{}                `json:"relation" specific:"type"`
	Rollup   *RollupPropertyItemData `json:"rollup" specific:"type"`
	People   struct{}                `json:"people" specific:"type"`
}

func (p PaginatedPropertyItem) MarshalJSON() ([]byte, error) {
	type Alias PaginatedPropertyItem
	return marshalByType(Alias(p), p.Type)
}

/*
Formula property value objects represent the result of evaluating a formula described in the
database's properties. These objects contain a type key and a key corresponding with the value of type. The value is an object containing type-specific data. The type-specific data are described in the sections below.
*/
type FormulaPropertyItemData struct {
	Type    string     `json:"type"`                    // The type of the formula result. Possible values are "string", "number", "boolean", and "date".
	String  *string    `json:"string" specific:"type"`  // String formula property values contain an optional string within the string property.
	Number  *float64   `json:"number" specific:"type"`  // Number formula property values contain an optional number within the number property.
	Boolean bool       `json:"boolean" specific:"type"` // Boolean formula property values contain a boolean within the boolean property.
	Date    *DateValue `json:"date" specific:"type"`    // Date formula property values contain an optional date property value within the date property.
}

func (p FormulaPropertyItemData) MarshalJSON() ([]byte, error) {
	type Alias FormulaPropertyItemData
	return marshalByType(Alias(p), p.Type)
}

/*
Rollup property value objects represent the result of evaluating a rollup described in the
database's properties. The property is returned as a list object of type property_item with a list of relation items used to computed the rollup under results.

A rollup property item is also returned under the property_type key that describes the rollup aggregation and computed result.

In order to avoid timeouts, if the rollup has a with a large number of aggregations or properties the endpoint returns a next_cursor value that is used to determinate the aggregation value so far for the subset of relations that have been paginated through.

Once has_more is false, then the final rollup value is returned.  See the Pagination documentation for more information on pagination in the Notion API.

Computing the values of following aggregations are not supported. Instead the endpoint returns a list of property_item objects for the rollup:
* show_unique (Show unique values)
* unique (Count unique values)
* median(Median)
*/
type RollupPropertyItemData struct {
	Type       string     `json:"type"`                       // The type of rollup. Possible values are "number", "date", "array", "unsupported" and "incomplete".
	Function   string     `json:"function"`                   // Describes the aggregation used.  Possible values include: count_all, count_values, count_unique_values, count_empty, count_not_empty, percent_empty, percent_not_empty, sum, average, median, min, max, range, show_original
	Number     float64    `json:"number" specific:"type"`     // Number rollup property values contain a number within the number property.
	Date       *DateValue `json:"date" specific:"type"`       // Date rollup property values contain a date property value within the date property.
	Array      []struct{} `json:"array" specific:"type"`      // Array rollup property values contain an array of property_item objects within the results property.
	Incomplete *struct{}  `json:"incomplete" specific:"type"` // Rollups with an aggregation with more than one page of aggregated results will return a rollup object of type "incomplete". To obtain the final value paginate through the next values in the rollup using the next_cursor or next_url property.
}

func (p RollupPropertyItemData) MarshalJSON() ([]byte, error) {
	type Alias RollupPropertyItemData
	return marshalByType(Alias(p), p.Type)
}
