package notionapi

// Code generated by notion.codegen; DO NOT EDIT.

// https://developers.notion.com/reference/post-database-query-filter

type Filter interface {
	filter()
}

// Each database property filter object must contain a property key and a key corresponding with the type of the database property identified by property. The value is an object containing a type-specific filter condition.
type PropertyFilter struct {
	Property       string                      `json:"property"` // The name or ID of the property to filter on.
	Title          *TextFilterCondition        `json:"title,omitempty"`
	RichText       *TextFilterCondition        `json:"rich_text,omitempty"`
	URL            *TextFilterCondition        `json:"url,omitempty"`
	Email          *TextFilterCondition        `json:"email,omitempty"`
	PhoneNumber    *TextFilterCondition        `json:"phone_number,omitempty"`
	Number         *NumberFilterCondition      `json:"number,omitempty"`
	Checkbox       *CheckboxFilterCondition    `json:"checkbox,omitempty"`
	Select         *SelectFilterCondition      `json:"select,omitempty"`
	MultiSelect    *MultiSelectFilterCondition `json:"multi_select,omitempty"`
	Status         *StatusFilterCondition      `json:"status,omitempty"`
	Date           *DateFilterCondition        `json:"date,omitempty"`
	CreatedTime    *DateFilterCondition        `json:"created_time,omitempty"`
	LastEditedTime *DateFilterCondition        `json:"last_edited_time,omitempty"`
	People         *PeopleFilterCondition      `json:"people,omitempty"`
	CreatedBy      *PeopleFilterCondition      `json:"created_by,omitempty"`
	LastEditedBy   *PeopleFilterCondition      `json:"last_edited_by,omitempty"`
	Files          *FilesFilterCondition       `json:"files,omitempty"`
	Relation       *RelationFilterCondition    `json:"relation,omitempty"`
	Rollup         *RollupFilterCondition      `json:"rollup,omitempty"`
	Formula        *FormulaFilterCondition     `json:"formula,omitempty"`
}

func (f PropertyFilter) filter() {}

var _ Filter = PropertyFilter{}

// A timestamp filter object must contain a timestamp key corresponding to the type of timestamp and a key matching that timestamp type which contains a date filter condition.
type TimestampFilter struct {
	Timestamp string `json:"timestamp"` // The type of timestamp to filter on. Possible values are: "created_time", "last_edited_time".
}

func (f TimestampFilter) filter() {}

var _ Filter = TimestampFilter{}

/*
A compound filter object combines several filter objects together using a logical operator and or or. A compound filter can even be combined within a compound filter, but only up to two nesting levels deep.
The compound filter object contains one of the following keys:
*/
type CompoundFilter struct {
	Or  []Filter `json:"or,omitempty"`  // Returns pages when any of the filters inside the provided array match.
	And []Filter `json:"and,omitempty"` // Returns pages when all of the filters inside the provided array match.
}

func (f CompoundFilter) filter() {}

var _ Filter = CompoundFilter{}

// A text filter condition can be applied to database properties of types "title", "rich_text", "url", "email", and "phone_number".
type TextFilterCondition struct {
	Equals         string `json:"equals,omitempty"`           // Only return pages where the page property value matches the provided value exactly.
	DoesNotEqual   string `json:"does_not_equal,omitempty"`   // Only return pages where the page property value does not match the provided value exactly.
	Contains       string `json:"contains,omitempty"`         // Only return pages where the page property value contains the provided value.
	DoesNotContain string `json:"does_not_contain,omitempty"` // Only return pages where the page property value does not contain the provided value.
	StartsWith     string `json:"starts_with,omitempty"`      // Only return pages where the page property value starts with the provided value.
	EndsWith       string `json:"ends_with,omitempty"`        // Only return pages where the page property value ends with the provided value.
	IsEmpty        bool   `json:"is_empty,omitempty"`         // Only return pages where the page property value is empty.
	IsNotEmpty     bool   `json:"is_not_empty,omitempty"`     // Only return pages where the page property value is present.
}

