package notionapi

// Code generated by notion.codegen; DO NOT EDIT.

// https://developers.notion.com/reference/property-value-object

// A property value defines the identifier, type, and value of a page property in a page object. It's used when retrieving and updating pages, ex: Create and Update pages.
type PropertyValue struct {
	// Each page property value object contains the following keys. In addition, it contains a key corresponding with the value of type. The value is an object containing type-specific data. The type-specific data are described in the sections below.
	ID             string                   `json:"id,omitempty"`                     // Underlying identifier for the property. This identifier is guaranteed to remain constant when the property name changes. It may be a UUID, but is often a short random string.  The id may be used in place of name when creating or updating pages.
	Type           string                   `json:"type,omitempty"`                   // Type of the property. Possible values are "rich_text", "number", "select", "multi_select", "status", "date", "formula", "relation", "rollup", "title", "people", "files", "checkbox", "url", "email", "phone_number", "created_time", "created_by", "last_edited_time", and "last_edited_by".
	Title          []RichText               `json:"title" specific:"type"`            // Title property value objects contain an array of rich text objects within the title property.
	RichText       []RichText               `json:"rich_text" specific:"type"`        // Rich Text property value objects contain an array of rich text objects within the rich_text property.
	Number         *float64                 `json:"number" specific:"type"`           // Number property value objects contain a number within the number property.
	Select         *SelectOption            `json:"select" specific:"type"`           // Select property value objects contain the following data within the select property:
	Status         *StatusOption            `json:"status" specific:"type"`           // Status property value objects contain the following data within the status property:
	MultiSelect    []SelectOption           `json:"multi_select" specific:"type"`     // Multi-select property value objects contain an array of multi-select option values within the multi_select property.
	Date           *DateValue               `json:"date" specific:"type"`             // Date property value objects contain the following data within the date property:
	Formula        FormulaPropertyValueData `json:"formula" specific:"type"`          // Formula property value objects represent the result of evaluating a formula described in the  database's properties. These objects contain a type key and a key corresponding with the value of type. The value of a formula cannot be updated directly.
	Relation       []PageReference          `json:"relation" specific:"type"`         // Relation property value objects contain an array of page references within the relation property. A page reference is an object with an id property, with a string value (UUIDv4) corresponding to a page ID in another database. Updating with an empty array will clear the list. Note: a relation property has a maximum of 100 pages.
	Rollup         RollupPropertyValueData  `json:"rollup" specific:"type"`           // Rollup property value objects represent the result of evaluating a rollup described in the  database's properties. These objects contain a type key and a key corresponding with the value of type. The value of a rollup cannot be updated directly.
	People         []User                   `json:"people" specific:"type"`           // People property value objects contain an array of user objects within the people property.
	Files          []File                   `json:"files" specific:"type"`            // File property value objects contain an array of file references within the files property. A file reference is an object with a File Object and name property, with a string value corresponding to a filename of the original file upload (i.e. "Whole_Earth_Catalog.jpg").
	Checkbox       *bool                    `json:"checkbox" specific:"type"`         // Checkbox property value objects contain a boolean within the checkbox property.
	URL            *string                  `json:"url" specific:"type"`              // URL property value objects contain a non-empty string within the url property. The string describes a web address (i.e. "http://worrydream.com/EarlyHistoryOfSmalltalk/").
	Email          *string                  `json:"email" specific:"type"`            // Email property value objects contain a string within the email property. The string describes an email address (i.e. "hello@example.org").
	PhoneNumber    *string                  `json:"phone_number" specific:"type"`     // Phone number property value objects contain a string within the phone_number property. No structure is enforced.
	CreatedTime    *string                  `json:"created_time" specific:"type"`     // Created time property value objects contain a string within the created_time property. The string contains the date and time when this page was created. It is formatted as an ISO 8601 date time string (i.e. "2020-03-17T19:10:04.968Z"). The value of created_time cannot be updated. See the Property Item Object to see how these values are returned.
	CreatedBy      *User                    `json:"created_by" specific:"type"`       // Created by property value objects contain a user object within the created_by property. The user object describes the user who created this page. The value of created_by cannot be updated. See the Property Item Object to see how these values are returned.
	LastEditedTime *string                  `json:"last_edited_time" specific:"type"` // Last edited time property value objects contain a string within the last_edited_time property. The string contains the date and time when this page was last updated. It is formatted as an ISO 8601 date time string (i.e. "2020-03-17T19:10:04.968Z"). The value of last_edited_time cannot be updated. See the Property Item Object to see how these values are returned.
	LastEditedBy   *User                    `json:"last_edited_by" specific:"type"`   // Last edited by property value objects contain a user object within the last_edited_by property. The user object describes the user who last updated this page. The value of last_edited_by cannot be updated. See the Property Item Object to see how these values are returned.
}

func (p PropertyValue) MarshalJSON() ([]byte, error) {
	type Alias PropertyValue
	return marshalByType(Alias(p), p.Type)
}

// Date property value objects contain the following data within the date property:
type DateValue struct {
	Start    ISO8601String  `json:"start"`     // An ISO 8601 format date, with optional time.
	End      *ISO8601String `json:"end"`       // An ISO 8601 formatted date, with optional time. Represents the end of a date range.  If null, this property's date value is not a range.
	TimeZone *string        `json:"time_zone"` // Time zone information for start and end. Possible values are extracted from the IANA database and they are based on the time zones from Moment.js.  When time zone is provided, start and end should not have any UTC offset. In addition, when time zone  is provided, start and end cannot be dates without time information.  If null, time zone information will be contained in UTC offsets in start and end.
}

/*
Formula property value objects represent the result of evaluating a formula described in the
database's properties. These objects contain a type key and a key corresponding with the value of type. The value of a formula cannot be updated directly.
*/
type FormulaPropertyValueData struct {
	Type    string     `json:"type"`
	String  *string    `json:"string" specific:"type"`  // String formula property values contain an optional string within the string property.
	Number  *float64   `json:"number" specific:"type"`  // Number formula property values contain an optional number within the number property.
	Boolean bool       `json:"boolean" specific:"type"` // Boolean formula property values contain a boolean within the boolean property.
	Date    *DateValue `json:"date" specific:"type"`    // Date formula property values contain an optional date property value within the date property.
}

func (p FormulaPropertyValueData) MarshalJSON() ([]byte, error) {
	type Alias FormulaPropertyValueData
	return marshalByType(Alias(p), p.Type)
}

/*
Rollup property value objects represent the result of evaluating a rollup described in the
database's properties. These objects contain a type key and a key corresponding with the value of type. The value of a rollup cannot be updated directly.
*/
type RollupPropertyValueData struct {
	Type     string          `json:"type"`                   // These objects contain a type key and a key corresponding with the value of type.
	Function string          `json:"function"`               // undocumented
	String   *string         `json:"string" specific:"type"` // String rollup property values contain an optional string within the string property.
	Number   float64         `json:"number" specific:"type"` // Number rollup property values contain a number within the number property.
	Date     *DateValue      `json:"date" specific:"type"`   // Date rollup property values contain a date property value within the date property.
	Array    []PropertyValue `json:"array" specific:"type"`  // Array rollup property values contain an array of number, date, or string objects within the results property.
}

func (p RollupPropertyValueData) MarshalJSON() ([]byte, error) {
	type Alias RollupPropertyValueData
	return marshalByType(Alias(p), p.Type)
}
