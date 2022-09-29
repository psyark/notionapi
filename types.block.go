package notionapi

import "encoding/json"

// Code generated by notion.codegen; DO NOT EDIT.

// https://developers.notion.com/reference/block

// Each block object contains the following keys. In addition, it must contain a key corresponding with the value of type. The value is an object containing a type-specific block information. The type-specific block information is described in the sections below.
type Block struct {
	Object           string                    `json:"object"`                             // Always "block".
	ID               UUIDString                `json:"id"`                                 // Identifier for the block.
	Parent           *Parent                   `json:"parent"`                             // Information about the block's parent. See Parent object.
	Type             string                    `json:"type"`                               // Type of block. Possible values include "paragraph", "heading_1", "heading_2", "heading_3", "bulleted_list_item", "numbered_list_item", "to_do", "toggle", "child_page","child_database", "embed", "image", "video", "file", "pdf", "bookmark", "callout",  "quote", "equation", "divider", "table_of_contents", "column", "column_list", "link_preview", "synced_block", "template", "link_to_page", "table"' "table_row", and "unsupported".
	CreatedTime      ISO8601String             `json:"created_time"`                       // Date and time when this block was created. Formatted as an ISO 8601 date time string.
	CreatedBy        *PartialUser              `json:"created_by"`                         // User who created the block.
	LastEditedTime   ISO8601String             `json:"last_edited_time"`                   // Date and time when this block was last updated. Formatted as an ISO 8601 date time string.
	LastEditedBy     *PartialUser              `json:"last_edited_by"`                     // User who last edited the block.
	Archived         bool                      `json:"archived"`                           // The archived status of the block.
	HasChildren      bool                      `json:"has_children"`                       // Whether or not the block has children blocks nested within it.
	Paragraph        ParagraphBlockData        `json:"paragraph" specific:"type"`          // Paragraph block objects contain the following information within the paragraph property:
	Heading1         HeadingBlockData          `json:"heading_1" specific:"type"`          // Heading one block objects contain the following information within the heading_1 property:
	Heading2         HeadingBlockData          `json:"heading_2" specific:"type"`          // Heading two block objects contain the following information within the heading_2 property:
	Heading3         HeadingBlockData          `json:"heading_3" specific:"type"`          // Heading three block objects contain the following information within the heading_3 property:
	Callout          CalloutBlockData          `json:"callout" specific:"type"`            // Callout block objects contain the following information within the callout property:
	Quote            QuoteBlockData            `json:"quote" specific:"type"`              // Quote block objects contain the following information within the quote property
	BulletedListItem BulletedListItemBlockData `json:"bulleted_list_item" specific:"type"` // Bulleted list item block objects contain the following information within the bulleted_list_item property:
	NumberedListItem NumberedListItemBlockData `json:"numbered_list_item" specific:"type"` // Numbered list item block objects contain the following information within the numbered_list_item property:
	ToDo             ToDoBlockData             `json:"to_do" specific:"type"`              // To do block objects contain the following information within the to_do property:
	Toggle           ToggleBlockData           `json:"toggle" specific:"type"`             // Toggle block objects contain the following information within the toggle property:
	Code             CodeBlockData             `json:"code" specific:"type"`               // Code block objects contain the following information within the code property:
	ChildPage        ChildPageBlockData        `json:"child_page" specific:"type"`         // Child page block objects contain the following information within the child_page property:
	ChildDatabase    ChildDatabaseBlockData    `json:"child_database" specific:"type"`     // Child database block objects contain the following information within the child_database property:
	Embed            EmbedBlockData            `json:"embed" specific:"type"`              // Embed blocks include block types that allow displaying another website within Notion. These block types are:  * Framer * Twitter (tweets) * Google Drive documents * Gist * Figma * Invision, * Loom * Typeform * Codepen * PDFs * Google Maps * Whimisical * Miro * Abstract * excalidraw * Sketch * Replit   There is no need to specify the specific embed type, only the URL.  Embed block objects contain the following information within the embed property:
	Image            *ImageFile                `json:"image" specific:"type"`              // Includes supported image urls (i.e. ending in .png, .jpg, .jpeg, .gif, .tif, .tiff, .bmp, .svg, or .heic). Note that the url property only accepts direct urls to an image. The image must be directly hosted. In other words, the url cannot point to a service that retrieves the image.
	Video            *ImageFile                `json:"video" specific:"type"`              // Includes supported video urls (e.g. ending in .mkv, .flv, .gifv, .avi, .mov, .qt, .wmv, .asf, .amv, .mp4, .m4v, .mpeg, .mpv, .mpg, .f4v, etc.)
	File             FileBlockData             `json:"file" specific:"type"`
	PDF              PDFBlockData              `json:"pdf" specific:"type"`
	Bookmark         BookmarkBlockData         `json:"bookmark" specific:"type"`          // Bookmark block objects contain the following information within the bookmark property:
	Equation         EquationBlockData         `json:"equation" specific:"type"`          // Equation block objects contain the following information within the equation property
	Divider          struct{}                  `json:"divider" specific:"type"`           // Divider block objects do not contain any information within the divider property
	TableOfContents  TableOfContentsBlockData  `json:"table_of_contents" specific:"type"` // Table of contents block objects contain the following information within the table_of_contents property
	Breadcrumb       struct{}                  `json:"breadcrumb" specific:"type"`        // Breadcrumb block objects do not contain any information within the breadcrumb  property
	ColumnList       *ColumnListBlocks         `json:"column_list" specific:"type"`       // Column Lists are parent blocks for column children. They do not contain any information within the column_list property and can only contain children of type column.  Columns are parent blocks for any supported block children, excluding columns. They do not contain any information within the column property. They can only be appended to column_lists.  When creating a column list block via Append block children, the column_list must have at least 2 columns, and those columns must have at least one child each.  When fetching content for a column_list, first fetch the the column children via Retrieve block children. Then fetch the children for each column block. Column List blocks contain the following information in the column_list property:  Column blocks contain the following information in the column property.
	Column           *ColumnBlocks             `json:"column" specific:"type"`            // Column Lists are parent blocks for column children. They do not contain any information within the column_list property and can only contain children of type column.  Columns are parent blocks for any supported block children, excluding columns. They do not contain any information within the column property. They can only be appended to column_lists.  When creating a column list block via Append block children, the column_list must have at least 2 columns, and those columns must have at least one child each.  When fetching content for a column_list, first fetch the the column children via Retrieve block children. Then fetch the children for each column block. Column List blocks contain the following information in the column_list property:  Column blocks contain the following information in the column property.
	LinkPreview      LinkPreviewBlockData      `json:"link_preview" specific:"type"`      // Link Preview block objects return the originally pasted url. NOTE: The link_preview block will only be returned as part of a response. It cannot be created via the API.
	Template         TemplateBlockData         `json:"template" specific:"type"`          // Template block objects contain the following information within the template property:
	LinkToPage       LinkToPageBlockData       `json:"link_to_page" specific:"type"`      // Link to page objects contain a key corresponding with the value of type. The value is a type-specific string as described below.
	SyncedBlock      *SyncedBlockBlocks        `json:"synced_block" specific:"type"`      // "References" synced block objects contain the following information within the synced_block property:
	Table            TableBlockData            `json:"table" specific:"type"`             // Tables are parent blocks for table row children. They can only contain children of type table_row.  When creating a table block via the Append block children endpoint, the table must have at least 1 table_row whose cells array has the same length as the table_width. To fetch content for a table, fetch the the table_row children via Retrieve block children. The table block itself only contains formatting data, no content. Table blocks contain the following within the table property:
	TableRow         TableRowBlockData         `json:"table_row" specific:"type"`         // Table row blocks contain the following within the table_row property:
}

