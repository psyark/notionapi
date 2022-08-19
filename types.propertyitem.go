package notionapi

import "fmt"

// Code generated by notion.codegen; DO NOT EDIT.

// https://developers.notion.com/reference/property-item-object

type PropertyItem interface {
	getCommon() *PropertyItemCommon
}

// Each page property item object contains the following keys. In addition, it will contain a key corresponding with the value of type. The value is an object containing type-specific data. The type-specific data are described in the sections below.
type PropertyItemCommon struct {
	Object  string  `json:"object"` // Always "property_item".
	ID      string  `json:"id"`     // Underlying identifier for the property. This identifier is guaranteed to remain constant when the property name changes. It may be a UUID, but is often a short random string.The id may be used in place of name when creating or updating pages.
	Type    string  `json:"type"`   // Type of the property. Possible values are "rich_text", "number", "select", "multi_select", "date", "formula", "relation", "rollup", "title", "people", "files", "checkbox", "url", "email", "phone_number", "created_time", "created_by", "last_edited_time", and "last_edited_by".
	NextURL *string `json:"-"`      // Only present in paginated property values (see below) with another page of results. If present, the url the user can request to get the next page of results.
}

func (i *PropertyItemCommon) getCommon() *PropertyItemCommon {
	return i
}

// Title property value objects contain an array of rich text objects within the title property.
type TitlePropertyItem struct {
	PropertyItemCommon
	Title RichText `json:"title"` // Title property value objects contain an array of rich text objects within the title property.
}

// Rich Text property value objects contain an array of rich text objects within the rich_text property.
type RichTextPropertyItem struct {
	PropertyItemCommon
	RichText RichText `json:"rich_text"` // Rich Text property value objects contain an array of rich text objects within the rich_text property.
}

// Number property value objects contain a number within the number property.
type NumberPropertyItem struct {
	PropertyItemCommon
	Number float64 `json:"number"` // Number property value objects contain a number within the number property.
}

// Select property value objects contain the following data within the select property:
type SelectPropertyItem struct {
	PropertyItemCommon
	Select SelectPropertyItemData `json:"select"`
}

type SelectPropertyItemData struct {
	ID    UUIDString `json:"id"`    // ID of the option.When updating a select property, you can use either name or id.
	Name  string     `json:"name"`  // Name of the option as it appears in Notion.If the select database property does not yet have an option by that name, it will be added to the database schema if the integration also has write access to the parent database.Note: Commas (",") are not valid for select values.
	Color string     `json:"color"` // Color of the option. Possible values are: "default", "gray", "brown", "red", "orange", "yellow", "green", "blue", "purple", "pink". Defaults to "default".Not currently editable.
}

// Multi-select property value objects contain an array of multi-select option values within the multi_select property.
type MultiSelectPropertyItem struct {
	PropertyItemCommon
	MultiSelect []MultiSelectOptionValues `json:"multi_select"` // Multi-select property value objects contain an array of multi-select option values within the multi_select property.
}

type MultiSelectOptionValues struct {
	ID    UUIDString `json:"id"`    // ID of the option.When updating a multi-select property, you can use either name or id.
	Name  string     `json:"name"`  // Name of the option as it appears in Notion.If the multi-select database property does not yet have an option by that name, it will be added to the database schema if the integration also has write access to the parent database.Note: Commas (",") are not valid for select values.
	Color string     `json:"color"` // Color of the option. Possible values are: "default", "gray", "brown", "red", "orange", "yellow", "green", "blue", "purple", "pink". Defaults to "default".Not currently editable.
}

// Date property value objects contain the following data within the date property:
type DatePropertyItem struct {
	PropertyItemCommon
	Date DatePropertyItemData `json:"date"`
}

type DatePropertyItemData struct {
	Start    ISO8601String  `json:"start"`     // An ISO 8601 format date, with optional time.
	End      *ISO8601String `json:"end"`       // An ISO 8601 formatted date, with optional time. Represents the end of a date range.If null, this property's date value is not a range.
	TimeZone *string        `json:"time_zone"` // Time zone information for start and end. Possible values are extracted from the IANA database and they are based on the time zones from Moment.js.When time zone is provided, start and end should not have any UTC offset. In addition, when time zone  is provided, start and end cannot be dates without time information.If null, time zone information will be contained in UTC offsets in start and end.
}

