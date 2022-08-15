package notionapi

import (
	"context"
	"fmt"
)

// Code generated by notion.codegen; DO NOT EDIT.

// Query a database
// https://developers.notion.com/reference/post-database-query
func (c *Client) QueryDatabase(ctx context.Context, database_id string, options *QueryDatabaseOptions) (*Pagination, error) {
	result := &Pagination{}
	return result, c.call(ctx, "POST", fmt.Sprintf("/v1/databases/%v/query", database_id), options, result)
}

type QueryDatabaseOptions struct {
	Filter      map[string]interface{} `json:"filter,omitempty"`       // When supplied, limits which pages are returned based on the [filter conditions](ref:post-database-query-filter).
	Sorts       []interface{}          `json:"sorts,omitempty"`        // When supplied, orders the results based on the provided [sort criteria](ref:post-database-query-sort).
	StartCursor string                 `json:"start_cursor,omitempty"` // When supplied, returns a page of results starting after the cursor provided. If not supplied, this endpoint will return the first page of results.
	PageSize    int                    `json:"page_size,omitempty"`    // The number of items from the full list desired in the response. Maximum: 100
}

// Create a database
// https://developers.notion.com/reference/create-a-database
func (c *Client) CreateDatabase(ctx context.Context, options *CreateDatabaseOptions) (*Database, error) {
	result := &Database{}
	return result, c.call(ctx, "POST", fmt.Sprintf("/v1/databases"), options, result)
}

type CreateDatabaseOptions struct {
	Parent     *Parent                `json:"parent"`          // A [page parent](/reference/database#page-parent)
	Title      []RichText             `json:"title,omitempty"` // Title of database as it appears in Notion. An array of [rich text objects](ref:rich-text).
	Properties map[string]interface{} `json:"properties"`      // Property schema of database. The keys are the names of properties as they appear in Notion and the values are [property schema objects](https://developers.notion.com/reference/property-schema-object).
}

// Update database
// https://developers.notion.com/reference/update-a-database
func (c *Client) UpdateDatabase(ctx context.Context, database_id string, options *UpdateDatabaseOptions) (*Database, error) {
	result := &Database{}
	return result, c.call(ctx, "PATCH", fmt.Sprintf("/v1/databases/%v", database_id), options, result)
}

type UpdateDatabaseOptions struct {
	Title      []RichText             `json:"title,omitempty"`      // Title of database as it appears in Notion. An array of [rich text objects](ref:rich-text). If omitted, the database title will remain unchanged.
	Properties map[string]interface{} `json:"properties,omitempty"` // Updates to the property schema of a database. If updating an existing property, the keys are the names or IDs of the properties as they appear in Notion and the values are [property schema objects](ref:property-schema-object). If adding a new property, the key is the name of the database property and the value is a [property schema object](ref:property-schema-object).
}

// Retrieve a database
// https://developers.notion.com/reference/retrieve-a-database
func (c *Client) RetrieveDatabase(ctx context.Context, database_id string) (*Database, error) {
	result := &Database{}
	return result, c.call(ctx, "GET", fmt.Sprintf("/v1/databases/%v", database_id), nil, result)
}

// Retrieve a page
// https://developers.notion.com/reference/retrieve-a-page
func (c *Client) RetrievePage(ctx context.Context, page_id string) (*Page, error) {
	result := &Page{}
	return result, c.call(ctx, "GET", fmt.Sprintf("/v1/pages/%v", page_id), nil, result)
}

// Create a page
// https://developers.notion.com/reference/post-page
func (c *Client) CreatePage(ctx context.Context, options *CreatePageOptions) (*Page, error) {
	result := &Page{}
	return result, c.call(ctx, "POST", fmt.Sprintf("/v1/pages"), options, result)
}

type CreatePageOptions struct {
	Parent     *Parent                `json:"parent"`             // A [database parent](/reference/page#database-parent) or [page parent](/reference/page#page-parent)
	Properties map[string]interface{} `json:"properties"`         // Property values of this page. The keys are the names or IDs of the [property](ref:database#property-object) and the values are [property values](property-value-object)
	Children   []interface{}          `json:"children,omitempty"` // Page content for the new page as an array of [block objects](ref:block)
	Icon       map[string]interface{} `json:"icon,omitempty"`     // Page icon for the new page
	Cover      map[string]interface{} `json:"cover,omitempty"`    // Page cover for the new page
}

// Update page
// https://developers.notion.com/reference/patch-page
func (c *Client) UpdatePage(ctx context.Context, page_id string, options *UpdatePageOptions) (*Page, error) {
	result := &Page{}
	return result, c.call(ctx, "PATCH", fmt.Sprintf("/v1/pages/%v", page_id), options, result)
}

type UpdatePageOptions struct {
	Properties map[string]interface{} `json:"properties,omitempty"` // Property values to update for this page. The keys are the names or IDs of the [property](ref:database#property-object) and the values are [property values](ref:page#property-value-object).
	Archived   *bool                  `json:"archived,omitempty"`   // Set to true to archive (delete) a page. Set to false to un-archive (restore) a page.
	Icon       map[string]interface{} `json:"icon,omitempty"`       // Page icon for the new page.
	Cover      map[string]interface{} `json:"cover,omitempty"`      // Page cover for the new page
}