func (p Block) MarshalJSON() ([]byte, error) {
	type Alias Block
	return marshalByType(Alias(p), p.Type)
}

// Paragraph block objects contain the following information within the paragraph property:
type ParagraphBlockData struct {
	RichText []RichText `json:"rich_text"`          // Rich text in the paragraph block.
	Color    string     `json:"color"`              // Color of the block. Possible values are: "default", "gray", "brown", "orange", "yellow", "green", "blue", "purple", "pink", "red", "gray_background", "brown_background", "orange_background", "yellow_background", "green_background", "blue_background", "purple_background", "pink_background", "red_background".
	Children []Block    `json:"children,omitempty"` // Any nested children blocks of the paragraph block.
}

/*
Heading one block objects contain the following information within the heading_1 property:
Heading two block objects contain the following information within the heading_2 property:
Heading three block objects contain the following information within the heading_3 property:
*/
type HeadingBlockData struct {
	RichText     []RichText `json:"rich_text"`     // Rich text in the heading block.
	Color        string     `json:"color"`         // Color of the block. Possible values are: "default", "gray", "brown", "orange", "yellow", "green", "blue", "purple", "pink", "red", "gray_background", "brown_background", "orange_background", "yellow_background", "green_background", "blue_background", "purple_background", "pink_background", "red_background".
	IsToggleable bool       `json:"is_toggleable"` // Whether or not the heading block is a toggle heading or not. If true, the heading block has toggle and can support children. If false, the heading block is a normal heading block.
}

