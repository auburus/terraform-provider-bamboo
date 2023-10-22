# Terraform Provider Bamboo

> This is a work in progress and it won't be considered stable until
> version 1.0 is released (at which time I'll seek publishing this
> to the registry).

This manages bamboo resources using terraform.

## Why?
There are more up to date tools than Bamboo for CI, but my employer mandates its use
so lets make the best out of it.

[Bamboo API Docs](https://developer.atlassian.com/server/bamboo/rest/api-group-api/#api-api-latest-get)


# Using the provider
WIP

# Developing the Provider

### Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

```shell
make build
```

This will create the provier compiled binary in `~/go/bin`.
Add this folder to the "overrides" settings for terraform.

```hcl
# ~/.terraformrc
provider_installation {
    dev_overrides {
        # For some reason `~` is not supported, need to set the full path
        "local/providers/bamboo" = "/home/<your-username>/go/bin/"
    }

    # This line is mandatory to make other providers use the
    # terraform registry as ususal
    direct {}
}

```
