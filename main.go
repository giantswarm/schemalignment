package main

import (
	"encoding/json"
	"log"

	"github.com/giantswarm/microerror"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"github.com/skratchdot/open-golang/open"

	"github.com/giantswarm/schemalignment/pkg/analysis"
	"github.com/giantswarm/schemalignment/pkg/server"
)

var (
	// Define all cluster apps to analyse schema for here.
	clusterApps = []analysis.ClusterApp{
		{
			ProviderName:  "AWS",
			RepositoryURL: "https://github.com/giantswarm/cluster-aws",
			SchemaURL:     "https://raw.githubusercontent.com/giantswarm/cluster-aws/master/helm/cluster-aws/values.schema.json",
		},
		{
			ProviderName:  "Azure",
			RepositoryURL: "https://github.com/giantswarm/cluster-azure",
			SchemaURL:     "https://raw.githubusercontent.com/giantswarm/cluster-azure/main/helm/cluster-azure/values.schema.json",
		},
		{
			ProviderName:  "Cloud Director",
			RepositoryURL: "https://github.com/giantswarm/cluster-cloud-director",
			SchemaURL:     "https://raw.githubusercontent.com/giantswarm/cluster-cloud-director/main/helm/cluster-cloud-director/values.schema.json",
		},
		{
			ProviderName:  "GCP",
			RepositoryURL: "https://github.com/giantswarm/cluster-gcp",
			SchemaURL:     "https://raw.githubusercontent.com/giantswarm/cluster-gcp/main/helm/cluster-gcp/values.schema.json",
		},
		{
			ProviderName:  "VSphere",
			RepositoryURL: "https://github.com/giantswarm/cluster-vsphere",
			SchemaURL:     "https://raw.githubusercontent.com/giantswarm/cluster-vsphere/main/helm/cluster-vsphere/values.schema.json",
		},
	}

	url = "http://localhost:8080/"
)

// Data is a big data structure we deliver to the web UI as JSON,
// containing all the information we want to display to users.
type Data struct {
	ClusterApps []analysis.ClusterApp
	Providers   []string

	// List of all properties with hierarchical name
	PropertyKeys []string

	// Map of properties (key) and information per provider.
	// Example:
	// "/foo/bar": {"AWS": {"type": ["string"], "title": "...", "description": "..."}}
	PropertiesAndProviders map[string]map[string]analysis.ProviderPropertySummary

	// Map of features (key) and array of locations where this feature occurs.
	Features map[string][]string
}

func main() {
	analyser, err := analysis.New(clusterApps)
	if err != nil {
		log.Fatal(microerror.Mask(err))
	}

	data := Data{
		ClusterApps:            clusterApps,
		Providers:              analyser.Providers(),
		PropertyKeys:           analyser.HierarchicalKeys(),
		PropertiesAndProviders: analyser.MergedSchemas(),
		Features:               analyser.Features(),
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

	err = server.Serve(8080, dataJson)
	if err != nil {
		log.Fatal(err)
	}
}
