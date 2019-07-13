package mill

import (
	"fmt"
	"testing"
)

func TestRoles_Mill(t *testing.T) {
	m := &Roles{}

	sharedBlog := `
{
  "default": "ANNOTATE",
  "accounts": {
    "P7BRCWyWkL4NZvxPx4mZFiYjwyq7aGwLqhzstqGSBkYe5sSH": "WRITE",
    "P6a4L9A4QGrm9tagB7cu7sBSYKHd73cua5jjRwRJictZdS6m": "WRITE",
    "P6YU9UzpG5rKUnDbsYozt1ebSaVjzEZJE1V7xubWnUz9k1uK": "NO_ACCESS"
  }
}
`

	res, err := m.Mill([]byte(sharedBlog), "test")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(res.File))
}
