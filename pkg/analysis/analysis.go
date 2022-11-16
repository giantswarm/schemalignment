// Package analysis analyses schemas and creates comparison.
package analysis

import (
	"sort"

	"github.com/giantswarm/microerror"
	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
)

type Analyser struct {
	SchemaUrls      map[string]string
	Schemas         map[string]*jsonschema.Schema
	FlattenedSchema map[string]map[string]*jsonschema.Schema
}

type PropertyDetails struct {
	Key             string
	KeyHierarchical string
	Types           []string
	DefaultValue    string
}

func New(schemaUrls map[string]string) (*Analyser, error) {
	a := &Analyser{
		SchemaUrls: schemaUrls,
		Schemas:    make(map[string]*jsonschema.Schema),
	}

	for provider, url := range schemaUrls {
		schema, err := jsonschema.Compile(url)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		a.Schemas[provider] = schema
	}

	a.FlattenedSchema = make(map[string]map[string]*jsonschema.Schema)
	for provider := range schemaUrls {
		a.FlattenedSchema[provider] = make(map[string]*jsonschema.Schema)
		a.FlattenedSchema[provider] = a.flattened(provider)
	}

	return a, nil
}

func (a *Analyser) FullSchemas() map[string][]string {
	// Create complete list of all property keys with a list of the providers having them.
	fullSchemas := make(map[string][]string)
	for provider := range a.SchemaUrls {
		// collect all keys
		for key := range a.FlattenedSchema[provider] {
			_, ok := fullSchemas[key]
			if ok {
				fullSchemas[key] = append(fullSchemas[key], provider)
			} else {
				fullSchemas[key] = []string{provider}
			}
		}
	}

	return fullSchemas
}

// Flattens the schema, using a hierarchical key.
func (a *Analyser) flattened(provider string) map[string]*jsonschema.Schema {
	var mymap = make(map[string]*jsonschema.Schema)
	return flattened(mymap, a.Schemas[provider], "", 0)
}

// Return keys of the provider's schema
func (a *Analyser) Keys(provider string) []string {
	var keys []string
	for key := range a.FlattenedSchema[provider] {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// Returns map of properties in this schema
func flattened(mymap map[string]*jsonschema.Schema, schema *jsonschema.Schema, parentKey string, level int) map[string]*jsonschema.Schema {
	for propertyName, propertySchema := range schema.Properties {
		key := parentKey + "/" + propertyName
		mymap[key] = propertySchema
		if level < 10 {
			mymap = flattened(mymap, propertySchema, key, level+1)
		}
	}

	return mymap
}
