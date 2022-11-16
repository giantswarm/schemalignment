package main

import (
	"encoding/json"
	"log"
	"sort"

	"github.com/giantswarm/microerror"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"github.com/skratchdot/open-golang/open"

	"github.com/giantswarm/schemalignment/pkg/analysis"
	"github.com/giantswarm/schemalignment/pkg/server"
)

var (
	schemaUrls = map[string]string{
		"AWS":            "https://raw.githubusercontent.com/giantswarm/cluster-aws/master/helm/cluster-aws/values.schema.json",
		"Cloud Director": "https://raw.githubusercontent.com/giantswarm/cluster-cloud-director/main/helm/cluster-cloud-director/values.schema.json",
		"GCP":            "https://raw.githubusercontent.com/giantswarm/cluster-gcp/main/helm/cluster-gcp/values.schema.json",
		"OpenStack":      "https://raw.githubusercontent.com/giantswarm/cluster-openstack/main/helm/cluster-openstack/values.schema.json",
		"VSphere":        "https://raw.githubusercontent.com/giantswarm/cluster-vsphere/main/helm/cluster-vsphere/values.schema.json",

		// TODO: add Azure
		// once it's in https://github.com/giantswarm/cluster-azure/tree/main/helm/cluster-azure
		//"Azure": ""
		// "https://raw.githubusercontent.com/giantswarm/cluster-azure/main/helm/cluster-azure/values.schema.json",
	}

	url = "http://localhost:8080/"
)

type Data struct {
	Providers []string `json:"providers"`

	// List of all properties with hierarchical name
	PropertyKeys []string `json:"property_keys"`

	// Map of properties (key) and array of provides per key
	PropertiesAndProviders map[string][]string `json:"properties_and_providers"`
}

func main() {
	analyser, err := analysis.New(schemaUrls)
	if err != nil {
		log.Fatal(microerror.Mask(err))
	}

	full := analyser.FullSchemas()
	var keys []string
	for key := range full {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var providers []string
	for provider := range schemaUrls {
		providers = append(providers, provider)
	}
	sort.Strings(providers)

	data := Data{
		Providers:              providers,
		PropertyKeys:           keys,
		PropertiesAndProviders: full,
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Opening browser at %s", url)
	err = open.Start(url)
	if err != nil {
		log.Fatal(err)
	}

	server.Serve(8080, dataJson)
}
