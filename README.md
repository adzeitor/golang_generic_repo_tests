# This is WORK IN PROGRESS repository. Be careful to use it in production code. I think better way is to copy/paste this code to your project.


## Generic tests

Test your storage/repository layer with
generic tests similar to quick check with "one line".


## How it works

This tests use (Property testing)[https://en.wikipedia.org/wiki/Property_testing] approach similar to quick check tests.

Usually repository have several generic properties likes this:
- After create record you could Load record with exact fields of reate
- After update you could Load record with next fields

So this kind of tests is here.

## How to add

- We need adapter of your repository interface to generic repository of this package.
- Also your Model/Aggregate should have property ID, because Load method needs it.

So in your test file:

```go
// ...
import genericrepotests "github.com/adzeitor/golang_generic_repo_tests"

type RepoAdapter struct {
	Repo MyRepo
}

func (r *RepoAdapter) Create(ctx context.Context, t *Model) error {
	return r.Repo.Save(ctx, t)
}

func (r *RepoAdapter) Update(ctx context.Context, t *Model) error {
	return r.Repo.Save(ctx, t)
}

func (r *RepoAdapter) Load(ctx context.Context, id uuid.UUID) (Model, error){
	return r.Repo.GetByID(ctx, id)
}


func TestGeneric(t *testing.T) {
	myRepo := MyRepo.New()
	genericrepotests.GenericTest(RepoAdapter{Repo: myRepo}, t)
}
```

## In case of matched interface (rare case)

```go
// ...
import "github.com/adzeitor/golang_generic_repo_tests"

/...

func TestGeneric(t *testing.T) {
        // this repo should have Create/Update/Load methods
	myRepo := MyRepo.New()
	GenericTest(myRepo, t)
}
```

## Copy/Paste

Another way just copy/paste this code to your repo and change method names to your convention.