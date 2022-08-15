package codegen

type SSRProps struct {
	BaseURL                  string                 `json:"baseUrl"`
	APIBaseURL               string                 `json:"apiBaseUrl"`
	Config                   interface{}            `json:"config"`
	Context                  map[string]interface{} `json:"context"`                  // map[project:map[appearance:map[body:map[style:n...
	Doc                      SSRPropsDoc            `json:"doc"`                      // map[__v:57 _id:609176570b6bf20019821ce5 api:map...
	GlossaryTerms            []interface{}          `json:"glossaryTerms"`            // [map[_id:608fece52e9bb4000fea89a2 definition:TO...
	TideTOC                  bool                   `json:"hideTOC"`                  // false
	IsDetachedProductionSite bool                   `json:"isDetachedProductionSite"` // false
	Lang                     string                 `json:"lang"`                     // en
	LangFull                 string                 `json:"langFull"`                 // English
	LoginURL                 string                 `json:"loginUrl"`                 // https://dash.readme.com/to/notion-group
	Meta                     map[string]interface{} `json:"meta"`                     // map[hidden:false title:Query a database type:re...
	OasDefinition            map[string]interface{} `json:"oasDefinition"`            // map[_id:606ecc2cd9e93b0044cf6e47:609176570b6bf2...
	OasPublicURL             interface{}            `json:"oasPublicUrl"`             // @notionapi/v1#4rkc1lkz7dw9us
	Oauth                    bool                   `json:"oauth"`                    // false
	RdmdOpts                 map[string]interface{} `json:"rdmdOpts"`                 // map[compatibilityMode:false correctnewlines:fal...
	ReqURL                   string                 `json:"reqUrl"`                   // /reference/post-database-query
	Search                   map[string]interface{} `json:"search"`                   // map[UrlManager:map[defaults:map[child:<nil> lan...
	Sidebars                 interface{}            `json:"sidebars"`                 // map[docs:[map[__v:0 _id:6038057d9c4b200067ba3ca...
	SuggestedEdits           bool                   `json:"suggestedEdits"`           // true
	Variables                map[string]interface{} `json:"variables"`                // map[defaults:[map[_id:60c02cf4f26fa80064c30a79 ...
	Version                  map[string]interface{} `json:"version"`                  // map[__v:2 _id:6038057d9c4b200067ba3c9f categori...
}

type SSRPropsDoc struct {
	V                     float64                `json:"__v"`                   // 57
	ID                    string                 `json:"_id"`                   // 609176570b6bf20019821ce5
	API                   SSRPropsDocAPI         `json:"api"`                   // map[apiSetting:606ecc2cd9e93b0044cf6e47 auth:re...
	Body                  string                 `json:"body"`                  // Gets a list of [Pages](ref:page) contained in t...
	Category              string                 `json:"category"`              // 6091386ce2ca9200479fb438
	Children              []interface{}          `json:"children"`              // [map[__v:0 _id:6098885974ae4300418f9a18 api:map...
	ChildrenPages         []interface{}          `json:"childrenPages"`         // [map[__v:0 _id:6098885974ae4300418f9a18 api:map...
	CreatedAt             string                 `json:"createdAt"`             // 2021-05-04T16:29:11.027Z
	Deprecated            bool                   `json:"deprecated"`            // false
	Excerpt               string                 `json:"excerpt"`               //
	Hidden                bool                   `json:"hidden"`                // false
	Icon                  interface{}            `json:"icon,omitempty"`        //
	IsReference           bool                   `json:"isReference"`           // true
	LinkExternal          bool                   `json:"link_external"`         // false
	LinkURL               string                 `json:"link_url"`              //
	Metadata              map[string]interface{} `json:"metadata"`              // map[description: image:[] title:]
	Next                  map[string]interface{} `json:"next"`                  // map[description: pages:[]]
	Order                 float64                `json:"order"`                 // 1
	ParentDoc             interface{}            `json:"parentDoc"`             // <nil>
	PendingAlgoliaPublish bool                   `json:"pendingAlgoliaPublish"` // false
	PreviousSlug          string                 `json:"previousSlug"`          // post-databases-query
	Project               string                 `json:"project"`               // 6038057d9c4b200067ba3c9a
	Slug                  string                 `json:"slug"`                  // post-database-query
	SlugUpdatedAt         string                 `json:"slugUpdatedAt"`         // 2021-05-10T00:46:29.470Z
	Swagger               map[string]interface{} `json:"swagger,omitempty"`     // map[path:/v1/databases/{database_id}/query]
	SyncUnique            string                 `json:"sync_unique"`           //
	Title                 string                 `json:"title"`                 // Query a database
	Type                  string                 `json:"type"`                  // endpoint
	UpdatedAt             string                 `json:"updatedAt"`             // 2021-12-23T16:56:23.254Z
	Updates               []interface{}          `json:"updates"`               // []
	User                  string                 `json:"user"`                  // 60917de732252800631fcd43
	Version               string                 `json:"version"`               // 6038057d9c4b200067ba3c9f
}

type SSRPropsDocAPI struct {
	APISetting string                 `json:"apiSetting,omitempty"` // 606ecc2cd9e93b0044cf6e47
	Auth       string                 `json:"auth"`                 // required
	Examples   map[string]interface{} `json:"examples,omitempty"`   // map[codes:[map[code:const { Client } = require(...
	Method     string                 `json:"method"`               // post
	Params     []SSRPropsDocAPIParam  `json:"params"`               // [map[_id:609176570b6bf20019821ce8 default: desc...
	Results    map[string]interface{} `json:"results"`              // map[codes:[map[code:{"object": "list","resu...
	URL        string                 `json:"url"`                  // /v1/databases/{database_id}/query
}

type SSRPropsDocAPIParam struct {
	ID         string `json:"_id"`        // 609176570b6bf20019821ce8
	Default    string `json:"default"`    //
	Desc       string `json:"desc"`       // When supplied, limits which pages are returned ...
	EnumValues string `json:"enumValues"` //
	In         string `json:"in"`         // body
	Name       string `json:"name"`       // filter
	Ref        string `json:"ref"`        //
	Required   bool   `json:"required"`   // false
	Type       string `json:"type"`       // json
}
