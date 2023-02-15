// Package analysis analyses schemas and creates comparison.
package analysis

import (
	"log"
	"sort"

	"github.com/giantswarm/microerror"
	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"golang.org/x/exp/slices"
)

// ClusterApp holds information on a Giant Swarm cluster app for a Cluster API provider.
type ClusterApp struct {
	// User-friendly name of the Cluster API infrastructure provider.
	ProviderName string

	// URL of the GitHub repo landing page.
	RepositoryURL string

	// URL of the schema file in JSON for download.
	SchemaURL string
}

// Analyser is the agent that performs comparison and analysis on the schemas.
type Analyser struct {
	// The cluster apps handed over to the agent.
	ClusterApps []ClusterApp

	// Schemas holds the full schema information in the original hierarchical
	// form. Map key is the provider name.
	Schemas map[string]*jsonschema.Schema

	// FlattenedSchema holds a flattened schema where all property names
	// are brought into a path like `/main/sub/subsub`.
	// Top level map key is the provider name.
	// Second level map key is the property name in path form.
	FlattenedSchema map[string]map[string]*jsonschema.Schema
}

// PropertyDetails is an extract of information regarding a single property in a schema.
type PropertyDetails struct {
	Key             string
	KeyHierarchical string
	Types           []string
	DefaultValue    string
}

func New(clusterApps []ClusterApp) (*Analyser, error) {
	a := &Analyser{
		ClusterApps:     clusterApps,
		Schemas:         make(map[string]*jsonschema.Schema),
		FlattenedSchema: make(map[string]map[string]*jsonschema.Schema),
	}

	for _, clusterApp := range clusterApps {
		schema, err := jsonschema.Compile(clusterApp.SchemaURL)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		a.Schemas[clusterApp.ProviderName] = schema
		a.FlattenedSchema[clusterApp.ProviderName] = make(map[string]*jsonschema.Schema)
		a.FlattenedSchema[clusterApp.ProviderName] = flattenedSchema(schema, clusterApp.ProviderName)
	}

	return a, nil
}

// Providers returns the list of provider names in the order of definition.
func (a *Analyser) Providers() (providers []string) {
	for _, clusterApp := range a.ClusterApps {
		providers = append(providers, clusterApp.ProviderName)
	}

	return providers
}

// Returns a list of all hierarchical property keys from all schemas.
func (a *Analyser) HierarchicalKeys() (keys []string) {
	for key := range a.MergedSchemas() {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// MergedSchemas returns a map of properties occuring over all analysed schemas,
// with the list of providers as their value.
func (a *Analyser) MergedSchemas() map[string][]string {
	// Create complete list of all property keys with a list of the providers having them.
	fullSchemas := make(map[string][]string)
	for _, clusterApp := range a.ClusterApps {
		// Collect all keys
		for key := range a.FlattenedSchema[clusterApp.ProviderName] {
			_, ok := fullSchemas[key]
			if ok {
				fullSchemas[key] = append(fullSchemas[key], clusterApp.ProviderName)
			} else {
				fullSchemas[key] = []string{clusterApp.ProviderName}
			}
		}
	}

	return fullSchemas
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

// Flattens the schema, using a hierarchical key.
func flattenedSchema(schema *jsonschema.Schema, providerName string) map[string]*jsonschema.Schema {
	var mymap = make(map[string]*jsonschema.Schema)
	return flattened(mymap, schema, providerName, "", 0)
}

// Returns map of properties in this schema by
// recursing into a schema's properties.
func flattened(mymap map[string]*jsonschema.Schema, schema *jsonschema.Schema, providerName, parentKey string, level int) map[string]*jsonschema.Schema {
	if slices.Contains(schema.Types, "object") {
		if len(schema.Properties) == 0 && schema.AdditionalProperties == nil {
			log.Printf("Warning: provider %q object %s has no 'properties'/'additionalProperties' defined.", providerName, parentKey)
		} else {
			for propertyName, propertySchema := range schema.Properties {
				key := parentKey + "/" + propertyName
				if slices.Contains(propertySchema.Types, "array") {
					key = key + "[*]"
				}
				mymap[key] = propertySchema

				if level < 10 {
					mymap = flattened(mymap, propertySchema, providerName, key, level+1)
				}
			}
		}
	} else if slices.Contains(schema.Types, "array") {
		// Array properties.
		//parentKey = parentKey + "[*]"

		if schema.Items2020 != nil {
			if len(schema.Items2020.Types) == 1 {
				itemType := schema.Items2020.Types[0]
				if itemType == "object" {
					for propertyName, propertySchema := range schema.Items2020.Properties {
						key := parentKey + "/" + propertyName
						mymap[key] = propertySchema

						if level < 10 {
							mymap = flattened(mymap, propertySchema, providerName, key, level+1)
						}
					}
				} else if itemType == "string" {
					// Fine to ignore this. Nothing to add for string and number types.
				} else if itemType == "number" {
					// Fine to ignore this. Nothing to add for string and number types.
				} else if itemType == "integer" {
					// Fine to ignore this. Nothing to add for string and number types.
				} else {
					log.Printf("Debug: provider %q array %s items skipped because of unhandled type %s", providerName, parentKey, itemType)
				}
			} else {
				log.Printf("Warning: provider %q array %s has multiple types %s, skipped.", providerName, parentKey, schema.Items2020.Types)
			}
		} else if schema.Items != nil {
			switch v := schema.Items.(type) {
			case *jsonschema.Schema:
				if len(v.Types) == 1 {
					itemType := v.Types[0]
					log.Printf("Debug: provider %q array %s items have type %q", providerName, parentKey, itemType)

					if itemType == "object" {
						for propertyName, propertySchema := range v.Properties {
							key := parentKey + "/" + propertyName
							mymap[key] = propertySchema

							if level < 10 {
								mymap = flattened(mymap, propertySchema, providerName, key, level+1)
							}
						}
					} else if itemType == "string" {
						// Fine to ignore this. Nothing to add for string and number types.
					} else if itemType == "number" {
						// Fine to ignore this. Nothing to add for string and number types.
					} else if itemType == "integer" {
						// Fine to ignore this. Nothing to add for string and number types.
					} else {
						log.Printf("Debug: provider %q array %s items skipped because of unhandled type %s", providerName, parentKey, itemType)
					}
				} else {
					log.Printf("Warning: provider %q array %s has multiple types %s, skipped.", providerName, parentKey, v.Types)
				}
			case []*jsonschema.Schema:
				log.Printf("Debug: items is type []*jsonschema.Schema")
			}
		} else {
			mymap[parentKey] = nil
			log.Printf("Warning: provider %q array %s provides no 'items' schema.", providerName, parentKey)
		}
	}

	return mymap
}
