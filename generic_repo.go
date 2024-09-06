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

type Config[T any] struct {
	Compare  func(t *testing.T, expected T, got T)
	MakeData func(t *testing.T) T
	FieldID  string
}

func DefaultConfig[T any]() Config[T] {
	return Config[T]{
		Compare: func(t *testing.T, expected T, got T) {
			assert.Equal(t, expected, got)
		},
		MakeData: func(t *testing.T) T {
			return testdata.Make[T](t)
		},
		FieldID: "ID",
	}
}

func GenericTest[T any, ID any](repo Repo[T, ID], t *testing.T) {
	GenericTestWithConfig(repo, DefaultConfig[T](), t)
}

func GenericTestWithConfig[T any, ID any](
	repo Repo[T, ID],
	newCfg Config[T],
	t *testing.T,
) {
	cfg := DefaultConfig[T]()
	if newCfg.Compare != nil {
		cfg.Compare = newCfg.Compare
	}
	if newCfg.MakeData != nil {
		cfg.MakeData = newCfg.MakeData
	}
	if newCfg.FieldID != "" {
		cfg.FieldID = newCfg.FieldID
	}

	getID := func(aggregate T) ID {
		field := reflect.
			ValueOf(aggregate).
			FieldByName(cfg.FieldID)
		if field == (reflect.Value{}) {
			t.Fatal("For generic tests your model/aggregate should have ID property, check config FieldID property")
		}
		return field.Interface().(ID)
	}
	setID := func(aggregate *T, id ID) {
		reflect.ValueOf(aggregate).
			Elem().
			FieldByName(cfg.FieldID).
			Set(reflect.ValueOf(id))
	}

	t.Run("Generic!", func(t *testing.T) {
		t.Run("Create/Load", func(t *testing.T) {
			// arrange
			ctx := context.Background()
			aggregate := cfg.MakeData(t)
			err := repo.Create(ctx, &aggregate)
			assert.NoError(t, err)

			got, err := repo.Load(ctx, getID(aggregate))
			assert.NoError(t, err)

			// assert
			cfg.Compare(t, *got, aggregate)
		})

		t.Run("Create/Update/Load", func(t *testing.T) {
			// arrange
			ctx := context.Background()
			aggregate := cfg.MakeData(t)
			err := repo.Create(ctx, &aggregate)
			assert.NoError(t, err)

			// act
			updated := cfg.MakeData(t)
			id := getID(aggregate)
			setID(&updated, id)
			err = repo.Update(ctx, &updated)
			assert.NoError(t, err)

			// assert
			got, err := repo.Load(ctx, id)
			assert.NoError(t, err)
			cfg.Compare(t, *got, updated)
		})
	})
}
