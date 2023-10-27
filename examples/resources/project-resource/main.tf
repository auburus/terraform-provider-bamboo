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

resource "bamboo_project" "new_project" {
  key = "P4"
  name = "Project4"
  description = "Terraform generated project"
}

output "project" {
  value = bamboo_project.new_project
}
