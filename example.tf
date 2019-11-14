variable "dsv_client_id" {
  type = string
}

variable "dsv_client_secret" {
  type = string
}

variable "dsv_tenant" {
  type = string
}

variable "dsv_secret_path" {
  type = "string"
}

provider "dsv" {
  client_id     = var.dsv_client_id
  client_secret = var.dsv_client_secret
  tenant        = var.dsv_tenant
}

data "dsv_secret" "password" {
  path    = var.dsv_secret_path
  element = "password"
}

output "password" {
  value = data.dsv_secret.password.contents
}
