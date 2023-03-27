terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.51.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = "us-east1"
  zone    = "us-east1-a"
}

# resource "google_compute_network" "vpc_network" {
#   name = "terraform-network"
# }

resource "google_cloudbuild_trigger" "app-trigger" {
  location = "us-central1"

  github {
    owner = "packit461"
    name = "packit23"

    push {
      branch = "terraform-container-build"
    }
  }
  # included_files = ["test-app/*"] # Only update container if the folder is updated

  filename = "test-app/cloudbuild.yaml"

  # build {
  #   step {
  #     name = "gcr.io/cloud-builders/docker"
  #     args = ["build", "-t", "gcr.io/${var.project_id}/quickstart-image:$COMMIT_SHA", "test-app/"]
  #   }
  #   step {
  #     name = "gcr.io/cloud-builders/docker"
  #     args = ["push", "gcr.io/${var.project_id}/quickstart-image:$COMMIT_SHA"]
  #   }    
  # }
}

resource "google_artifact_registry_repository" "my-repo" {
  location      = "us-central1"
  repository_id = "my-docker-repo"
  description   = "example docker repository"
  format        = "DOCKER"
}