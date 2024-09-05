package genericrepotests

import (
	"context"
	"reflect"
	"testing"

	"github.com/kyuff/testdata"
	"github.com/stretchr/testify/assert"
)

type Repo[T any, ID any] interface {
	Create(context.Context, *T) error
	Update(context.Context, *T) error
	Load(context.Context, ID) (*T, error)
}

func GenericTest[T any, ID any](repo Repo[T, ID], t *testing.T) {
	getID := func(aggregate T) ID {
		field := reflect.ValueOf(aggregate).FieldByName("ID")
		if field == (reflect.Value{}) {
			t.Fatal("For generic tests your model/aggregate should have ID property")
		}
		return field.Interface().(ID)
	}
	setID := func(aggregate *T, id ID) {
		reflect.ValueOf(aggregate).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
	}

	t.Run("Generic!", func(t *testing.T) {
		t.Run("Create/Load", func(t *testing.T) {
			// arrange
			ctx := context.Background()
			aggregate := testdata.Make[T](t)
			err := repo.Create(ctx, &aggregate)
			assert.NoError(t, err)

			got, err := repo.Load(ctx, getID(aggregate))
			assert.NoError(t, err)

			// assert
			assert.Equal(t, *got, aggregate)
		})

		t.Run("Create/Update/Load", func(t *testing.T) {
			// arrange
			ctx := context.Background()
			aggregate := testdata.Make[T](t)
			err := repo.Create(ctx, &aggregate)
			assert.NoError(t, err)

			// act
			updated := testdata.Make[T](t)
			id := getID(aggregate)
			setID(&updated, id)
			err = repo.Update(ctx, &updated)
			assert.NoError(t, err)

			// assert
			got, err := repo.Load(ctx, id)
			assert.NoError(t, err)
			assert.Equal(t, *got, updated)
		})
	})
}
