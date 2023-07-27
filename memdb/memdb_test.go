package memdb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/flinnb/memdb/memdb"
)

func TestMain(m *testing.M) {
}

func TestMemDB(t *testing.T) {
	t.Run("Case 1", func(t *testing.T) {
		memdb.Init()

		v := memdb.Get("a")
		assert.Equal(t, "NULL", v)

		memdb.Set("a", "foo")
		memdb.Set("B", "foo")
		c := memdb.Count("foo")
		assert.Equal(t, 2, c)

		c = memdb.Count("bar")
		assert.Equal(t, 0, c)

		memdb.Delete("a")
		c = memdb.Count("foo")
		assert.Equal(t, 1, c)

		memdb.Set("b", "baz")
		c = memdb.Count("foo")
		assert.Equal(t, 0, c)

		v = memdb.Get("b")
		assert.Equal(t, "baz", v)

		v = memdb.Get("B")
		assert.Equal(t, "NULL", v)
	})

	t.Run("Case 2", func(t *testing.T) {
		memdb.Init()

		memdb.Set("a", "foo")
		memdb.Set("a", "foo")
		c := memdb.Count("foo")
		assert.Equal(t, 1, c)

		v := memdb.Get("a")
		assert.Equal(t, "foo", v)

		memdb.Delete("a")
		v = memdb.Get("a")
		assert.Equal(t, "NULL", v)

		c = memdb.Count("foo")
		assert.Equal(t, 0, c)
	})

	t.Run("Case 3", func(t *testing.T) {
		memdb.Init()

		memdb.Begin()
		memdb.Set("a", "foo")
		v := memdb.Get("a")
		assert.Equal(t, "foo", v)

		memdb.Begin()
		memdb.Set("a", "bar")
		v = memdb.Get("a")
		assert.Equal(t, "bar", v)

		memdb.Set("a", "baz")
		memdb.Rollback()
		v = memdb.Get("a")
		assert.Equal(t, "foo", v)

		memdb.Rollback()
		v = memdb.Get("a")
		assert.Equal(t, "NULL", v)
	})

	t.Run("Case 4", func(t *testing.T) {
		memdb.Init()

		memdb.Set("a", "foo")
		memdb.Set("b", "baz")
		memdb.Begin()
		v := memdb.Get("a")
		assert.Equal(t, "foo", v)

		memdb.Set("a", "bar")
		c := memdb.Count("bar")
		assert.Equal(t, 1, c)

		memdb.Begin()
		c = memdb.Count("bar")
		assert.Equal(t, 1, c)

		memdb.Delete("a")
		v = memdb.Get("a")
		assert.Equal(t, "NULL", v)
		c = memdb.Count("bar")
		assert.Equal(t, 0, c)

		memdb.Rollback()
		v = memdb.Get("a")
		assert.Equal(t, "bar", v)
		c = memdb.Count("bar")
		assert.Equal(t, 1, c)

		memdb.Commit()
		v = memdb.Get("a")
		assert.Equal(t, "bar", v)

		v = memdb.Get("b")
		assert.Equal(t, "baz", v)
	})
}
