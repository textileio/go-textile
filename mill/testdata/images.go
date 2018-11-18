package testdata

type TestImage struct {
	Path    string
	Format  string
	HasExif bool
	Width   int
	Height  int
	Created int64
}

var Images = []TestImage{
	{
		Path:    "testdata/image.jpg",
		Format:  "jpeg",
		HasExif: true,
		Width:   3024,
		Height:  4032,
		Created: 1526926068,
	},
	{
		Path:    "testdata/image.png",
		Format:  "png",
		HasExif: false,
		Width:   3024,
		Height:  4032,
		Created: -62135596800,
	},
	{
		Path:    "testdata/image.gif",
		Format:  "gif",
		HasExif: false,
		Width:   320,
		Height:  240,
		Created: -62135596800,
	},
}
