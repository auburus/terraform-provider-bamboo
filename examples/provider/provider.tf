terraform {
  required_providers {
    bamboo = {
      source = "local/providers/bamboo"
    }
  }
}
provider "bamboo" {
  # example configuration here
}


data "bamboo_unknown_data_source" "example" {
}