terraform {
  required_providers {
    bamboo = {
      source = "local/providers/bamboo"
    }
  }
}
provider "bamboo" {
  url = "http://localhost:8085/"
  username = "admin"
  password = "admin"
}

data "bamboo_project" "project1" {
  key = "P1"
}


output "project1" {
  value = data.bamboo_project.project1
}