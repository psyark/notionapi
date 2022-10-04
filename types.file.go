package notionapi

// Code generated by notion.codegen; DO NOT EDIT.

// https://developers.notion.com/reference/file-object

// File objects contain data about files uploaded to Notion as well as external files linked in Notion.
type File struct {
	// Each file object contains the following keys. In addition, it must contain a key corresponding with the value of type. The value is an object containing type-specific configuration. The type-specific configurations are described in the sections below.
	Type    string         `json:"type"`              // Type of this file object. Possible values are: "external", "file".
	Caption *RichTextArray `json:"caption,omitempty"` // undocumented

	File     *FilesUploadedToNotionData `json:"file" specific:"type"`     // All files hosted by Notion have a type of "file".  File objects contain the following information within the file property:
	External *ExternalFileData          `json:"external" specific:"type"` // All external file objects have a type of "external".  An external file is any URL that isn't hosted by Notion.  External file objects contain the following information within the external property:
}

func (c *File) fileOrEmoji() {}
func (p File) MarshalJSON() ([]byte, error) {
	type Alias File
	return marshalByType(Alias(p), p.Type)
}

/*
All files hosted by Notion have a type of "file".

File objects contain the following information within the file property:
*/
type FilesUploadedToNotionData struct {
	URL        string        `json:"url"`         // Authenticated S3 URL to the file. The file URL will be valid for 1 hour but updated links can be requested if required.
	ExpiryTime ISO8601String `json:"expiry_time"` // Date and time when this will expire. Formatted as an ISO 8601 date time string.
}

/*
All external file objects have a type of "external".

An external file is any URL that isn't hosted by Notion.

External file objects contain the following information within the external property:
*/
type ExternalFileData struct {
	URL string `json:"url"` // Link to the externally hosted content.
}
