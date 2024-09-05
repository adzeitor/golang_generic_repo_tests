package genericrepotests

import (
	"context"
	"testing"
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
