package testdata

type TestImage struct {
	Path    string
	Format  string
	HasExif bool
	Width   int
	Height  int
}

var Images = []TestImage{
	{
		Path:    "testdata/image.jpg",
		Format:  "jpeg",
		HasExif: true,
		Width:   3024,
		Height:  4032,
	},
	{
		Path:    "testdata/image.png",
		Format:  "png",
		HasExif: false,
		Width:   3024,
		Height:  4032,
	},
	{
		Path:    "testdata/image.gif",
		Format:  "gif",
		HasExif: false,
		Width:   320,
		Height:  240,
	},
}
