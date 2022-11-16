# Cluster app config schema alignment

Example screenshot

<img width="1439" alt="image" src="https://user-images.githubusercontent.com/273727/202226009-9e43bbac-9130-4007-8a81-419e77168799.png">

## Usage

Simply execute the program like this ...

```nohighlight
go run main.go
```

... and see the result in the opening browser.

## Background

Cluster apps are our means to provision workload clusters with Cluster API. Giant Swarm provides a specific cluster app for each provider.

These apps are configured via user values. The structure of these is determined by the cluster app's user values schema.

| Provider | Cluster app repository | Schema |
|-|-|-|
| AWS | [cluster-aws](https://github.com/giantswarm/cluster-aws) | [Schema](https://raw.githubusercontent.com/giantswarm/cluster-aws/master/helm/cluster-aws/values.schema.json) |
| Cloud Director | [cluster-cloud-director](https://github.com/giantswarm/cluster-cloud-director) | [Schema](https://raw.githubusercontent.com/giantswarm/cluster-cloud-director/main/helm/cluster-cloud-director/values.schema.json) |
| GCP | [cluster-gcp](https://github.com/giantswarm/cluster-gcp) | [Schema](https://raw.githubusercontent.com/giantswarm/cluster-gcp/main/helm/cluster-gcp/values.schema.json) |
| OpenStack | [cluster-openstack](https://github.com/giantswarm/cluster-openstack) | [Schema](https://raw.githubusercontent.com/giantswarm/cluster-openstack/main/helm/cluster-openstack/values.schema.json) |
| VSphere | [cluster-vsphere](https://github.com/giantswarm/cluster-vsphere) | [Schema](https://raw.githubusercontent.com/giantswarm/cluster-vsphere/main/helm/cluster-vsphere/values.schema.json) |

For the sake of simplicity and good UX, we think that these schemas should be aligned as much as possible across providers.

This repository provides tooling to help with alignment.
