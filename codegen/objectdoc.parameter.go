package codegen

import (
	"bytes"
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

type PropertyOption struct {
	Nullable     bool
	OmitEmpty    bool
	TypeSpecific bool
}

// Property は *このドキュメントの情報から当然読み取れる範囲で* Propertyへの変換を試みます
// ドキュメントの記述と実際のAPIの挙動が一致せず、正しい変換にさらなる知識を要する場合、
// この関数を変更するのではなく呼び出し側で例外処理を行ってください
// TODO: TypeSpecificのセット
func (p ObjectDocParameter) Property(opt *PropertyOption) (*Property, error) {
	prop := &Property{
		Name:        strings.TrimSuffix(p.Name, "*"),
		Description: p.Description,
	}

	switch prop.Name {
	case "object", "type":
		// "user" "property_item" "list" など、文字列リテラルの場合
		if strings.HasPrefix(p.Type, `"`) && strings.HasSuffix(p.Type, `"`) {
			p.Type = "string"
		}
	}

	const string_enum_language = "string (enum) \n \t    \nType of block. Possible values include: \n\"abap\", \"arduino\", \"bash\", \"basic\", \"c\", \"clojure\", \"coffeescript\", \"c++\", \"c#\", \"css\", \"dart\", \"diff\", \"docker\", \"elixir\", \"elm\", \"erlang\", \"flow\", \"fortran\", \"f#\", \"gherkin\", \"glsl\", \"go\", \"graphql\", \"groovy\", \"haskell\", \"html\", \"java\", \"javascript\", \"json\", \"julia\", \"kotlin\", \"latex\", \"less\", \"lisp\", \"livescript\", \"lua\", \"makefile\", \"markdown\", \"markup\", \"matlab\", \"mermaid\", \"nix\", \"objective-c\", \"ocaml\", \"pascal\", \"perl\", \"php\", \"plain text\", \"powershell\", \"prolog\", \"protobuf\", \"python\", \"r\", \"reason\", \"ruby\", \"rust\", \"sass\", \"scala\", \"scheme\", \"scss\", \"shell\", \"sql\", \"swift\", \"typescript\", \"vb.net\", \"verilog\", \"vhdl\", \"visual basic\", \"webassembly\", \"xml\", \"yaml\", and \"java/c/c++/c#\""

	// prop.Typeが構造体となる場合、原則として構造体ポインタとしてください
	switch p.Type {
	case "string", "string enum", "string (enum)", "non-empty string", string_enum_language:
		prop.Type = jen.String()
	case "string (optional)", "string (optional, enum)", "string (optional enum)", "string or null", "optional string": // APIの挙動でnullを確認 (User.avatar_url, RichText.href)
		prop.Type = jen.Op("*").String()
	case "string (UUID)", "string (UUIDv4)":
		prop.Type = jen.Id("UUIDString")
	case "string (ISO 8601 date time)", "string (ISO 8601 date and time)", "string (ISO 8601 date)":
		prop.Type = jen.Id("ISO8601String")
	case "string (optional, ISO 8601 date and time)":
		prop.Type = jen.Op("*").Id("ISO8601String")
	case "number":
		prop.Type = jen.Float64()
	case "optional number":
		prop.Type = jen.Op("*").Float64()
	case "integer":
		prop.Type = jen.Int64()
	case "boolean", "boolean (optional)", "boolean (only true)":
		prop.Type = jen.Bool()
	case "array of string (UUID)":
		prop.Type = jen.Index().Id("UUIDString")
	case "array of rich text objects", "array of Rich text object text objects":
		prop.Type = jen.Index().Id("RichText")
	case "array of array of Rich text objects":
		prop.Type = jen.Index().Index().Id("RichText")
	case "array of block objects", "array of table_row block objects":
		prop.Type = jen.Index().Id("Block")
	case "array of file references":
		prop.Type = jen.Index().Id("File")
	case "array of user objects":
		prop.Type = jen.Index().Id("User")
	case "array of page references":
		prop.Type = jen.Index().Id("PageReference")
	case "date property value", "optional date property value":
		prop.Type = jen.Op("*").Id("DateValue")
	case "array of select option objects.", "array of multi-select option objects.", "array of multi-select option values":
		prop.Type = jen.Index().Id("SelectOption")
	case "array of status option objects.":
		prop.Type = jen.Index().Id("StatusOption")
	case "array of status group objects.":
		prop.Type = jen.Index().Id("StatusGroup")
	case "user object":
		prop.Type = jen.Op("*").Id("User")
	case "Partial User":
		prop.Type = jen.Op("*").Id("PartialUser")
	case "File Object", `File object (only type of "external" is supported currently)`:
		prop.Type = jen.Op("*").Id("File")
	case "File Object or Emoji object", `File Object (only type of "external" is supported currently) or Emoji object`:
		prop.Type = jen.Id("FileOrEmoji")
		prop.Union = &UnionInfo{
			InterfaceName: "FileOrEmoji",
			TypeProp:      "type",
			Map:           map[string]string{"file": "File", "emoji": "Emoji"},
		}
	case "Synced From Object":
		prop.Type = jen.Op("*").Id("SyncedFrom")
	case "object (number filter condition)", "object (date filter condition)", "object (text filter condition)", "object (checkbox filter condition)":
		name := strings.TrimSuffix(strings.TrimPrefix(p.Type, "object ("), ")")
		name = getName(strings.ReplaceAll(name, " ", "_"))
		prop.Type = jen.Op("*").Id(name)
	case "object (empty)":
		prop.Type = jen.Op("*").Struct()
	case "object", "object (optional)":
		switch p.Name {
		case "any", "every", "none":
			prop.Type = jen.Interface()
		case "parent", "user", "annotations", "link", "property_item":
			prop.Type = jen.Op("*").Id(strings.Title(p.Name))
		default:
			return nil, fmt.Errorf("unknown name for object: %v", p.Name)
		}
	case "list":
		switch p.Description {
		case "List of property_item objects.":
			prop.Type = jen.Index().Id("PropertyItem")
		default:
			return nil, fmt.Errorf("unknown description for list: %v", p.Description)
		}
	default:
		return nil, fmt.Errorf("unknown type: %v (name=%v)", p.Type, p.Name)
	}

	if opt == nil {
		opt = &PropertyOption{}
	}

	prop.OmitEmpty = opt.OmitEmpty
	prop.TypeSpecific = opt.TypeSpecific

	if opt.Nullable {
		buffer := bytes.NewBuffer(nil)
		if err := jen.Var().Id("_").Add(prop.Type).Render(buffer); err != nil {
			return nil, err
		}

		code := string(buffer.Bytes())
		if code == "var _ interface{}" {
			// インターフェイスなのでNullable
		} else if strings.HasPrefix(code, "var _ *") {
			// 既にNullable
		} else if strings.HasPrefix(code, "var _ []") {
			// スライスはNullableにしない
		} else {
			prop.Type = jen.Op("*").Add(prop.Type) // それ以外をNullableにする
		}
	}

	return prop, nil
}
