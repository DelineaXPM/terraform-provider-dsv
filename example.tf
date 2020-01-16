variable "dsv_client_id" {
  type = string
}

variable "dsv_client_secret" {
  type = string
}

variable "dsv_tenant" {
  type = string
}

variable "dsv_role_name" {
  type = "string"
}

variable "dsv_secret_path" {
  type = "string"
}

provider "dsv" {
  client_id     = var.dsv_client_id
  client_secret = var.dsv_client_secret
  tenant        = var.dsv_tenant
}

data "dsv_secret" "username" {
  path    = var.dsv_secret_path
  element = "username"
}

data "dsv_secret" "password" {
  path    = var.dsv_secret_path
  element = "password"
}

data "dsv_role" "existing_role" {
  name = var.dsv_role_name
}

resource "dsv_client" "new_client" {
  role = data.dsv_role.existing_role.name
}

output "client_id" {
  value = dsv_client.new_client.client_id
}

output "client_secret" {
  value = dsv_client.new_client.client_secret
}

output "username" {
  value = data.dsv_secret.username.contents
}

output "password" {
  value = data.dsv_secret.password.contents
}
