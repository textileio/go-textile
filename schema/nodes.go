package schema

var Photo = &Node{
	Pin: true,
	Nodes: map[string]*Node{
		"raw": {
			Use:  ":original",
			Mill: "/blob",
		},
		"exif": {
			Use:  ":original",
			Mill: "/image/exif",
		},
		"large": {
			Use:  ":original",
			Mill: "/image/resize",
			Opts: map[string]interface{}{
				"width":   1600,
				"quality": 80,
			},
		},
		"medium": {
			Use:  ":original",
			Mill: "/image/resize",
			Opts: map[string]interface{}{
				"width":   800,
				"quality": 80,
			},
		},
		"small": {
			Use:  ":original",
			Mill: "/image/resize",
			Opts: map[string]interface{}{
				"width":   320,
				"quality": 80,
			},
		},
		"thumb": {
			Use:  ":original",
			Mill: "/image/resize",
			Opts: map[string]interface{}{
				"width":   100,
				"quality": 80,
			},
			Pin: true,
		},
	},
}

/*

Example JSON mill usage:

var Todo = &Node{
	Use: ":original",
	Mill: "/json",
	Schema: jsonschema.Reflect(&TodoSchema{}),
	Pin: true,
}

type TodoSchema struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Complete bool `json:"complete"`
}

*/
