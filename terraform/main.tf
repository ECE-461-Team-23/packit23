# Settings
locals {
  # General
  github_branch = "containerize_package_rater"
  artifact_registry_repo_name = "container-repo"

  # package-rater-app
  package_rater_app_cloud_run_name = "package-rater-app"
  package_rater_app_image_name = "package-rater-image"
}

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
  region  = "us-central1"
  zone    = "us-central1-a"
}

resource "google_artifact_registry_repository" "container_repo" {
  location = "us-central1"
  repository_id = local.artifact_registry_repo_name
  description   = "Repository to store containers and artifacts"
  format        = "DOCKER"
  depends_on = [google_project_service.artifact_registry_api]
}

# Automatically build container for package-rater-app
resource "google_cloudbuild_trigger" "package_rater_app_trigger" {
  location = "us-central1"

  github {
    owner = "packit461"
    name = "packit23"

    push {
      branch = local.github_branch
    }
  }

  filename = "package-rater-app/cloudbuild.yaml"
  depends_on = [google_project_service.cloud_build_api]
}

# Run containers for package-rater-app (container image is overwritten in cloudbuild.yaml)
resource "google_cloud_run_service" "run_service" {
  name = local.package_rater_app_cloud_run_name
  location = "us-central1"

  template {
    spec {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/placeholder:latest" # Placeholder
        # "us-central1-docker.pkg.dev/${var.project_id}/${local.artifact_registry_repo_name}/${local.package_rater_app_image_name}:latest"
        env {
          name = "GITHUB_TOKEN"
          value_from {
            secret_key_ref {
              name = "GITHUB_TOKEN"
              key  = "latest"
            }
          }
        }
        env {
          name = "LOG_FILE"
          value = "/var/log/output.log"
        }
        env {
          name = "LOG_LEVEL"
          value = "2"
        }
      }
      service_account_name = google_service_account.package_rater_service_account.email
    }
  }
  # https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_run_service#metadata
  # How to connect a container to a SQL database

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [google_project_service.cloud_run_api,  # Waits for the Cloud Run API to be enabled
                google_secret_manager_secret_iam_member.package_rater_access] # Make sure service account is attached to policy to give access to secret token
}

# Allow unauthenticated users to invoke the Cloud Run service
resource "google_cloud_run_service_iam_member" "run_all_users" {
  service  = google_cloud_run_service.run_service.name
  location = google_cloud_run_service.run_service.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# Create Secret Manager to store value of Github Token
resource "google_secret_manager_secret" "github_token_manager" {
  secret_id = "GITHUB_TOKEN"

  replication {
    automatic = true
  }

  depends_on = [ google_project_service.secret_manager_api ]
}

# Create a new version of "Github Token" secret
resource "google_secret_manager_secret_version" "github_token_manager_version" {
  secret   = google_secret_manager_secret.github_token_manager.id
  # version  = 1
  secret_data = var.github_token
}

# Give Package Rater SA access to "Github Token" secret
resource "google_secret_manager_secret_iam_member" "package_rater_access" {
  secret_id = google_secret_manager_secret.github_token_manager.secret_id
  role = "roles/secretmanager.secretAccessor"
  member = "serviceAccount:${google_service_account.package_rater_service_account.email}"
}


output "service_url" {
  value = google_cloud_run_service.run_service.status[0].url
}

## Service Accounts ##

resource "google_service_account" "package_rater_service_account" {
  account_id = "package-rater-sa"
  display_name = "Service account for package rater containers"
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

resource "google_project_service" "secret_manager_api" {
  service = "secretmanager.googleapis.com"
  disable_on_destroy = true
}
