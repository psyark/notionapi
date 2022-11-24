package notionapi

import uuid "github.com/google/uuid"

// Code generated by notion.codegen; DO NOT EDIT.

// https://developers.notion.com/reference/property-object

// Metadata that controls how a database property behaves.
type Property struct {
	// Each database property object contains the following keys. In addition, it must contain a key corresponding with the value of type. The value is an object containing type-specific configuration. The type-specific configurations are described in the sections below.
	ID   string `json:"id"`   // The ID of the property, usually a short string of random letters and symbols. Some automatically generated property types have special human-readable IDs. For example, all Title properties have an ID of "title".
	Type string `json:"type"` // Type that controls the behavior of the property. Possible values are: "title", "rich_text", "number", "select", "multi_select", "date", "people", "files", "checkbox", "url", "email", "phone_number", "formula", "relation", "rollup", "created_time", "created_by", "last_edited_time", "last_edited_by", "status".
	Name string `json:"name"` // The name of the property as it appears in Notion.

	Title          struct{}                  `json:"title" specific:"type"`            // Each database must have exactly one database property of type "title". This database property controls the title that appears at the top of the page when the page is opened. Title database property objects have no additional configuration within the title property.
	RichText       struct{}                  `json:"rich_text" specific:"type"`        // Text database property objects have no additional configuration within the rich_text property.
	Number         *NumberConfiguration      `json:"number" specific:"type"`           // Number database property objects contain the following configuration within the number property:
	Select         *SelectConfiguration      `json:"select" specific:"type"`           // Select database property objects contain the following configuration within the select property:
	Status         *StatusConfiguration      `json:"status" specific:"type"`           // Status database property objects contain the following configuration within the status property:
	MultiSelect    *MultiSelectConfiguration `json:"multi_select" specific:"type"`     // Multi-select database property objects contain the following configuration within the multi_select property:
	Date           struct{}                  `json:"date" specific:"type"`             // Date database property objects have no additional configuration within the date property.
	People         struct{}                  `json:"people" specific:"type"`           // People database property objects have no additional configuration within the people property.
	Files          struct{}                  `json:"files" specific:"type"`            // Files database property objects have no additional configuration within the files property.
	Checkbox       struct{}                  `json:"checkbox" specific:"type"`         // Checkbox database property objects have no additional configuration within the checkbox property.
	URL            struct{}                  `json:"url" specific:"type"`              // URL database property objects have no additional configuration within the url property.
	Email          struct{}                  `json:"email" specific:"type"`            // Email database property objects have no additional configuration within the email property.
	PhoneNumber    struct{}                  `json:"phone_number" specific:"type"`     // Phone number database property objects have no additional configuration within the phone_number property.
	Formula        *FormulaConfiguration     `json:"formula" specific:"type"`          // Formula database property objects contain the following configuration within the formula property:
	Relation       *RelationConfiguration    `json:"relation" specific:"type"`         // Relation database property objects contain the following configuration within the relation property. In addition, they must contain a key corresponding with the value of type. The value is an object containing type-specific configuration. The type-specific configurations are defined below.
	Rollup         *RollupConfiguration      `json:"rollup" specific:"type"`           // Rollup database property objects contain the following configuration within the rollup property:
	CreatedTime    struct{}                  `json:"created_time" specific:"type"`     // Created time database property objects have no additional configuration within the created_time property.
	CreatedBy      struct{}                  `json:"created_by" specific:"type"`       // Created by database property objects have no additional configuration within the created_by property.
	LastEditedTime struct{}                  `json:"last_edited_time" specific:"type"` // Last edited time database property objects have no additional configuration within the last_edited_time property.
	LastEditedBy   struct{}                  `json:"last_edited_by" specific:"type"`   // Last edited by database property objects have no additional configuration within the last_edited_by property.
}

func (p Property) MarshalJSON() ([]byte, error) {
	type Alias Property
	return marshalByType(Alias(p), p.Type)
}

type PropertyMap map[string]Property

func (m PropertyMap) Get(id string) *Property {
	for _, pv := range m {
		if pv.ID == id {
			return &pv
		}
	}
	return nil
}

// Number database property objects contain the following configuration within the number property:
type NumberConfiguration struct {
	Format string `json:"format"` // How the number is displayed in Notion. Potential values include: number, number_with_commas, percent, dollar, canadian_dollar, euro, pound, yen, ruble, rupee, won, yuan, real, lira, rupiah, franc, hong_kong_dollar, new_zealand_dollar, krona, norwegian_krone, mexican_peso, rand, new_taiwan_dollar, danish_krone, zloty, baht, forint, koruna, shekel, chilean_peso, philippine_peso, dirham, colombian_peso, riyal, ringgit, leu, argentine_peso, uruguayan_peso, singapore_dollar.
}

