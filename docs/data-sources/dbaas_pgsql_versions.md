---
subcategory: "Database as a Service - Postgres"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_pg_versions"
sidebar_current: "docs-ionoscloud_pg_versions"
description: |-
  Get information on DbaaS PgSql Versions
---

# ionoscloud\_pg_versions

The DbaaS Postgres Versions data source can be used to search for and retrieve list of available postgres versions for a specific cluster or for all clusters.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### Retrieve list of postgres versions for a specific cluster
```hcl
data "ionoscloud_pg_versions" "test_ds_dbaas_versions" {
	cluster_id = ionoscloud_pg_cluster.test_dbaas_cluster.id
}
```

### Retrieve list of postgres versions for all clusters
```hcl
data "ionoscloud_pg_versions" "test_ds_dbaas_versions" {
}
```

## Argument Reference

* `cluster_id` - (Optional) The unique ID of the cluster.

If `cluster_id` is not provided the data source will return the list of postgres version for all cluster.

## Attributes Reference

The following attributes are returned by the datasource:

* `cluster_id` - Id of the cluster
* `postgres_versions` - list of PostgreSQL versions.