// Callout block objects contain the following information within the callout property:
type CalloutBlockData struct {
	RichText []RichText  `json:"rich_text"`          // Rich text in the heading block.
	Icon     FileOrEmoji `json:"icon"`               // Page icon.
	Color    string      `json:"color"`              // Color of the block. Possible values are: "default", "gray", "brown", "orange", "yellow", "green", "blue", "purple", "pink", "red", "gray_background", "brown_background", "orange_background", "yellow_background", "green_background", "blue_background", "purple_background", "pink_background", "red_background".
	Children []Block     `json:"children,omitempty"` // Any nested children blocks of the callout block.
}

func (p *CalloutBlockData) UnmarshalJSON(data []byte) error {
	p.Icon = newFileOrEmoji(getChild(data, "icon"))
	type Alias CalloutBlockData
	return json.Unmarshal(data, (*Alias)(p))
}

// Quote block objects contain the following information within the quote property
type QuoteBlockData struct {
	RichText []RichText `json:"rich_text"`          // Rich text in the quote block.
	Color    string     `json:"color"`              // Color of the block. Possible values are: "default", "gray", "brown", "orange", "yellow", "green", "blue", "purple", "pink", "red", "gray_background", "brown_background", "orange_background", "yellow_background", "green_background", "blue_background", "purple_background", "pink_background", "red_background".
	Children []Block    `json:"children,omitempty"` // Any nested children blocks of the quote block.
}

// Bulleted list item block objects contain the following information within the bulleted_list_item property:
type BulletedListItemBlockData struct {
	RichText []RichText `json:"rich_text"`          // Rich text in the bulleted_list_item block.
	Color    string     `json:"color"`              // Color of the block. Possible values are: "default", "gray", "brown", "orange", "yellow", "green", "blue", "purple", "pink", "red", "gray_background", "brown_background", "orange_background", "yellow_background", "green_background", "blue_background", "purple_background", "pink_background", "red_background".
	Children []Block    `json:"children,omitempty"` // Any nested children blocks of the bulleted_list_item block.
}

// Numbered list item block objects contain the following information within the numbered_list_item property:
type NumberedListItemBlockData struct {
	RichText []RichText `json:"rich_text"`          // Rich text in the numbered_list_item block.
	Color    string     `json:"color"`              // Color of the block. Possible values are: "default", "gray", "brown", "orange", "yellow", "green", "blue", "purple", "pink", "red", "gray_background", "brown_background", "orange_background", "yellow_background", "green_background", "blue_background", "purple_background", "pink_background", "red_background".
	Children []Block    `json:"children,omitempty"` // Any nested children blocks of the numbered_list_item block.
}

// To do block objects contain the following information within the to_do property:
type ToDoBlockData struct {
	RichText []RichText `json:"rich_text"`          // Rich text in the to_do block.
	Checked  bool       `json:"checked"`            // Whether the to_do is checked or not.
	Color    string     `json:"color"`              // Color of the block. Possible values are: "default", "gray", "brown", "orange", "yellow", "green", "blue", "purple", "pink", "red", "gray_background", "brown_background", "orange_background", "yellow_background", "green_background", "blue_background", "purple_background", "pink_background", "red_background".
	Children []Block    `json:"children,omitempty"` // Any nested children blocks of the to_do block.
}

// Toggle block objects contain the following information within the toggle property:
type ToggleBlockData struct {
	RichText []RichText `json:"rich_text"`          // Rich text in the toggle block.
	Color    string     `json:"color"`              // Color of the block. Possible values are: "default", "gray", "brown", "orange", "yellow", "green", "blue", "purple", "pink", "red", "gray_background", "brown_background", "orange_background", "yellow_background", "green_background", "blue_background", "purple_background", "pink_background", "red_background".
	Children []Block    `json:"children,omitempty"` // Any nested children blocks of the toggle block.
}

// Code block objects contain the following information within the code property:
type CodeBlockData struct {
	RichText []RichText `json:"rich_text"` // Rich text in code block
	Caption  []RichText `json:"caption"`   // Rich text in caption of the code block
	Language string     `json:"language"`  // Coding language in code block
}

// Child page block objects contain the following information within the child_page property:
type ChildPageBlockData struct {
	Title string `json:"title"` // Plain text of page title.
}

// Child database block objects contain the following information within the child_database property:
type ChildDatabaseBlockData struct {
	Title string `json:"title"` // Plain text of the database title
}

/*
Embed blocks include block types that allow displaying another website within Notion.
These block types are:
* Framer
* Twitter (tweets)
* Google Drive documents
* Gist
* Figma
* Invision,
* Loom
* Typeform
* Codepen
* PDFs
* Google Maps
* Whimisical
* Miro
* Abstract
* excalidraw
* Sketch
* Replit

There is no need to specify the specific embed type, only the URL.

Embed block objects contain the following information within the embed property:
*/
type EmbedBlockData struct {
	URL     string     `json:"url"`     // Link to website the embed block will display.
	Caption []RichText `json:"caption"` // undocumented
}

