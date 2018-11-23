package mill

import (
	"testing"
)

func TestJson_Mill(t *testing.T) {
	m := &Json{
		Opts: JsonOpts{
			Schema: `
{
  "$id": "https://example.com/person.schema.json",
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Person",
  "type": "object",
  "properties": {
    "firstName": {
      "type": "string",
      "description": "The person's first name."
    },
    "lastName": {
      "type": "string",
      "description": "The person's last name."
    },
    "age": {
      "description": "Age in years which must be equal to or greater than zero.",
      "type": "integer",
      "minimum": 0
    }
  },
  "required": [
    "age"
  ],
  "additionalProperties": false
}
`,
		},
	}

	doc1 := `
{"firstName": "Grigori", "lastName": "Rasputin"}
`
	doc2 := `
{"firstName": "Grigori", "lastName": "Rasputin", "age": "foo"}
`
	doc3 := `
{"firstName": "Grigori", "lastName": "Rasputin", "age": 47}
`
	doc4 := `
{"firstName": "Grigori", "lastName": "Rasputin", "age": 47, "whacko": true}
`

	if _, err := m.Mill([]byte(doc1), "test"); err == nil {
		t.Error("age should be required")
	}
	if _, err := m.Mill([]byte(doc2), "test"); err == nil {
		t.Error("age should be int type")
	}
	if _, err := m.Mill([]byte(doc3), "test"); err != nil {
		t.Error(err)
	}
	if _, err := m.Mill([]byte(doc4), "test"); err == nil {
		t.Error("should not allow extra props")
	}
}
