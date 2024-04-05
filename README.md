# Terraform Provider String Template

A simple provider that exposes a singular provider function named `template`. This function consumes an
[HCL string template](https://developer.hashicorp.com/terraform/language/expressions/strings) along with a values map to
use in the templating process.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.8.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

```hcl
terraform {
  required_providers {
    vault = {
      source = "mbillow/string-template"
      version = "???"
    }
  }
}

variable "demo_template" {
  type = string
  description = "User provided HCL template string."
  // Since we are defining a default inline, we need to use %% and $$ to get the literals % and $.
  // Your templates don't need double characters if they are coming from outside your HCL.
  default = <<-EOT
  %%{ for ip in split(",", ip_addresses) ~}
  $${ ip }
  %%{ endfor ~}
  EOT
}

locals = {
  template_vars = {
    ip_addresses = "1.1.1.1,2.2.2.2,3.3.3.3"
  }
}

output "demo" {
  value = provider::string-template::template(
    var.demo_template,
    local.template_vars
  )
}
```

Output:
```text
1.1.1.1
2.2.2.2
3.3.3.3
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (
see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin`
directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
