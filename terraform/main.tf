# # API Gateway

# resource "google_api_gateway_api" "api_gw" {
#   provider = google-beta
#   api_id = "my-api"
# }

# resource "google_api_gateway_api_config" "api_cfg" {
#   provider = google-beta
#   api = google_api_gateway_api.api_gw.api_id
#   api_config_id = "api-config"

#   openapi_documents {
#     document {
#       path = "api_spec.yaml"
#       contents = base64encode("api_spec.yaml")
#     }
#   }
#   lifecycle {
#     create_before_destroy = true
#   }
# }

# resource "google_api_gateway_gateway" "gw" {
#   provider = google-beta
#   region = var.region
#   api_config = google_api_gateway_api_config.api_cfg.id
#   gateway_id = "api_gw"

#   depends_on = [google_api_gateway_api_config.api_cfg]

# }


# Settings
locals {
  # General
  github_branch = "terraform"
  artifact_registry_repo_name = "container-repo"
  region = "us-central1"

  # package-rater-app
  package_rater_app_cloud_run_name = "package-rater-app"
  package_rater_app_image_name = "package-rater-image"

  # SQL
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
  region  = local.region
  zone    = "us-central1-a"
}

resource "google_artifact_registry_repository" "container_repo" {
  location = local.region
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

  filename = "package_rater/cloudbuild.yaml"
  depends_on = [google_project_service.cloud_build_api]
}

# Run containers for package-rater-app (container image is overwritten in cloudbuild.yaml)
resource "google_cloud_run_service" "run_service" {
  name = local.package_rater_app_cloud_run_name
  location = "us-central1"

  template {
    spec {
      containers {
        # image = "us-docker.pkg.dev/cloudrun/container/placeholder:latest" # Placeholder
        image = "us-central1-docker.pkg.dev/${var.project_id}/${local.artifact_registry_repo_name}/${local.package_rater_app_image_name}:latest"
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

      timeout_seconds = 90
      container_concurrency = 5
      service_account_name = google_service_account.package_rater_service_account.email
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "10"
      }
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

# SQL Database
resource "google_sql_database_instance" "mysql-instance" {
  name             = local.google_sql_database_instance
  region           = local.region
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

# SQL users
resource "google_sql_user" "read-user" {
  name     = "read-user"
  instance = local.google_sql_database_instance
  host     = "%"
  password = "strongpassword"
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