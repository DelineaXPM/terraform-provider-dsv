# Thycotic DevOps Secrets Vault - Terraform Provider

The [Thycotic](https://thycotic.com/)
[DevOps Secrets Vault](https://thycotic.com/products/devops-secrets-vault-password-management/)
[Terraform](https://www.terraform.io/) Provider makes Secrets data available and provisions client
secrets for existing roles.

## Installation

Terraform has [installation instructions](https://www.terraform.io/docs/plugins/basics.html#installing-plugins).
The binaries can be downloaded [here]().

## Usage

To run the example, create a `terraform.tfvars`:

```terraform
dsv_client_id     = "a54bc1b6-7dd7-4fb1-a8ba-bbfa81820e40"
dsv_client_secret = "xxxxxxxxxxxxxxxxxxxxxxxxx-xxxxxxxxxxx-xxxxx"
dsv_tenant        = "mytenant"
dsv_secret_path   = "/path/to/a/test/secret"
```
