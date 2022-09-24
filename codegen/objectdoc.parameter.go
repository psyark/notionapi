package codegen

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
)

type ObjectDocParameter struct {
	Name         string `json:"Property"`
	Type         string
	Description  string
	ExampleValue string `json:"Example value"`
}

// Property は *このドキュメントの情報から当然読み取れる範囲で* Propertyへの変換を試みます
// ドキュメントの記述と実際のAPIの挙動が一致せず、正しい変換にさらなる知識を要する場合、
// この関数を変更するのではなく呼び出し側で例外処理を行ってください
func (p *ObjectDocParameter) Property() (*Property, error) {
	prop := &Property{
		Name:        strings.TrimSuffix(p.Name, "*"),
		Description: p.Description,
	}

	const string_enum_language = "string (enum) \n \t    \nType of block. Possible values include: \n\"abap\", \"arduino\", \"bash\", \"basic\", \"c\", \"clojure\", \"coffeescript\", \"c++\", \"c#\", \"css\", \"dart\", \"diff\", \"docker\", \"elixir\", \"elm\", \"erlang\", \"flow\", \"fortran\", \"f#\", \"gherkin\", \"glsl\", \"go\", \"graphql\", \"groovy\", \"haskell\", \"html\", \"java\", \"javascript\", \"json\", \"julia\", \"kotlin\", \"latex\", \"less\", \"lisp\", \"livescript\", \"lua\", \"makefile\", \"markdown\", \"markup\", \"matlab\", \"mermaid\", \"nix\", \"objective-c\", \"ocaml\", \"pascal\", \"perl\", \"php\", \"plain text\", \"powershell\", \"prolog\", \"protobuf\", \"python\", \"r\", \"reason\", \"ruby\", \"rust\", \"sass\", \"scala\", \"scheme\", \"scss\", \"shell\", \"sql\", \"swift\", \"typescript\", \"vb.net\", \"verilog\", \"vhdl\", \"visual basic\", \"webassembly\", \"xml\", \"yaml\", and \"java/c/c++/c#\""

	// prop.Typeが構造体となる場合、原則として構造体ポインタとしてください
	switch p.Type {
	case "string", "string enum", "string (enum)", "string (optional, enum)", string_enum_language:
		prop.Type = jen.String()
	case "string (optional)": // APIの挙動でnullを確認 (User.avatar_url, RichText.href)
		prop.Type = jen.Op("*").String()
	case "string (UUID)", "string (UUIDv4)":
		prop.Type = jen.Id("UUIDString")
	case "string (ISO 8601 date time)", "string (ISO 8601 date and time)", "string (ISO 8601 date)":
		prop.Type = jen.Id("ISO8601String")
	case "number":
		prop.Type = jen.Float64()
	case "integer":
		prop.Type = jen.Int64()
	case "boolean", "boolean (optional)", "boolean (only true)":
		prop.Type = jen.Bool()
	case "array of rich text objects", "array of Rich text object text objects":
		prop.Type = jen.Index().Id("RichText")
	case "array of array of Rich text objects":
		prop.Type = jen.Index().Index().Id("RichText")
	case "array of block objects", "array of table_row block objects":
		prop.Type = jen.Index().Id("Block")
	case "Partial User":
		prop.Type = jen.Op("*").Id("PartialUser")
	case "File Object", `File object (only type of "external" is supported currently)`:
		prop.Type = jen.Op("*").Id("File")
	case `File Object (only type of "external" is supported currently) or Emoji object`:
		prop.Type = jen.Op("*").Id("FileOrEmoji")
	case "File Object or Emoji object":
		prop.Type = jen.Op("*").Id("FileOrEmoji")
	case "Synced From Object":
		prop.Type = jen.Op("*").Id("SyncedFrom")
	case "object (number filter condition)":
		prop.Type = jen.Op("*").Id("NumberFilterCondition")
	case "object (date filter condition)":
		prop.Type = jen.Op("*").Id("DateFilterCondition")
	case "object (text filter condition)":
		prop.Type = jen.Op("*").Id("TextFilterCondition")
	case "object (checkbox filter condition)":
		prop.Type = jen.Op("*").Id("CheckboxFilterCondition")
	case "object (empty)":
		prop.Type = jen.Op("*").Struct()
	case "object", "object (optional)":
		switch p.Name {
		case "any", "every", "none":
			prop.Type = jen.Interface()
		case "parent":
			prop.Type = jen.Op("*").Id("Parent")
		case "user":
			prop.Type = jen.Op("*").Id("User")
		case "annotations":
			prop.Type = jen.Op("*").Id("Annotations")
		case "link":
			prop.Type = jen.Op("*").Id("Link")
		default:
			return nil, fmt.Errorf("unknown name for object: %v", p.Name)
		}
	default:
		return nil, fmt.Errorf("unknown type: %v (name=%v)", p.Type, p.Name)
	}
	return prop, nil
}
