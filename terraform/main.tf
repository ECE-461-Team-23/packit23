# Settings
locals {
  # General
  github_branch = "terraform"
  artifact_registry_repo_name = "container-repo"

  # test-app
  test_app_cloud_run_name = "test-app"
  test_app_image_name = "python-app2"

  # sql db
  google_sql_database_instance = "mysql-instance"

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
  region  = var.region
  zone    = "us-central1-a"
}

resource "google_artifact_registry_repository" "container-repo" {
  location = var.region
  repository_id = local.artifact_registry_repo_name
  description   = "Repository to store containers and artifacts"
  format        = "DOCKER"
  depends_on = [google_project_service.artifact_registry_api]
}

# Automatically build container for test-app
resource "google_cloudbuild_trigger" "app-trigger" {
  location = var.region

  github {
    owner = "packit461"
    name = "packit23"

    push {
      branch = local.github_branch
    }
  }

  filename = "test-app/cloudbuild.yaml"
  depends_on = [google_project_service.cloud_build_api]
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

resource "google_api_gateway_api" "api_gw" {
  provider = google-beta
  api_id = "my-api"
  project = var.project_id
}

resource "google_api_gateway_api_config" "api_cfg" {
  provider = google-beta
  api = google_api_gateway_api.api_gw.api_id
  api_config_id = "api-config"
  project = var.project_id

  openapi_documents {
    document {
      path = "api_spec.yaml"
      contents = filebase64("api_spec.yaml")
    }
  }
  lifecycle {
    create_before_destroy = true
  }
}

resource "google_api_gateway_gateway" "gw" {
  provider = google-beta
  region = var.region
  project = var.project_id
  api_config = google_api_gateway_api_config.api_cfg.id
  depends_on = [google_api_gateway_api_config.api_cfg]

}

# SQL
resource "google_sql_database_instance" "mysql-instance" {
  name             = local.google_sql_database_instance
  region           = var.region
  database_version = "MYSQL_8_0"
  settings {
    tier = "db-f1-micro"
  }

  deletion_protection  = "true"
}

resource "google_sql_database" "database" {
  name = "mysql_db"
  instance = local.google_sql_database_instance
}

resource "random_id" "db_name_suffix" {
  byte_length = 4
}

resource "google_sql_user" "read-user" {
  name     = "read-user"
  instance = local.google_sql_database_instance
  host     = "%"
  password = "strongpassword"
}

# Run containers for test-app (container image is overwritten in cloudbuild.yaml)
resource "google_cloud_run_service" "run_service" {
  name = local.test_app_cloud_run_name
  location = var.region

  template {
    spec {
      containers {
        #image = "us-docker.pkg.dev/cloudrun/container/placeholder:latest" # Placeholder
        image = "us-central1-docker.pkg.dev/${var.project_id}/${local.artifact_registry_repo_name}/${local.test_app_image_name}:latest"
      }
    }
  }
  # https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloud_run_service#metadata
  # How to connect a container to a SQL database

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [google_project_service.cloud_run_api]  # Waits for the Cloud Run API to be enabled
}

# Allow unauthenticated users to invoke the Cloud Run service
resource "google_cloud_run_service_iam_member" "run_all_users" {
  service  = google_cloud_run_service.run_service.name
  location = google_cloud_run_service.run_service.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

output "service_url" {
  value = google_cloud_run_service.run_service.status[0].url
}