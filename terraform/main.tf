# Settings
locals {
  # General
  github_branch = "terraform"
  artifact_registry_repo_name = "container-repo"
  region = "us-central1"

  # package-rater-app
  package_rater_app_cloud_run_name = "package-rater-app"
  package_rater_app_image_name = "package-rater-image"

  # read-apis-app
  read_db_user_name = "read-user"
  read_apis_app_cloud_run_name = "read-apis-app"
  read_apis_app_image_name = "read-apis-image"

  # write-apis-app
  write_db_user_name = "write-user"
  write_apis_app_cloud_run_name = "write-apis-app"
  write_apis_app_image_name = "write-apis-image"

  # SQL
  mysql_db_name = "mysql-db"
  mysql_db_instance_name = "mysql-instance"
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

# Automatically build containers
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

resource "google_cloudbuild_trigger" "read_apis_app_trigger" {
  location = "us-central1"

  github {
    owner = "packit461"
    name = "packit23"

    push {
      branch = local.github_branch
    }
  }

  filename = "read_apis/cloudbuild.yaml"
  depends_on = [google_project_service.cloud_build_api]
}

resource "google_cloudbuild_trigger" "write_apis_app_trigger" {
  location = "us-central1"

  github {
    owner = "packit461"
    name = "packit23"

    push {
      branch = local.github_branch
    }
  }

  filename = "delete_write_apis/cloudbuild.yaml"
  depends_on = [google_project_service.cloud_build_api]
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

resource "google_secret_manager_secret" "read_user_password_manager" {
  secret_id = "READ_USER_PASSWORD"

  replication {
    automatic = true
  }

  depends_on = [ google_project_service.secret_manager_api ]
}

resource "google_secret_manager_secret" "write_user_password_manager" {
  secret_id = "WRITE_USER_PASSWORD"

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

resource "google_secret_manager_secret_version" "read_user_password_secret" {
  secret   = google_secret_manager_secret.read_user_password_manager.id
  secret_data = var.read_user_password
}

resource "google_secret_manager_secret_version" "write_user_password_secret" {
  secret   = google_secret_manager_secret.write_user_password_manager.id
  secret_data = var.write_user_password
}

# Give Package Rater SA access to "Github Token" secret
resource "google_secret_manager_secret_iam_member" "package_rater_access" {
  secret_id = google_secret_manager_secret.github_token_manager.secret_id
  role = "roles/secretmanager.secretAccessor"
  member = "serviceAccount:${google_service_account.package_rater_service_account.email}"
}

resource "google_secret_manager_secret_iam_member" "read_apis_access" {
  secret_id = google_secret_manager_secret.read_user_password_manager.secret_id
  role = "roles/secretmanager.secretAccessor"
  member = "serviceAccount:${google_service_account.read_apis_service_account.email}"
}

resource "google_secret_manager_secret_iam_member" "write_apis_access" {
  secret_id = google_secret_manager_secret.write_user_password_manager.secret_id
  role = "roles/secretmanager.secretAccessor"
  member = "serviceAccount:${google_service_account.write_apis_service_account.email}"
}

output "service_url" {
  value = google_cloud_run_service.run_service.status[0].url
}

# SQL Database
resource "google_sql_database_instance" "mysql-instance" {
  name             = local.mysql_db_instance_name
  region           = local.region
  database_version = "MYSQL_8_0"
  settings {
    tier = "db-f1-micro"
  }

  deletion_protection  = "true"
}

resource "google_sql_database" "database" {
  name = local.mysql_db_name
  instance = local.mysql_db_instance_name
}

resource "random_id" "db_name_suffix" {
  byte_length = 4
}

# SQL users
resource "google_sql_user" "read-user" {
  name     = local.read_user_name
  instance = local.mysql_db_instance_name
  host     = "%"
  password = var.read_user_password
}

resource "google_sql_user" "write-user" {
  name     = local.write_user_name
  instance = local.mysql_db_instance_name
  host     = "%"
  password = var.write_user_password
}

# API Gateway
resource "google_api_gateway_api" "api_gw" {
  project = var.project_id
  provider = google-beta
  api_id = "my-api"
}

resource "google_api_gateway_api_config" "api_cfg" {
  project = var.project_id
  provider = google-beta
  api = google_api_gateway_api.api_gw.api_id
  api_config_id = "api-config"

  openapi_documents {
    document {
      path = "api_spec.yaml"
      contents = base64encode("api_spec.yaml")
    }
  }
  lifecycle {
    create_before_destroy = true
  }
}

resource "google_api_gateway_gateway" "gw" {
  project = var.project_id
  provider = google-beta
  api_config = google_api_gateway_api_config.api_cfg.id
  gateway_id = "api_gw"

  depends_on = [google_api_gateway_api_config.api_cfg]

}

## Service Accounts ##

resource "google_service_account" "package_rater_service_account" {
  account_id = "package-rater-sa"
  display_name = "Service account for package rater containers"
}

resource "google_service_account" "write_apis_service_account" {
  account_id = "package-rater-sa"
  display_name = "Service account for delete/write apis containers"
}

resource "google_service_account" "read_apis_service_account" {
  account_id = "package-rater-sa"
  display_name = "Service account for read apis containers"
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

resource "google_project_service" "cloud_sql_api" {
  service = "sql-component.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "api_gateway_api" {
  service = "apigateway.googleapis.com"
  disable_on_destroy = true
}

resource "google_project_service" "service_control_api" {
  service = "servicecontrol.googleapis.com"
  disable_on_destroy = true
}

resource "google_project_service" "service_management_api" {
  service = "servicemanagement.googleapis.com"
  disable_on_destroy = true
}