type FileBlockData struct {
	File    *File      `json:"file"`    // File reference
	Caption []RichText `json:"caption"` // Caption of the file block
}

type PDFBlockData struct {
	PDF *File `json:"pdf"` // PDF file reference
}

// Bookmark block objects contain the following information within the bookmark property:
type BookmarkBlockData struct {
	URL     string     `json:"url"`     // Bookmark link
	Caption []RichText `json:"caption"` // Caption of the bookmark block
}

// Equation block objects contain the following information within the equation property
type EquationBlockData struct {
	Expression string `json:"expression"` // A KaTeX compatible string
}

// Table of contents block objects contain the following information within the table_of_contents property
type TableOfContentsBlockData struct {
	Color string `json:"color"` // Color of the block. Possible values are: "default", "gray", "brown", "orange", "yellow", "green", "blue", "purple", "pink", "red", "gray_background", "brown_background", "orange_background", "yellow_background", "green_background", "blue_background", "purple_background", "pink_background", "red_background".
}

/*
Column Lists are parent blocks for column children. They do not contain any information within the column_list property and can only contain children of type column.

Columns are parent blocks for any supported block children, excluding columns. They do not contain any information within the column property. They can only be appended to column_lists.

When creating a column list block via Append block children, the column_list must have at least 2 columns, and those columns must have at least one child each.

When fetching content for a column_list, first fetch the the column children via Retrieve block children. Then fetch the children for each column block.
Column List blocks contain the following information in the column_list property:

Column blocks contain the following information in the column property.
*/
type ColumnListBlocks struct{}

/*
Column Lists are parent blocks for column children. They do not contain any information within the column_list property and can only contain children of type column.

Columns are parent blocks for any supported block children, excluding columns. They do not contain any information within the column property. They can only be appended to column_lists.

When creating a column list block via Append block children, the column_list must have at least 2 columns, and those columns must have at least one child each.

When fetching content for a column_list, first fetch the the column children via Retrieve block children. Then fetch the children for each column block.
Column List blocks contain the following information in the column_list property:

Column blocks contain the following information in the column property.
*/
type ColumnBlocks struct{}

// Link Preview block objects return the originally pasted url. NOTE: The link_preview block will only be returned as part of a response. It cannot be created via the API.
type LinkPreviewBlockData struct{}

// Template block objects contain the following information within the template property:
type TemplateBlockData struct {
	RichText []RichText `json:"rich_text"`          // rich text in the title of the template
	Children []Block    `json:"children,omitempty"` // Any nested children blocks of the template block. These blocks will be duplicated when the template block is used in the UI.
}

// Link to page objects contain a key corresponding with the value of type. The value is a type-specific string as described below.
type LinkToPageBlockData struct {
	Type       string     `json:"type"`        // Type of this link to page object. Possible values are: "page_id", "database_id".
	PageID     UUIDString `json:"page_id"`     // Identifier for a page
	DatabaseID UUIDString `json:"database_id"` // Identifier for a database page
}

// "References" synced block objects contain the following information within the synced_block property:
type SyncedBlockBlocks struct {
	SyncedFrom *SyncedFrom `json:"synced_from"` // Object that contains the id of the original synced_block
}

/*
Synced From Object
Synced From objects contain a key corresponding with the value of type. The value is a type-specific string  as described below.
*/
type SyncedFrom struct {
	Type    string     `json:"type"`     // Type of this synced from object. Possible values are: "block_id".
	BlockID UUIDString `json:"block_id"` // Identifier of an original synced_block
}

/*
Tables are parent blocks for table row children. They can only contain children of type table_row.

When creating a table block via the Append block children endpoint, the table must have at least 1 table_row whose cells array has the same length as the table_width.
To fetch content for a table, fetch the the table_row children via Retrieve block children. The table block itself only contains formatting data, no content.
Table blocks contain the following within the table property:
*/
type TableBlockData struct {
	TableWidth      int64   `json:"table_width"`        // Number of columns in the table. Note that this cannot be changed via the public API once a table is created.
	HasColumnHeader bool    `json:"has_column_header"`  // Whether or not the table has a column header. If true, the first row in the table will appear visually distinct from the other rows.
	HasRowHeader    bool    `json:"has_row_header"`     // Whether or not the table has a header row. If true, the first column in the table will appear visually distinct from the other columns.
	Children        []Block `json:"children,omitempty"` // List of table_row children for this table.
}

// Table row blocks contain the following within the table_row property:
type TableRowBlockData struct {
	Cells [][]RichText `json:"cells"` // Array of cell contents in horizontal display order. Each cell itself is an array of rich text objects.
}
