package gens

// const (
// 	fieldName  = `"field":`
// 	methodName = `"function":`
// )

// type JSONFormat struct {
// }

// type jsonFormat struct {
// 	Field    map[string]fieldType
// 	Function []methodType
// }

// func (JSONFormat) Extension() string {
// 	return ".json"
// }

// func (JSONFormat) Start() string {
// 	return "{"
// }

// func (JSONFormat) Separate() string {
// 	return ","
// }

// func (this *JSONFormat) Struct(str *Struct) string {
// 	sb := strings.Builder{}
// 	strtype := structType{
// 		Name:        str.Name,
// 		Description: str.Comment,
// 	}

// 	bytes, _ := json.Marshal(strtype)
// 	sb.WriteString("\"struct\":")
// 	sb.Write(bytes)
// 	return sb.String()
// }

// func (this *JSONFormat) Methods(methods []*Method) string {
// 	sb := strings.Builder{}
// 	ms := make(map[string]methodType)
// 	for _, m := range methods {
// 		jm := methodType{
// 			Name:        m.Name(),
// 			Description: m.Comment,
// 		}

// 		params := m.Params()
// 		jm.Params = make([]fieldType, 0, len(params))
// 		for _, p := range params {
// 			jm.Params = append(jm.Params, fieldType{
// 				Name: p.Name(),
// 				Type: p.Type().String(),
// 			})
// 		}

// 		if r := m.Return(); r != nil {
// 			jm.Return = fieldType{
// 				Name: r.Name(),
// 				Type: r.Type().String(),
// 			}
// 		}

// 		ms[jm.Name] = jm
// 		// ms = append(ms, jm)
// 	}

// 	bytes, _ := json.Marshal(ms)
// 	sb.WriteString(methodName)
// 	sb.Write(bytes)

// 	return sb.String()
// }

// func (this *JSONFormat) Fields(fields []*Field) string {
// 	var (
// 		sb       = strings.Builder{}
// 		varTypes = make(map[string]interface{}, len(fields))
// 	)

// 	for _, f := range fields {
// 		if str := f.Struct(); str != nil {
// 			g := NewInformationGenerator(str, this)
// 			_ = g.Generate(context.TODO())

// 			jsonFormat := jsonFormat{}
// 			_ = json.Unmarshal(g.Bytes(), &jsonFormat)

// 			varTypes[f.Name()] = jsonFormat
// 		} else {
// 			varTypes[f.Name()] = f.Type().String()
// 		}

// 	}

// 	bytes, _ := json.Marshal(varTypes)
// 	sb.WriteString(fieldName)
// 	sb.Write(bytes)

// 	return sb.String()
// }

// func (this *JSONFormat) End() string {
// 	return "}"
// }

// func NewJSONFormat() *JSONFormat {
// 	return &JSONFormat{}
// }
