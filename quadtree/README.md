# geo/quadtree [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/geo)](https://pkg.go.dev/github.com/pchchv/geo/quadtree)

The *quadtree* package implements quadtree using rectangular partitions. Each point exists in a unique node.

## API

```go
func New(bound geo.Bound) *Quadtree
func (q *Quadtree) Bound() geo.Bound
func (q *Quadtree) Add(p geo.Pointer) error
func (q *Quadtree) Remove(p geo.Pointer, eq FilterFunc) bool
func (q *Quadtree) Find(p geo.Point) geo.Pointer
func (q *Quadtree) Matching(p geo.Point, f FilterFunc) geo.Pointer
func (q *Quadtree) KNearest(buf []geo.Pointer, p geo.Point, k int, maxDistance ...float64) []geo.Pointer
func (q *Quadtree) KNearestMatching(buf []geo.Pointer, p geo.Point, k int, f FilterFunc, maxDistance ...float64) []geo.Pointer
func (q *Quadtree) InBound(buf []geo.Pointer, b geo.Bound) []geo.Pointer
func (q *Quadtree) InBoundMatching(buf []geo.Pointer, b geo.Bound, f FilterFunc) []geo.Pointer
```

## Examples

```go
func ExampleQuadtree_Find() {
    r := rand.New(rand.NewSource(42)) // to make things reproducible
    qt := quadtree.New(geo.Bound{Min: geo.Point{0, 0}, Max: geo.Point{1, 1}})
    // add 1000 random points
    for i := 0; i < 1000; i++ {
        qt.Add(geo.Point{r.Float64(), r.Float64()})
    }

    nearest := qt.Find(geo.Point{0.5, 0.5})

    fmt.Printf("nearest: %+v\n", nearest)
    // Output:
    // nearest: [0.4930591659434973 0.5196585530161364]
}
```