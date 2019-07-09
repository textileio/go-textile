package mill

import (
	"testing"
)

func TestRoles_Mill(t *testing.T) {
	m := &Roles{}

	sharedBlog := `
{
  "default": 2,
  "accounts": {
    "P7BRCWyWkL4NZvxPx4mZFiYjwyq7aGwLqhzstqGSBkYe5sSH": 3,
    "P6a4L9A4QGrm9tagB7cu7sBSYKHd73cua5jjRwRJictZdS6m": 3,
    "P6YU9UzpG5rKUnDbsYozt1ebSaVjzEZJE1V7xubWnUz9k1uK": 0
  }
}
`

	if _, err := m.Mill([]byte(sharedBlog), "test"); err != nil {
		t.Fatal(err)
	}
}
