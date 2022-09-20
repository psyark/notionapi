package notionapi

// Code generated by notion.codegen; DO NOT EDIT.

// https://developers.notion.com/reference/user

/*
User objects appear in the API in nearly all objects returned by the API, including:

Block object under created_by and last_edited_by.
Page object under created_by and last_edited_by and in people property items.
Database object under created_by and last_edited_by.
Rich text object, as user mentions.
Property object when the property is a people property.

User objects will always contain object and id keys, as described below. The remaining properties may appear if the user is being rendered in a rich text or page property context, and the bot has the correct capabilities to access those properties. For more about capabilities, see the Capabilities guide and the Authorization guide.
*/
type User struct {
	// These fields are shared by all users, including people and bots. Fields marked with * are always present.
	Object    string     `json:"object"`               // Always "user"
	ID        UUIDString `json:"id"`                   // Unique identifier for this user.
	Type      string     `json:"type,omitempty"`       // Type of the user. Possible values are "person" and "bot".
	Name      *string    `json:"name,omitempty"`       // User's name, as displayed in Notion.
	AvatarURL *string    `json:"avatar_url,omitempty"` // Chosen avatar image.
	// Properties only present for non-bot users.
	Person *Person `json:"person,omitempty"`
	// Properties only present for bot users.If viewing your own bot with GET /v1/users/me or GET /v1/users/{{your_bot_id}}, this field will be populated with more information about the bot.
	Bot *Bot `json:"bot,omitempty"`
}

// User objects that represent people have the type property set to "person". These objects also have the following properties:
type Person struct {
	Email string `json:"email"` // Email address of person. This is only present if an integration has user capabilities that allow access to email addresses.
}

// User objects that represent bots have the type property set to "bot". These objects also have the following properties:
type Bot struct {
	Owner *Owner `json:"owner,omitempty"` // Information about who owns this bot.
}

// User objects that represent bots have the type property set to "bot". These objects also have the following properties:
type Owner struct {
	Type      string `json:"type"`      // The type of owner - either "workspace" or "user".
	Workspace bool   `json:"workspace"` // Always true. Only present if owner.type is "workspace".
	User      *User  `json:"user"`      // A user object describing who owns this bot. Currently only "person" users can own bots. See the People reference above for more detail. The properties in the user object are based on the bot capabilities.
}

type UserPropertyValueData struct {
	// These fields are shared by all users, including people and bots. Fields marked with * are always present.
	Object    string     `json:"object"`     // Always "user"
	ID        UUIDString `json:"id"`         // Unique identifier for this user.
	Type      string     `json:"type"`       // Type of the user. Possible values are "person" and "bot".
	Name      *string    `json:"name"`       // User's name, as displayed in Notion.
	AvatarURL *string    `json:"avatar_url"` // Chosen avatar image.
	// Properties only present for non-bot users.
	Person *Person `json:"person" specific:"type"`
	// Properties only present for bot users.If viewing your own bot with GET /v1/users/me or GET /v1/users/{{your_bot_id}}, this field will be populated with more information about the bot.
	Bot *Bot `json:"bot" specific:"type"`
}

func (p UserPropertyValueData) MarshalJSON() ([]byte, error) {
	type Alias UserPropertyValueData
	return marshalByType(Alias(p), p.Type)
}