/*
Formula property value objects represent the result of evaluating a formula described in the
database's properties. These objects contain a type key and a key corresponding with the value of type. The value is an object containing type-specific data. The type-specific data are described in the sections below.
*/
type FormulaPropertyItem struct {
	PropertyItemCommon
}

// String formula property values contain an optional string within the string property.
type StringFormulaPropertyItem struct {
	PropertyItemCommon
	String string `json:"string"` // String formula property values contain an optional string within the string property.
}

// Number formula property values contain an optional number within the number property.
type NumberFormulaPropertyItem struct {
	PropertyItemCommon
	Number float64 `json:"number"` // Number formula property values contain an optional number within the number property.
}

// Boolean formula property values contain a boolean within the boolean property.
type BooleanFormulaPropertyItem struct {
	PropertyItemCommon
	Boolean *bool `json:"boolean"` // Boolean formula property values contain a boolean within the boolean property.
}

// Date formula property values contain an optional date property value within the date property.
type DateFormulaPropertyItem struct {
	PropertyItemCommon
	Date DatePropertyItem `json:"date"` // Date formula property values contain an optional date property value within the date property.
}

// Relation property value objects contain an array of relation property items with a pagereferences within the relation property. A page reference is an object with an id property, with a string value (UUIDv4) corresponding to a page ID in another database
type RelationPropertyItem struct {
	PropertyItemCommon
	Relation PageReference `json:"relation"` // Relation property value objects contain an array of relation property items with a pagereferences within the relation property. A page reference is an object with an id property, with a string value (UUIDv4) corresponding to a page ID in another database
}

/*
Rollup property value objects represent the result of evaluating a rollup described in the
database's properties. The property is returned as a list object of type property_item with a list of relation items used to computed the rollup under results.
A rollup property item is also returned under the property_type key that describes the rollup aggregation and computed result.
In order to avoid timeouts, if the rollup has a with a large number of aggregations or properties the endpoint returns a next_cursor value that is used to determinate the aggregation value so far for the subset of relations that have been paginated through.
Once has_more is false, then the final rollup value is returned.  See the Pagination documentation for more information on pagination in the Notion API.
Computing the values of following aggregations are not supported. Instead the endpoint returns a list of property_item objects for the rollup:

show_unique (Show unique values)
unique (Count unique values)
median(Median)
*/
type RollupPropertyItem struct {
	PropertyItemCommon
}

// Number rollup property values contain a number within the number property.
type NumberRollupPropertyItem struct {
	PropertyItemCommon
	Number float64 `json:"number"` // Number rollup property values contain a number within the number property.
}

// Date rollup property values contain a date property value within the date property.
type DateRollupPropertyItem struct {
	PropertyItemCommon
	Date DatePropertyItem `json:"date"` // Date rollup property values contain a date property value within the date property.
}

// Array rollup property values contain an array of property_item objects within the results property.
type ArrayRollupPropertyItem struct {
	PropertyItemCommon
	Results []DatePropertyItem `json:"results"` // Array rollup property values contain an array of property_item objects within the results property.
}

// Rollups with an aggregation with more than one page of aggregated results will return a rollup object of type "incomplete". To obtain the final value paginate through the next values in the rollup using the next_cursor or next_url property.
type IncompleteRollupPropertyItem struct {
	PropertyItemCommon
}

// People property value objects contain an array of user objects within the people property.
type PeoplePropertyItem struct {
	PropertyItemCommon
	People User `json:"people"` // People property value objects contain an array of user objects within the people property.
}

// File property value objects contain an array of file references within the files property. A file reference is an object with a File Object and name property, with a string value corresponding to a filename of the original file upload (i.e. "Whole_Earth_Catalog.jpg").
type FilesPropertyItem struct {
	PropertyItemCommon
	Files []File `json:"files"` // File property value objects contain an array of file references within the files property. A file reference is an object with a File Object and name property, with a string value corresponding to a filename of the original file upload (i.e. "Whole_Earth_Catalog.jpg").
}

// Checkbox property value objects contain a boolean within the checkbox property.
type CheckboxPropertyItem struct {
	PropertyItemCommon
	Checkbox *bool `json:"checkbox"` // Checkbox property value objects contain a boolean within the checkbox property.
}

// URL property value objects contain a non-empty string within the url property. The string describes a web address (i.e. "http://worrydream.com/EarlyHistoryOfSmalltalk/").
type URLPropertyItem struct {
	PropertyItemCommon
	URL string `json:"url"` // URL property value objects contain a non-empty string within the url property. The string describes a web address (i.e. "http://worrydream.com/EarlyHistoryOfSmalltalk/").
}

