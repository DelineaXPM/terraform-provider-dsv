# Thycotic DevOps Secrets Vault - Terraform Provider

The [Thycotic](https://thycotic.com/)
[DevOps Secrets Vault](https://thycotic.com/products/devops-secrets-vault-password-management/) (DSV)
[Terraform](https://www.terraform.io/) Provider makes Secrets data available and provisions client
secrets for existing roles.

## Installation

### Install the executable

```bash
go get github.com/thycotic/terraform-provider-dsv
```

The executable is now available in `$GOPATH/bin`

### Install Terraform

1. Download the platform-specific static executable [here](https://www.terraform.io/downloads.html).
2. Copy or link it to a directory in `$PATH` of the target environment.


### Make the executable available to Terraform

Copy or link the executable into the Terraform _plugins directory_ of the target
environment. Refer to the Terraform [installation instructions](https://www.terraform.io/docs/plugins/basics.html#installing-plugins)
for the platform-specific location. It is `~/.terraform.d/plugins` on Linux and
MacOS.

## Examples

`example.tf` retrieves a secret and role from DSV then creates a `client_id` for
the role.

### Run `example.tf`

1. Create a `terraform.tfvars`:

    ```terraform
    dsv_client_id     = "a54bc1b6-7dd7-4fb1-a8ba-bbfa81820e40"
    dsv_client_secret = "xxxxxxxxxxxxxxxxxxxxxxxxx-xxxxxxxxxxx-xxxxx"
    dsv_tenant        = "mytenant"
    dsv_role_name     = "example-role"
    dsv_secret_path   = "/path/to/a/test/secret"
    ```

2. "Apply" the example Terraform plan:

    ```bash
    terraform apply -auto-approve
    ```

3. (Optional) delete the newly created `client_id`

    ```bash
    terraform destroy -auto-approve
    ```
