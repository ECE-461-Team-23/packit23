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

resource "google_cloudbuild_trigger" "app-trigger" {
  location = "us-central1"

  github {
    owner = "packit461"
    name = "packit23"

    push {
      branch = "terraform-container-build"
    }
  }

  filename = "test-app/cloudbuild.yaml"
}

resource "google_artifact_registry_repository" "container-repo" {
  location      = "us-central1"
  repository_id = "container-repo"
  description   = "Repository to store containers and artifacts"
  format        = "DOCKER"
}

resource "google_cloud_run_service" "run_service" {
  name = "app"
  location = "us-central1"

  template {
    spec {
      containers {
        image = "us-central1-docker.pkg.dev/${var.project_id}/my-docker-repo/python-app:latest"
      }
    }
  }
  # https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_run_service#metadata
  # How to connect a container to a SQL database

  traffic {
    percent         = 100
    latest_revision = true
  }

  # Waits for the Cloud Run API to be enabled
  depends_on = [google_project_service.cloud_run_api]
}

# Allow unauthenticated users to invoke the service
resource "google_cloud_run_service_iam_member" "run_all_users" {
  service  = google_cloud_run_service.run_service.name
  location = google_cloud_run_service.run_service.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

output "service_url" {
  value = google_cloud_run_service.run_service.status[0].url
}


## Enable services ##

resource "google_project_service" "cloud_run_api" {
  service = "run.googleapis.com"
  disable_on_destroy = true
}

resource "google_project_service" "cloud_build_api" {
  service = "cloudbuild.googleapis.com"
  disable_on_destroy = true
}

resource "google_project_service" "artifact_registry_api" {
  service = "artifactregistry.googleapis.com"
  disable_on_destroy = true
}