// A number filter condition can be applied to database properties of type "number".
type NumberFilterCondition struct {
	Equals               *float64 `json:"equals,omitempty"`                   // Only return pages where the page property value matches the provided value exactly.
	DoesNotEqual         *float64 `json:"does_not_equal,omitempty"`           // Only return pages where the page property value does not match the provided value exactly.
	GreaterThan          *float64 `json:"greater_than,omitempty"`             // Only return pages where the page property value is greater than the provided value.
	LessThan             *float64 `json:"less_than,omitempty"`                // Only return pages where the page property value is less than the provided value.
	GreaterThanOrEqualTo *float64 `json:"greater_than_or_equal_to,omitempty"` // Only return pages where the page property value is greater than or equal to the provided value.
	LessThanOrEqualTo    *float64 `json:"less_than_or_equal_to,omitempty"`    // Only return pages where the page property value is less than or equal to the provided value.
	IsEmpty              bool     `json:"is_empty,omitempty"`                 // Only return pages where the page property value is empty.
	IsNotEmpty           bool     `json:"is_not_empty,omitempty"`             // Only return pages where the page property value is present.
}

// A checkbox filter condition can be applied to database properties of type "checkbox".
type CheckboxFilterCondition struct {
	Equals       *bool `json:"equals,omitempty"`         // Only return pages where the page property value matches the provided value exactly.
	DoesNotEqual *bool `json:"does_not_equal,omitempty"` // Only return pages where the page property value does not match the provided value exactly.
}

// A select filter condition can be applied to database properties of type "select".
type SelectFilterCondition struct {
	Equals       string `json:"equals,omitempty"`         // Only return pages where the page property value matches the provided value exactly.
	DoesNotEqual string `json:"does_not_equal,omitempty"` // Only return pages where the page property value does not match the provided value exactly.
	IsEmpty      bool   `json:"is_empty,omitempty"`       // Only return pages where the page property value is empty.
	IsNotEmpty   bool   `json:"is_not_empty,omitempty"`   // Only return pages where the page property value is present.
}

// A multi-select filter condition can be applied to database properties of type "multi_select".
type MultiSelectFilterCondition struct {
	Contains       string `json:"contains,omitempty"`         // Only return pages where the page property value contains the provided value.
	DoesNotContain string `json:"does_not_contain,omitempty"` // Only return pages where the page property value does not contain the provided value.
	IsEmpty        bool   `json:"is_empty,omitempty"`         // Only return pages where the page property value is empty.
	IsNotEmpty     bool   `json:"is_not_empty,omitempty"`     // Only return pages where the page property value is present.
}

// A status filter condition can be applied to database properties of type "status".
type StatusFilterCondition struct {
	Equals       string `json:"equals,omitempty"`         // Only return pages where the page property value matches the provided value exactly.
	DoesNotEqual string `json:"does_not_equal,omitempty"` // Only return pages where the page property value does not match the provided value exactly.
	IsEmpty      bool   `json:"is_empty,omitempty"`       // Only return pages where the page property value is empty.
	IsNotEmpty   bool   `json:"is_not_empty,omitempty"`   // Only return pages where the page property value is present.
}

// A date filter condition can be applied to database properties of types "date", "created_time", and "last_edited_time".
type DateFilterCondition struct {
	Equals     ISO8601String `json:"equals,omitempty"`       // Only return pages where the page property value matches the provided date exactly. If a date is provided, the comparison is done against the start and end of the UTC date.If a date with a time is provided, the comparison is done with millisecond precision.Note that if no timezone is provided, the default is UTC.
	Before     ISO8601String `json:"before,omitempty"`       // Only return pages where the page property value is before the provided date. If a date with a time is provided, the comparison is done with millisecond precision.Note that if no timezone is provided, the default is UTC.
	After      ISO8601String `json:"after,omitempty"`        // Only return pages where the page property value is after the provided date. If a date with a time is provided, the comparison is done with millisecond precision.Note that if no timezone is provided, the default is UTC.
	OnOrBefore ISO8601String `json:"on_or_before,omitempty"` // Only return pages where the page property value is on or before the provided date. If a date with a time is provided, the comparison is done with millisecond precision.Note that if no timezone is provided, the default is UTC.
	IsEmpty    bool          `json:"is_empty,omitempty"`     // Only return pages where the page property value is empty.
	IsNotEmpty bool          `json:"is_not_empty,omitempty"` // Only return pages where the page property value is present.
	OnOrAfter  ISO8601String `json:"on_or_after,omitempty"`  // Only return pages where the page property value is on or after the provided date. If a date with a time is provided, the comparison is done with millisecond precision.Note that if no timezone is provided, the default is UTC.
	PastWeek   struct{}      `json:"past_week,omitempty"`    // Only return pages where the page property value is within the past week.
	PastMonth  struct{}      `json:"past_month,omitempty"`   // Only return pages where the page property value is within the past month.
	PastYear   struct{}      `json:"past_year,omitempty"`    // Only return pages where the page property value is within the past year.
	NextWeek   struct{}      `json:"next_week,omitempty"`    // Only return pages where the page property value is within the next week.
	NextMonth  struct{}      `json:"next_month,omitempty"`   // Only return pages where the page property value is within the next month.
	NextYear   struct{}      `json:"next_year,omitempty"`    // Only return pages where the page property value is within the next year.
}

