package notionapi

// Code generated by notion.codegen; DO NOT EDIT.

// https://developers.notion.com/reference/pagination

// Responses from paginated endpoints contain the following properties:
type Pagination struct {
	HasMore    bool   `json:"has_more"`    // When the response includes the end of the list, false. Otherwise, true.
	NextCursor string `json:"next_cursor"` // Only available when has_more is true.Used to retrieve the next page of results by passing the value as the start_cursor parameter to the same endpoint.
	Object     string `json:"object"`      // Always list.
	Type       string `json:"type"`        // Type of the objects in results. Possible values include "block", "page", "user", "database", "property_item", "page_or_database".
}

type BlockPagination struct {
	Pagination
	Results []Block `json:"results"`
}

type PagePagination struct {
	Pagination
	Results []Page `json:"results"`
}

type UserPagination struct {
	Pagination
	Results []User `json:"results"`
}

type DatabasePagination struct {
	Pagination
	Results []Database `json:"results"`
}

type PropertyItemPagination struct {
	Pagination
	Results PropertyItems `json:"results"`
}

type PageOrDatabasePagination struct {
	Pagination
	Results []PageOrDatabase `json:"results"`
}
