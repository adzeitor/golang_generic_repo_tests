package genericrepotests

import (
	"context"
	"errors"
	"math/rand/v2"
	"strings"
	"testing"

	"github.com/kyuff/testdata"
)

type ResourceSlot struct {
	Resource string
	Amount   int
}

type CargoWagon struct {
	Resources []ResourceSlot
}

// Just demonstration of repository.
// We have no locks here, etc.
type Train struct {
	ID          int
	Name        string
	CargoWagons []CargoWagon
}

type TrainsRepo struct {
	NextID int
	Trains []Train
}

func (r *TrainsRepo) nextID() int {
	id := r.NextID
	r.NextID += 1
	return id
}

func (r *TrainsRepo) Create(_ context.Context, d *Train) error {
	d.ID = r.nextID()
	r.Trains = append(r.Trains, *d)
	return nil
}

func (r *TrainsRepo) Update(_ context.Context, d *Train) error {
	r.Trains[d.ID] = *d
	return nil
}

func (r *TrainsRepo) Load(_ context.Context, id int) (*Train, error) {
	return &r.Trains[id], nil
}

func TestGeneric(t *testing.T) {
	trains := &TrainsRepo{}
	GenericTest(trains, t)
}

type RepoAdapter struct {
	Repo TrainsRepo
}

func (r *RepoAdapter) Create(ctx context.Context, t *Train) error {
	return r.Repo.Create(ctx, t)
}

func (r *RepoAdapter) Update(ctx context.Context, t *Train) error {
	return r.Repo.Update(ctx, t)
}

func (r *RepoAdapter) Load(ctx context.Context, id int) (*Train, error) {
	return r.Repo.Load(ctx, id)
}

// TODO: add generic adapter with function
func TestGenericWithAdapter(t *testing.T) {
	withAdapter := &RepoAdapter{Repo: TrainsRepo{}}
	GenericTest(withAdapter, t)
}

type Item struct {
	ID   string
	Name string
}
type MyRepo struct {
	items map[string]Item
}

func (r *MyRepo) Create(_ context.Context, item *Item) error {
	if !strings.HasPrefix(item.Name, "+") {
		return errors.New("name should start with +")
	}
	r.items[item.ID] = *item
	return nil
}

func (r *MyRepo) Update(_ context.Context, item *Item) error {
	r.items[item.ID] = *item
	return nil
}

func (r *MyRepo) Load(_ context.Context, id string) (*Item, error) {
	item := r.items[id]
	return &item, nil
}

func TestGenericWithConfig(t *testing.T) {
	repo := &MyRepo{
		items: map[string]Item{},
	}
	cfg := testdata.NewConfig(
		testdata.WithGenerator(func(rand *rand.Rand) Item {
			return Item{
				ID:   testdata.Make[string](t),
				Name: "+" + testdata.Make[string](t),
			}
		}),
	)
	GenericTestWithConfig(repo, cfg, t)
}
