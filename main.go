package main

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"os"
	"strings"
	"text/template"
)

func main() {
	fMap := registerFunc()
	t, err := template.New("types").Funcs(fMap).ParseGlob("templates/*.gots")
	if err != nil {
		panic(err)
	}
	t = t.Funcs(fMap)
	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("test.yaml")
	if err != nil {
		panic(err)
	}
	file, err := os.OpenFile("schema.ts", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = t.ExecuteTemplate(file, "schema", swagger)

}

func registerFunc() template.FuncMap {
	fMap := make(template.FuncMap, 0)
	fMap["getType"] = getType
	fMap["scalarType"] = scalarType
	fMap["title"] = strings.Title
	fMap["reqProp"] = required
	fMap["getEnums"] = getEnums
	fMap["classFromRef"] = classFromRef
	fMap["nli"] = notLastIndex
	fMap["isEnum"] = isEnum
	fMap["schemas"] = schemas
	return fMap
}

func schemas(scm map[string]*openapi3.SchemaRef) (ret []*openapi3.SchemaRef) {
	if scm != nil {
		for name, ref := range scm {
			ref.Ref = "#/components/schemas/" + name
			ret = append(ret, ref)
		}
	}
	return
}

func isEnum(scr *openapi3.SchemaRef) bool {
	return scr.Value.Type == "string" && len(scr.Value.Enum) > 0
}

func notLastIndex(i, l int) bool {
	return (i + 1) != l
}

func classFromRef(ref string) string {
	return ref[21:]
}

func required(sc *openapi3.Schema, name string) string {
	for _, s := range sc.Required {
		if s == name {
			return ""
		}
	}
	return "?"
}

func scalarType(scr *openapi3.SchemaRef) string {
	sc := scr.Value
	switch sc.Type {
	case "array":
		return fmt.Sprintf("Array<%s>", getType(sc.Items))
	case "integer":
		return "number"
	default:
		return sc.Type
	}
}

func getType(scr *openapi3.SchemaRef) string {
	if scr.Ref != "" {
		return scr.Ref[21:]
	} else {
		sc := scr.Value
		switch sc.Type {
		case "array":
			return fmt.Sprintf("Array<%s>", getType(sc.Items))
		case "integer":
			return "number"
		default:
			return sc.Type
		}
	}
}

func getEnums(sc *openapi3.SchemaRef) []EnumData {
	ret := make([]EnumData, 0)
	for name, scf := range sc.Value.Properties {
		if scf.Ref == "" && scf.Value.Type == "string" && len(scf.Value.Enum) > 0 {
			enm := EnumData{
				ClassName: sc.Ref,
				EnumName:  name,
				Values:    scf.Value.Enum,
			}
			ret = append(ret, enm)
		}
	}
	return ret
}

type EnumData struct {
	ClassName string
	EnumName  string
	Values    []interface{}
}
