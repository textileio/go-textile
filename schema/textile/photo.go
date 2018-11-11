package textile

var Photo = `
{
  "pin": true,
  "nodes": {
    "raw": {
      "use": ":original",
      "mill": "/blob"
    },
    "large": {
      "use": ":original",
      "mill": "/image/resize",
      "opts": {
        "width": 1600,
        "quality": 80
      }
    },
    "medium": {
      "use": ":original",
      "mill": "/image/resize",
      "opts": {
        "width": 800,
        "quality": 80
      }
    },
    "small": {
      "use": ":original",
      "mill": "/image/resize",
      "opts": {
        "width": 320,
        "quality": 80
      }
    },
    "thumb": {
      "pin": true,
      "use": ":original",
      "mill": "/image/resize",
      "opts": {
        "width": 100,
        "quality": 80
      }
    }
  }
}
`

// Example schema w/ using the JSON mill:
/*

var Person = `
{
  "pin": true,
  "use": ":original",
  "mill": "/json",
  "schema": {
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
    }
  }
}
`
*/
