terraform {
  required_providers {
    loriot = {
      source = "registry.terraform.io/hashicorp/loriot"
    }
  }
}

provider "loriot" {
  host = "https://rentokil-stage.loriot.io"
  #host = "http://localhost:5555"
  key = "AAAAJAODvBmSt6tDCGBokPDEF_rZbMbNXkj_tOblnKvhq5aCg"
}

data "loriot_user" "example" {}
data "loriot_userusage" "example" {}
data "loriot_app" "example" {
  app_id = "BE01000D"
}

data "loriot_app" "working" {
  app_id = "BE010003"
}

resource "loriot_app" "my-app" {
  name                = "Terraform Test 2"
  devices_limit       = 200
  mcast_devices_limit = 1
}

data "loriot_apptoken" "example" {
  app_id = data.loriot_app.example.app_id
}

output "user" {
  value = data.loriot_user.example
}

output "userusage" {
  value = data.loriot_userusage.example
}

output "app" {
  value = data.loriot_app.example
}

output "app-working" {
  value = data.loriot_app.working
}
output "apptoken" {
  value = data.loriot_apptoken.example
}

output "my-app" {
  value = loriot_app.my-app
}