// Select database property objects contain the following configuration within the select property:
type SelectConfiguration struct {
	Options []SelectOption `json:"options"` // Sorted list of options available for this property.
}

type SelectOption struct {
	Name  string `json:"name,omitempty"`  // Name of the option as it appears in Notion.  Note: Commas (",") are not valid for select values.
	ID    string `json:"id,omitempty"`    // Identifier of the option, which does not change if the name is changed. These are sometimes, but not always, UUIDs.
	Color string `json:"color,omitempty"` // Color of the option. Possible values include: default, gray, brown, orange, yellow, green, blue, purple, pink, red.
}

// Status database property objects contain the following configuration within the status property:
type StatusConfiguration struct {
	Options []StatusOption `json:"options"` // Sorted list of options available for this property.
	Groups  []StatusGroup  `json:"groups"`  // Sorted list of groups available for this property.
}

type StatusOption struct {
	Name  string `json:"name,omitempty"`  // Name of the option as it appears in Notion.  Note: Commas (",") are not valid for select values.
	ID    string `json:"id,omitempty"`    // Identifier of the option, which does not change if the name is changed. These are sometimes, but not always, UUIDs.
	Color string `json:"color,omitempty"` // Color of the option. Possible values include: default, gray, brown, orange, yellow, green, blue, purple, pink, red.
}

type StatusGroup struct {
	Name      string      `json:"name"`       // Name of the option as it appears in Notion.  Note: Commas (",") are not valid for select values.
	ID        string      `json:"id"`         // Identifier of the option, which does not change if the name is changed. These are sometimes, but not always, UUIDs.
	Color     string      `json:"color"`      // Color of the option. Possible values include: default, gray, brown, orange, yellow, green, blue, purple, pink, red.
	OptionIds []uuid.UUID `json:"option_ids"` // Sorted list of ids of all options that belong to a group.
}

// Multi-select database property objects contain the following configuration within the multi_select property:
type MultiSelectConfiguration struct {
	Options []SelectOption `json:"options"` // Settings for multi select properties.
}

// Formula database property objects contain the following configuration within the formula property:
type FormulaConfiguration struct {
	Expression string `json:"expression"` // Formula to evaluate for this property. You can read more about the syntax for formulas in the help center.
}

// Relation database property objects contain the following configuration within the relation property. In addition, they must contain a key corresponding with the value of type. The value is an object containing type-specific configuration. The type-specific configurations are defined below.
type RelationConfiguration struct {
	DatabaseID     uuid.UUID                          `json:"database_id"`                     // The database this relation refers to. New linked pages must belong to this database in order to be valid.
	Type           string                             `json:"type"`                            // The type of the relation. Can be "single_property" or "dual_property".
	SingleProperty struct{}                           `json:"single_property" specific:"type"` // Single property relation objects have no additional configuration within the single_property property.
	DualProperty   *DualPropertyRelationConfiguration `json:"dual_property" specific:"type"`   // Dual property relation objects contain the following configuration within the dual_property property:
}

func (p RelationConfiguration) MarshalJSON() ([]byte, error) {
	type Alias RelationConfiguration
	return marshalByType(Alias(p), p.Type)
}

// Dual property relation objects contain the following configuration within the dual_property property:
type DualPropertyRelationConfiguration struct {
	SyncedPropertyName string `json:"synced_property_name"` // The relation is formed as two synced properties. If you make a change to one property, it updates the other property at the same time. synced_property_name  refers to the name  of the related property.
	SyncedPropertyID   string `json:"synced_property_id"`   // The relation is formed as two synced properties. If you make a change to one property, it updates the other property at the same time. synced_property_id refers to the id  of the related property. This is usually a short string of random letters and symbols.
}

// Rollup database property objects contain the following configuration within the rollup property:
type RollupConfiguration struct {
	RelationPropertyName string `json:"relation_property_name"` // The name of the relation property this property is responsible for rolling up.
	RelationPropertyID   string `json:"relation_property_id"`   // The id of the relation property this property is responsible for rolling up.
	RollupPropertyName   string `json:"rollup_property_name"`   // The name of the property of the pages in the related database that is used as an input to function.
	RollupPropertyID     string `json:"rollup_property_id"`     // The id of the property of the pages in the related database that is used as an input to function.
	Function             string `json:"function"`               // The function that is evaluated for every page in the relation of the rollup. Possible values include: count,  count_values,  empty,  not_empty,  unique,  show_unique,  percent_empty,  percent_not_empty,  sum,  average,  median,  min,  max,  range,  earliest_date,  latest_date,  date_range,  checked,  unchecked,  percent_checked,  percent_unchecked,  count_per_group,  percent_per_group,  show_original
}