// A people filter condition can be applied to database properties of types "people",  "created_by", and "last_edited_by".
type PeopleFilterCondition struct {
	Contains       UUIDString `json:"contains,omitempty"`         // Only return pages where the page property value contains the provided value.
	DoesNotContain UUIDString `json:"does_not_contain,omitempty"` // Only return pages where the page property value does not contain the provided value.
	IsEmpty        bool       `json:"is_empty,omitempty"`         // Only return pages where the page property value is empty.
	IsNotEmpty     bool       `json:"is_not_empty,omitempty"`     // Only return pages where the page property value is present.
}

// A files filter condition can be applied to database properties of type "files".
type FilesFilterCondition struct {
	IsEmpty    bool `json:"is_empty,omitempty"`     // Only return pages where the page property value is empty.
	IsNotEmpty bool `json:"is_not_empty,omitempty"` // Only return pages where the page property value is present.
}

// A relation filter condition can be applied to database properties of type "relation".
type RelationFilterCondition struct {
	Contains       UUIDString `json:"contains,omitempty"`         // Only return pages where the page property value contains the provided value.
	DoesNotContain UUIDString `json:"does_not_contain,omitempty"` // Only return pages where the page property value does not contain the provided value.
	IsEmpty        bool       `json:"is_empty,omitempty"`         // Only return pages where the page property value is empty.
	IsNotEmpty     bool       `json:"is_not_empty,omitempty"`     // Only return pages where the page property value is present.
}

/*
A rollup filter condition can be applied to database properties of type "rollup". Rollups which evaluate to arrays accept a filter with an any, every, or none condition; rollups which evaluate to numbers accept a filter with a number condition; and rollups which evaluate to dates accept a filter with a date condition.
Rollups which evaluate to arrays accept any kind of property in
*/
type RollupFilterCondition struct {
	Any    interface{}           `json:"any,omitempty"`    // For a rollup property which evaluates to an array, return the pages where any item in that rollup fits this criterion. The criterion itself can be any other property type.
	Every  interface{}           `json:"every,omitempty"`  // For a rollup property which evaluates to an array, return the pages where every item in that rollup fits this criterion. The criterion itself can be any other property type.
	None   interface{}           `json:"none,omitempty"`   // For a rollup property which evaluates to an array, return the pages where no item in that rollup fits this criterion. The criterion itself can be any other property type.
	Number NumberFilterCondition `json:"number,omitempty"` // For a rollup property which evaluates to an number, return the pages where the number fits this criterion.
	Date   DateFilterCondition   `json:"date,omitempty"`   // For a rollup property which evaluates to an date, return the pages where the date fits this criterion.
}

// A formula filter condition can be applied to database properties of type "formula".
type FormulaFilterCondition struct {
	String   TextFilterCondition     `json:"string,omitempty"`   // Only return pages where the result type of the page property formula is "string" and the provided string filter condition matches the formula's value.
	Checkbox CheckboxFilterCondition `json:"checkbox,omitempty"` // Only return pages where the result type of the page property formula is "checkbox" and the provided checkbox filter condition matches the formula's value.
	Number   NumberFilterCondition   `json:"number,omitempty"`   // Only return pages where the result type of the page property formula is "number" and the provided number filter condition matches the formula's value.
	Date     DateFilterCondition     `json:"date,omitempty"`     // Only return pages where the result type of the page property formula is "date" and the provided date filter condition matches the formula's value.
}
