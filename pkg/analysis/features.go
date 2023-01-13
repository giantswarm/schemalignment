package analysis

import (
	"fmt"
	"sort"
	"strings"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
)

const unspecifiedNumber = -1

// Return "features" found in all schemas.
//
// A feature is a unique combination of the JSON schema capabilities used.
func (a *Analyser) Features() map[string][]string {
	featuresMap := make(map[string][]string)

	for _, clusterApp := range a.ClusterApps {
		for key := range a.FlattenedSchema[clusterApp.ProviderName] {
			if a.FlattenedSchema[clusterApp.ProviderName][key] != nil {
				features := extractFeatures(a.FlattenedSchema[clusterApp.ProviderName][key])
				featuresString := strings.Join(features, " ")

				_, exists := featuresMap[featuresString]
				if exists {
					featuresMap[featuresString] = append(featuresMap[featuresString], a.FlattenedSchema[clusterApp.ProviderName][key].Location)
				} else {
					featuresMap[featuresString] = []string{a.FlattenedSchema[clusterApp.ProviderName][key].Location}
				}
			}
		}
	}

	return featuresMap
}

func extractFeatures(s *jsonschema.Schema) []string {
	// List of feature strings to collect.
	features := []string{}

	// type
	if len(s.Types) > 1 {
		features = append(features, "multiple_types")
		types := s.Types
		sort.Strings(types)
		features = append(features, fmt.Sprintf("type=%s", strings.Join(types, ",")))
	} else if len(s.Types) == 1 {
		features = append(features, "single_type")
		features = append(features, fmt.Sprintf("type=%s", s.Types[0]))
	} else if len(s.Types) == 0 {
		features = append(features, "no_type")
	}

	// additionalProperties
	if s.AdditionalProperties != nil {
		features = append(features, "additionalProperties")
		switch v := s.AdditionalProperties.(type) {
		case bool:
			features = append(features, "additional_properties_boolean")
		case *jsonschema.Schema:
			features = append(features, "additional_properties_object")
		default:
			features = append(features, fmt.Sprintf("additional_properties_unknown=%s", v))
		}
	}

	// patternProperties
	if s.PatternProperties != nil {
		features = append(features, "patternProperties")
		switch v := s.AdditionalProperties.(type) {
		case bool:
			features = append(features, "pattern_properties_boolean")
		case *jsonschema.Schema:
			features = append(features, "pattern_properties_object")
		default:
			features = append(features, fmt.Sprintf("pattern_properties_unknown=%s", v))
		}
	}

	// additionalItems
	if s.AdditionalItems != nil {
		features = append(features, "additionalItems")
		switch v := s.AdditionalProperties.(type) {
		case bool:
			features = append(features, "additional_items_boolean")
		case *jsonschema.Schema:
			features = append(features, "additional_items_object")
		default:
			features = append(features, fmt.Sprintf("additional_items_unknown=%s", v))
		}
	}

	if len(s.AllOf) > 0 {
		features = append(features, "allOf")
	}

	if len(s.AnyOf) > 0 {
		features = append(features, "anyOf")
	}

	if s.Constant != nil {
		features = append(features, "constant")
	}

	if s.Contains != nil {
		features = append(features, "contains")
	}

	if s.Default != nil {
		features = append(features, "default")
	}

	if s.DependentRequired != nil {
		features = append(features, "dependentRequired")
	}

	if s.Deprecated {
		features = append(features, "deprecated")
	}

	if s.Enum != nil {
		features = append(features, "enum")
	}

	if s.Examples != nil {
		features = append(features, "examples")
	}

	if s.Else != nil {
		features = append(features, "else")
	}

	if s.ExclusiveMaximum != nil {
		features = append(features, "exclusiveMaximum")
	}

	if s.ExclusiveMinimum != nil {
		features = append(features, "exclusiveMinimum")
	}

	if s.Format != "" {
		features = append(features, "format")
	}

	if s.If != nil {
		features = append(features, "if")
	}

	if s.MaxContains != unspecifiedNumber {
		features = append(features, "maxContains")
	}

	if s.MaxItems != unspecifiedNumber {
		features = append(features, "maxItems")
	}

	if s.MaxLength != unspecifiedNumber {
		features = append(features, "maxLength")
	}

	if s.MaxProperties != unspecifiedNumber {
		features = append(features, "maxProperties")
	}

	if s.Maximum != nil {
		features = append(features, "maximum")
	}

	if s.MinContains != 1 {
		features = append(features, "minContains")
	}

	if s.MinItems != unspecifiedNumber {
		features = append(features, "minItems")
	}

	if s.MinLength != unspecifiedNumber {
		features = append(features, "minLength")
	}

	if s.MinProperties != unspecifiedNumber {
		features = append(features, "minProperties")
	}

	if s.Minimum != nil {
		features = append(features, "minimum")
	}

	if s.MultipleOf != nil {
		features = append(features, "multipleOf")
	}

	if s.Not != nil {
		features = append(features, "not")
	}

	if len(s.OneOf) > 0 {
		features = append(features, "oneOf")
	}

	if s.Pattern != nil {
		features = append(features, "pattern")
	}

	if s.PrefixItems != nil {
		features = append(features, "prefixItems")
	}

	if len(s.Required) > 0 {
		features = append(features, "required_properties")
	}

	if s.ReadOnly {
		features = append(features, "readOnly")
	}

	if s.Then != nil {
		features = append(features, "then")
	}

	if s.UniqueItems {
		features = append(features, "uniqueItems")
	}

	if s.WriteOnly {
		features = append(features, "writeOnly")
	}

	sort.Strings(features)
	return features
}