// Email property value objects contain a string within the email property. The string describes an email address (i.e. "[email protected]").
type EmailPropertyItem struct {
	PropertyItemCommon
	Email string `json:"email"` // Email property value objects contain a string within the email property. The string describes an email address (i.e. "[email protected]").
}

// Phone number property value objects contain a string within the phone_number property. No structure is enforced.
type PhoneNumberPropertyItem struct {
	PropertyItemCommon
	PhoneNumber string `json:"phone_number"` // Phone number property value objects contain a string within the phone_number property. No structure is enforced.
}

// Created time property value objects contain a string within the created_time property. The string contains the date and time when this page was created. It is formatted as an ISO 8601 date time string (i.e. "2020-03-17T19:10:04.968Z").
type CreatedTimePropertyItem struct {
	PropertyItemCommon
	CreatedTime string `json:"created_time"` // Created time property value objects contain a string within the created_time property. The string contains the date and time when this page was created. It is formatted as an ISO 8601 date time string (i.e. "2020-03-17T19:10:04.968Z").
}

// Created by property value objects contain a user object within the created_by property. The user object describes the user who created this page.
type CreatedByPropertyItem struct {
	PropertyItemCommon
	CreatedBy User `json:"created_by"` // Created by property value objects contain a user object within the created_by property. The user object describes the user who created this page.
}

// Last edited time property value objects contain a string within the last_edited_time property. The string contains the date and time when this page was last updated. It is formatted as an ISO 8601 date time string (i.e. "2020-03-17T19:10:04.968Z").
type LastEditedTimePropertyItem struct {
	PropertyItemCommon
	LastEditedTime string `json:"last_edited_time"` // Last edited time property value objects contain a string within the last_edited_time property. The string contains the date and time when this page was last updated. It is formatted as an ISO 8601 date time string (i.e. "2020-03-17T19:10:04.968Z").
}

// Last edited by property value objects contain a user object within the last_edited_by property. The user object describes the user who last updated this page.
type LastEditedByPropertyItem struct {
	PropertyItemCommon
	LastEditedBy User `json:"last_edited_by"` // Last edited by property value objects contain a user object within the last_edited_by property. The user object describes the user who last updated this page.
}

func createPropertyItem(typeName string) (PropertyItem, error) {
	switch typeName {
	case "title":
		return &TitlePropertyItem{}, nil
	case "rich_text":
		return &RichTextPropertyItem{}, nil
	case "number":
		return &NumberPropertyItem{}, nil
	case "select":
		return &SelectPropertyItem{}, nil
	case "multi_select":
		return &MultiSelectPropertyItem{}, nil
	case "date":
		return &DatePropertyItem{}, nil
	case "formula":
		return &FormulaPropertyItem{}, nil
	case "string_formula":
		return &StringFormulaPropertyItem{}, nil
	case "number_formula":
		return &NumberFormulaPropertyItem{}, nil
	case "boolean_formula":
		return &BooleanFormulaPropertyItem{}, nil
	case "date_formula":
		return &DateFormulaPropertyItem{}, nil
	case "relation":
		return &RelationPropertyItem{}, nil
	case "rollup":
		return &RollupPropertyItem{}, nil
	case "number_rollup":
		return &NumberRollupPropertyItem{}, nil
	case "date_rollup":
		return &DateRollupPropertyItem{}, nil
	case "array_rollup":
		return &ArrayRollupPropertyItem{}, nil
	case "incomplete_rollup":
		return &IncompleteRollupPropertyItem{}, nil
	case "people":
		return &PeoplePropertyItem{}, nil
	case "files":
		return &FilesPropertyItem{}, nil
	case "checkbox":
		return &CheckboxPropertyItem{}, nil
	case "url":
		return &URLPropertyItem{}, nil
	case "email":
		return &EmailPropertyItem{}, nil
	case "phone_number":
		return &PhoneNumberPropertyItem{}, nil
	case "created_time":
		return &CreatedTimePropertyItem{}, nil
	case "created_by":
		return &CreatedByPropertyItem{}, nil
	case "last_edited_time":
		return &LastEditedTimePropertyItem{}, nil
	case "last_edited_by":
		return &LastEditedByPropertyItem{}, nil
	default:
		return nil, fmt.Errorf("unsupported type: %v", typeName)
	}
}
