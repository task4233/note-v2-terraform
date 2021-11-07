# ref: https://learn.hashicorp.com/collections/terraform/gcp-get-started

# variables set in `terraform.tfvars` which should be contained `.gitignore`.
# variable "project" {}
# variable "credentials_file" {}
# variable "region" {
#   default = "us-west1"
# }
# variable "zone" {
#   default = "us-west1-b"
# }
# variable "gce_ssh_user" {}
# variable "gce_ssh_pub_key_file" {}

terraform {
  required_providers {
    # google = {
    #   source  = "hashicorp/google"
    #   version = "3.5.0"
    # }
    log = {
      version = "0.1.0"
      source  = "terraform.local/local/log"
    }
  }
}

provider "log" {
  host = "http://localhost:19090"
}

resource "log_order" "item" {
  items = [
    {
      log = {
        body = "hoge"
      }
    },
    {
      log = {
        body = "fuga"
      }
    },
  ]
}

# provider "google" {
#   credentials = file(var.credentials_file)

#   project = var.project
#   region  = var.region
#   zone    = var.zone
# }

# # ready vpc_network for terraform
# resource "google_compute_network" "vpc_network" {
#   name = "terraform-network"
# }

# # firewall configuration for icmp & ssh
# resource "google_compute_firewall" "allow_icmp" {
#   name    = "allow-icmp"
#   network = google_compute_network.vpc_network.name
#   allow {
#     protocol = "icmp"
#   }
# }

# resource "google_compute_firewall" "allow_ssh" {
#   name    = "allow-ssh"
#   network = google_compute_network.vpc_network.name
#   allow {
#     protocol = "tcp"
#     ports    = ["2526"] # should be changed
#   }
# }

# resource "google_compute_firewall" "allow_web" {
#   name    = "allow-web"
#   network = google_compute_network.vpc_network.name
#   allow {
#     protocol = "tcp"
#     ports    = ["443"]
#   }
# }

# # use configuration in free tier
# # ref: https://cloud.google.com/free/docs/gcp-free-tier/#free-tier-usage-limits
# resource "google_compute_instance" "vm_instance" {
#   name         = "terraform-instance"
#   machine_type = "e2-micro"

#   boot_disk {
#     initialize_params {
#       size  = 30
#       type  = "pd-standard"
#       image = "debian-cloud/debian-9"
#     }
#   }

#   network_interface {
#     network = google_compute_network.vpc_network.name
#     access_config {
#     }
#   }

#   metadata = {
#     ssh-keys = "${var.gce_ssh_user}:${file(var.gce_ssh_pub_key_file)}" # use ssh-key set in `terraform.tfvars`.
#   }
# }
