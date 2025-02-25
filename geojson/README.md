# geo/geojson [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/geo)](https://pkg.go.dev/github.com/pchchv/geo/geojson)

This package **encodes** and **decodes** [GeoJSON](http://geojson.org/) in Go structs using the geometries in the [geo](https://github.com/pchchv/geo) package.

Supports both the [json.Marshaler](https://pkg.go.dev/encoding/json#Marshaler) and [json.Unmarshaler](https://pkg.go.dev/encoding/json#Unmarshaler) interfaces. Package also provides helper functions such as `UnmarshalFeatureCollection` and `UnmarshalFeature`.

Types also support *BSON* using [bson.Marshaler](https://pkg.go.dev/go.mongodb.org/mongo-driver/bson#Marshaler) and [bson.Unmarshaler](https://pkg.go.dev/go.mongodb.org/mongo-driver/bson#Unmarshaler) interfaces.
These types can be used directly when working with MongoDB.

### Unmarshalling (JSON -> Go)

```go
rawJSON := []byte(`
  { "type": "FeatureCollection",
    "features": [
      { "type": "Feature",
        "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
        "properties": {"prop0": "value0"}
      }
    ]
  }`)

fc, _ := geojson.UnmarshalFeatureCollection(rawJSON)
// or
fc := geojson.NewFeatureCollection()
err := json.Unmarshal(rawJSON, &fc)
// Geometry will be unmarshalled into the correct geo.Geometry type.
point := fc.Features[0].Geometry.(geo.Point)
```

## Marshalling (Go -> JSON)

```go
fc := geojson.NewFeatureCollection()
fc.Append(geojson.NewFeature(geo.Point{1, 2}))

rawJSON, _ := fc.MarshalJSON()
// or
blob, _ := json.Marshal(fc)
```

## Extra members in a feature collection

```go
rawJSON := []byte(`
  { "type": "FeatureCollection",
    "generator": "myapp",
    "timestamp": "2020-06-15T01:02:03Z",
    "features": [
      { "type": "Feature",
        "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
        "properties": {"prop0": "value0"}
      }
    ]
  }`)

fc, _ := geojson.UnmarshalFeatureCollection(rawJSON)
fc.ExtraMembers["generator"] // == "myApp"
fc.ExtraMembers["timestamp"] // == "2020-06-15T01:02:03Z"

// marshalling will include values in `ExtraMembers` in the base featureCollection object.
```

## Feature Properties

GeoJSON features can have properties of any type. This can cause issues in a statically typed language such as Go.
Included is a `Properties` type with some helper methods that will attempt to force convert a property.
An optional default value will be used if the property is missing or of the wrong type.

```go
f.Properties.MustBool(key string, def ...bool) bool
f.Properties.MustFloat64(key string, def ...float64) float64
f.Properties.MustInt(key string, def ...int) int
f.Properties.MustString(key string, def ...string) string
```