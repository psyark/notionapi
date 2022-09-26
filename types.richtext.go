package notionapi

// Code generated by notion.codegen; DO NOT EDIT.

// https://developers.notion.com/reference/rich-text

// Rich text objects contain data for displaying formatted text, mentions, and equations. A rich text object also contains annotations for style information. Arrays of rich text objects are used within property objects and property value objects to create what a user sees as a single text value in Notion.
type RichText struct {
	// Each rich text object contains the following keys. In addition, it must contain a key corresponding with the value of type. The value is an object containing type-specific configuration. The type-specific configurations are described in the sections below.
	PlainText   string       `json:"plain_text"`  // The plain text without annotations.
	Href        *string      `json:"href"`        // The URL of any link or internal Notion mention in this text, if any.
	Annotations *Annotations `json:"annotations"` // All annotations that apply to this rich text. Annotations include colors and bold/italics/underline/strikethrough.
	Type        string       `json:"type"`        // Type of this rich text object. Possible values are: "text", "mention", "equation".

	Text     *Text     `json:"text" specific:"type"`     // Text objects contain the following information within the text property:
	Mention  *Mention  `json:"mention" specific:"type"`  // Mention objects represent an inline mention of a user, page, database, or date. In the app these are created by typing @ followed by the name of a user, page, database, or a date.  Mention objects contain a type key. In addition, mention objects contain a key corresponding with the value of type. The value is an object containing type-specific configuration. The type-specific configurations are described in the sections below.
	Equation *Equation `json:"equation" specific:"type"` // Equation objects contain the following information within the equation property:
}

func (p RichText) MarshalJSON() ([]byte, error) {
	type Alias RichText
	return marshalByType(Alias(p), p.Type)
}

// Style information which applies to the whole rich text object.
type Annotations struct {
	Bold          bool   `json:"bold"`          // Whether the text is bolded.
	Italic        bool   `json:"italic"`        // Whether the text is italicized.
	Strikethrough bool   `json:"strikethrough"` // Whether the text is struck through.
	Underline     bool   `json:"underline"`     // Whether the text is underlined.
	Code          bool   `json:"code"`          // Whether the text is code style.
	Color         string `json:"color"`         // Color of the text. Possible values are: "default", "gray", "brown", "orange", "yellow", "green", "blue", "purple", "pink", "red", "gray_background", "brown_background", "orange_background", "yellow_background", "green_background", "blue_background", "purple_background", "pink_background", "red_background".
}

// Text objects contain the following information within the text property:
type Text struct {
	Content string `json:"content"` // Text content. This field contains the actual content of your text and is probably the field you'll use most often.
	Link    *Link  `json:"link"`    // Any inline link in this text. See link objects.
}

// Text link objects contain a type key whose value is always "url" and a url key whose value is a web address.
type Link struct{}

/*
Mention objects represent an inline mention of a user, page, database, or date. In the app these are created by typing @ followed by the name of a user, page, database, or a date.

Mention objects contain a type key. In addition, mention objects contain a key corresponding with the value of type. The value is an object containing type-specific configuration. The type-specific configurations are described in the sections below.
*/
type Mention struct {
	Type string `json:"type"` // Type of the inline mention. Possible values include: "user", "page", "database", "date", "link_preview".
	// User mentions contain a user object within the user property.
	/*
	   Page mentions contain a page reference within the page property. A page reference is an object with an id property, with a string value (UUIDv4) corresponding to a page ID.

	   If an integration does not have access to the mentioned page, the mention will be returned with just the ID but without detailed information (title will appear as "Unititled" and annotations will be default).
	*/
	/*
	   Database mentions contain a database reference within the database property. A database reference is an object with an id property, with a string value (UUIDv4) corresponding to a database ID.

	   If an integration does not have access to the mentioned database, the mention will be returned with just the ID but without detailed information (title will appear as "Unititled" and annotations will be default).
	*/
	// Date mentions contain a date property value object within the date property.
	/*
	   Template mentions represent mentions within a template button or page that refer to a date or user upon duplication.

	   Template mentions contain a template_mention object with a nested type key that can be either "template_mention_date" or "template_mention_user".

	   If the type key is "template_mention_date", the template mention will contain the following property:
	   If the type key is "template_mention_user", the template mention will contain the following property:
	*/
	TemplateMentionDate string `json:"template_mention_date"` // Type of the date mention. Possible values include: "today" and "now".
	TemplateMentionUser string `json:"template_mention_user"` // Type of the user mention. The only possible value is "me".
	// Link preview mentions contain the originally pasted url.
	URL string `json:"url"` // The originally pasted url used to create the mention
}

// Equation objects contain the following information within the equation property:
type Equation struct {
	Expression string `json:"expression"` // The LaTeX string representing this inline equation.
}
