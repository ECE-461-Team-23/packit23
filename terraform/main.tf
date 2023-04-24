# Settings
locals {
  # General
  github_branch = "main"
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
  mysql_db_name = "appdata"
  mysql_db_instance_name = "mysql-instance"

  # Bucket
  bucket_db_name = "binaries"
}

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.61.0"
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

## Enable services ##

resource "google_project_service" "cloud_run_api" {
  service = "run.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "iam_api" {
  service = "iam.googleapis.com"
  disable_on_destroy = false
}
    
resource "google_project_service" "cloud_resource_manager_api" {
  service = "cloudresourcemanager.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "cloud_build_api" {
  service = "cloudbuild.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "artifact_registry_api" {
  service = "artifactregistry.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "secret_manager_api" {
  service = "secretmanager.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "cloud_sql_api" {
  service = "sql-component.googleapis.com"
  disable_on_destroy = false # Need to stay false
}

resource "google_project_service" "api_gateway_api" {
  service = "apigateway.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "service_control_api" {
  service = "servicecontrol.googleapis.com"
  disable_on_destroy = false # Need to stay false
}

resource "google_project_service" "service_management_api" {
  service = "servicemanagement.googleapis.com"
  disable_on_destroy = false # Need to stay false
}

resource "google_project_service" "sql_admin_api" {
  service = "sqladmin.googleapis.com"
  disable_on_destroy = false
}

# Run containers for package-rater-app (container image is overwritten in cloudbuild.yaml)
resource "google_cloud_run_service" "package_rater_run_service" {
  name = local.package_rater_app_cloud_run_name
  location = local.region

  template {
    spec {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/placeholder:latest" # Placeholder
        # image = "us-central1-docker.pkg.dev/${var.project_id}/${local.artifact_registry_repo_name}/${local.package_rater_app_image_name}:latest"
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
                google_secret_manager_secret_iam_member.package_rater_token_access] # Make sure service account is attached to policy to give access to secret token
}

# Run containers for read apis
resource "google_cloud_run_service" "read_apis_run_service" {
  name = local.read_apis_app_cloud_run_name
  location = local.region

  template {
    spec {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/placeholder:latest" # Placeholder
        # image = "us-central1-docker.pkg.dev/${var.project_id}/${local.artifact_registry_repo_name}/${local.read_apis_app_image_name}:latest"
        env {
          name = "PROJECT_ID"
          value = var.project_id
        }
        env {
          name = "REGION"
          value = local.region
        }
        env {
          name = "PACKAGE_RATER_URL"
          value = google_cloud_run_service.package_rater_run_service.status[0].url
        }
        env {
          name = "BUCKET_NAME"
          value = google_storage_bucket.binary_bucket.name
        }
        env {
          name = "INSTANCE_CONNECTION_NAME"
          value = google_sql_database_instance.mysql_instance.connection_name
        }
        env {
          name = "DB_NAME"
          value = local.mysql_db_name
        }
        env {
          name = "DB_USER"
          value = local.read_db_user_name
        }
        env {
          name = "DB_PASS"
          value_from {
            secret_key_ref {
              name = "READ_USER_PASSWORD"
              key  = "latest"
            }
          }
        }
        env {
          name = "JWT_SECRET"
          value_from {
            secret_key_ref {
              name = "JWT_SECRET"
              key  = "latest"
            }
          }
        }
      }

      timeout_seconds = 300
      container_concurrency = 1
      service_account_name = google_service_account.read_apis_service_account.email
    }

    metadata {
      annotations = {
        "run.googleapis.com/client-name" = "terraform"
        "run.googleapis.com/cloudsql-instances" = google_sql_database_instance.mysql_instance.connection_name
        "autoscaling.knative.dev/maxScale" = "20"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [google_project_service.cloud_run_api,
                google_secret_manager_secret_iam_member.read_apis_access,
                google_cloud_run_service.package_rater_run_service,
                google_storage_bucket_iam_member.read_apis_bucket_access]
}

# Run containers for write apis
resource "google_cloud_run_service" "write_apis_run_service" {
  name = local.write_apis_app_cloud_run_name
  location = local.region

  template {
    spec {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/placeholder:latest" # Placeholder
        # image = "us-central1-docker.pkg.dev/${var.project_id}/${local.artifact_registry_repo_name}/${local.write_apis_app_image_name}:latest"
        env {
          name = "PROJECT_ID"
          value = var.project_id
        }
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
          name = "REGION"
          value = local.region
        }
        env {
          name = "PACKAGE_RATER_URL"
          value = google_cloud_run_service.package_rater_run_service.status[0].url
        }
        env {
          name = "BUCKET_NAME"
          value = google_storage_bucket.binary_bucket.name
        }
        env {
          name = "INSTANCE_CONNECTION_NAME"
          value = google_sql_database_instance.mysql_instance.connection_name
        }
        env {
          name = "DB_NAME"
          value = local.mysql_db_name
        }
        env {
          name = "DB_USER"
          value = local.write_db_user_name
        }
        env {
          name = "DB_PASS"
          value_from {
            secret_key_ref {
              name = "WRITE_USER_PASSWORD"
              key  = "latest"
            }
          }
        }
        env {
          name = "JWT_SECRET"
          value_from {
            secret_key_ref {
              name = "JWT_SECRET"
              key  = "latest"
            }
          }
        }
        env {
          name = "USER_LOGINS"
          value_from {
            secret_key_ref {
              name = "USER_LOGINS"
              key  = "latest"
            }
          }
        }
      }

      timeout_seconds = 300
      container_concurrency = 1
      service_account_name = google_service_account.write_apis_service_account.email
    }

    metadata {
      annotations = {
        "run.googleapis.com/client-name" = "terraform"
        "run.googleapis.com/cloudsql-instances" = google_sql_database_instance.mysql_instance.connection_name
        "autoscaling.knative.dev/maxScale" = "20"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [google_project_service.cloud_run_api,
                google_secret_manager_secret_iam_member.write_apis_access,
                google_cloud_run_service.package_rater_run_service,
                google_secret_manager_secret_iam_member.write_apis_token_access,
                google_storage_bucket_iam_member.write_apis_bucket_access]
}

# Automatically build containers
resource "google_cloudbuild_trigger" "package_rater_app_trigger" {
  location = "us-central1"
  name     = "package-rater-app-trigger"

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
  name     = "read-apis-app-trigger"

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
  name     = "write-apis-app-trigger"

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
resource "google_cloud_run_service_iam_member" "package_rater_users" {
  service  = google_cloud_run_service.package_rater_run_service.name
  location = google_cloud_run_service.package_rater_run_service.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

resource "google_cloud_run_service_iam_member" "read_apis_users" {
  service  = google_cloud_run_service.read_apis_run_service.name
  location = google_cloud_run_service.read_apis_run_service.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

resource "google_cloud_run_service_iam_member" "write_apis_users" {
  service  = google_cloud_run_service.write_apis_run_service.name
  location = google_cloud_run_service.write_apis_run_service.location
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

resource "google_secret_manager_secret" "jwt_secret_manager" {
  secret_id = "JWT_SECRET"

  replication {
    automatic = true
  }

  depends_on = [ google_project_service.secret_manager_api ]
}

resource "google_secret_manager_secret" "user_logins_manager" {
  secret_id = "USER_LOGINS"

  replication {
    automatic = true
  }

  depends_on = [ google_project_service.secret_manager_api ]
}

# Create a new version of "Github Token" secret
resource "google_secret_manager_secret_version" "github_token_manager_version" {
  secret   = google_secret_manager_secret.github_token_manager.id
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

resource "google_secret_manager_secret_version" "jwt_secret" {
  secret   = google_secret_manager_secret.jwt_secret_manager.id
  secret_data = var.jwt_secret
}

resource "google_secret_manager_secret_version" "user_logins_secret" {
  secret   = google_secret_manager_secret.user_logins_manager.id
  secret_data = jsonencode(var.user_logins)
}

# Give service accounts access to "Github Token" secret
resource "google_secret_manager_secret_iam_member" "package_rater_token_access" {
  secret_id = google_secret_manager_secret.github_token_manager.secret_id
  role = "roles/secretmanager.secretAccessor"
  member = "serviceAccount:${google_service_account.package_rater_service_account.email}"
}

resource "google_secret_manager_secret_iam_member" "write_apis_token_access" {
  secret_id = google_secret_manager_secret.github_token_manager.secret_id
  role = "roles/secretmanager.secretAccessor"
  member = "serviceAccount:${google_service_account.write_apis_service_account.email}"
}

# Give access to passwords for DB
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

# Give access to JWT secret
resource "google_secret_manager_secret_iam_member" "read_apis_jwt_access" {
  secret_id = google_secret_manager_secret.jwt_secret_manager.secret_id
  role = "roles/secretmanager.secretAccessor"
  member = "serviceAccount:${google_service_account.read_apis_service_account.email}"
}

resource "google_secret_manager_secret_iam_member" "write_apis_jwt_access" {
  secret_id = google_secret_manager_secret.jwt_secret_manager.secret_id
  role = "roles/secretmanager.secretAccessor"
  member = "serviceAccount:${google_service_account.write_apis_service_account.email}"
}

# Give access to user logins
resource "google_secret_manager_secret_iam_member" "write_apis_user_logins_access" {
  secret_id = google_secret_manager_secret.user_logins_manager.secret_id
  role = "roles/secretmanager.secretAccessor"
  member = "serviceAccount:${google_service_account.write_apis_service_account.email}"
}

# Give service accounts database access
resource "google_project_iam_member" "write_apis_db_access" {
  project = var.project_id
  role = "roles/cloudsql.editor"
  member = "serviceAccount:${google_service_account.write_apis_service_account.email}"
}

resource "google_project_iam_member" "read_apis_db_access" {
  project = var.project_id
  role = "roles/cloudsql.client"
  member = "serviceAccount:${google_service_account.read_apis_service_account.email}"
}

resource "google_project_iam_member" "write_apis_db_user" {  
 project = var.project_id
 role   = "roles/cloudsql.instanceUser"  
 member = "serviceAccount:${google_service_account.write_apis_service_account.email}"
}  

resource "google_project_iam_member" "read_apis_db_user" {  
 project = var.project_id
 role   = "roles/cloudsql.instanceUser"  
 member = "serviceAccount:${google_service_account.read_apis_service_account.email}"
}  

# Bucket Access
resource "google_storage_bucket_iam_member" "write_apis_bucket_access" {  
 bucket = google_storage_bucket.binary_bucket.name
 role   = "roles/storage.objectAdmin"  
 member = "serviceAccount:${google_service_account.write_apis_service_account.email}"
}

resource "google_storage_bucket_iam_member" "read_apis_bucket_access" {  
 bucket = google_storage_bucket.binary_bucket.name
 role   = "roles/storage.objectViewer"  
 member = "serviceAccount:${google_service_account.read_apis_service_account.email}"
}

# Outputs
output "package_rater_service_url" {
  value = google_cloud_run_service.package_rater_run_service.status[0].url
  description = "url for package rater service"
}

output "read_apis_service_url" {
  value = google_cloud_run_service.read_apis_run_service.status[0].url
  description = "url for read apis service"
}

output "write_apis_service_url" {
  value = google_cloud_run_service.write_apis_run_service.status[0].url
  description = "url for write apis service"
}

# Bucket
resource "google_storage_bucket" "binary_bucket" {
  name          = "binaries-${random_id.bucket_name_suffix.hex}"
  location      = "US"
  force_destroy = true

  uniform_bucket_level_access = true
  public_access_prevention = "enforced"

  versioning {
    enabled = true
  }

  lifecycle_rule {
    condition {
      num_newer_versions = 3
    }
    action {
      type = "Delete"
    }
  }
}

resource "random_id" "bucket_name_suffix" {
  byte_length = 8
}

# SQL Database
resource "google_sql_database_instance" "mysql_instance" {
  name             = local.mysql_db_instance_name
  region           = local.region
  database_version = "MYSQL_8_0"
  settings {
    tier = "db-f1-micro"
    disk_size = 100 # GB

    database_flags {
      name  = "cloudsql_iam_authentication" #"cloudsql.iam_authentication"
      value = "on"
    }

    # database_flags {
    #   name  = "cloudsql.iam_authentication"
    #   value = "on"
    # }
  }

  deletion_protection  = "true"
}

resource "google_sql_database" "database" {
  name = local.mysql_db_name
  instance = google_sql_database_instance.mysql_instance.name
}

# SQL users
resource "google_sql_user" "read-user" {
  name     = local.read_db_user_name
  instance = google_sql_database_instance.mysql_instance.name
  # type = "CLOUD_IAM_USER"
  host     = "%"
  password = var.read_user_password
}

resource "google_sql_user" "write-user" {
  name     = local.write_db_user_name
  instance = google_sql_database_instance.mysql_instance.name
  # type = "CLOUD_IAM_USER"
  host     = "%"
  password = var.write_user_password
}

# API Gateway - Create by hand for final project
# resource "google_api_gateway_api" "api_gw" {
#   project = var.project_id
#   provider = google-beta
#   api_id = "my-api"

#   depends_on = [ google_project_service.api_gateway_api ]
# }

# resource "google_api_gateway_api_config" "api_cfg" {
#   project = var.project_id
#   provider = google-beta
#   api = google_api_gateway_api.api_gw.api_id
#   api_config_id = "api-config"

#   openapi_documents {
#     document {
#       path = "api_spec.yaml"
#       contents = base64encode(templatefile("${path.module}/api_spec.yaml", { read_url=google_cloud_run_service.read_apis_run_service.status[0].url, write_url=google_cloud_run_service.write_apis_run_service.status[0].url }))
#     }
#   }
#   lifecycle {
#     create_before_destroy = true
#   }
# }

# resource "google_api_gateway_gateway" "gateway" {
#   project = var.project_id
#   region = local.region
#   provider = google-beta
#   api_config = google_api_gateway_api_config.api_cfg.id
#   gateway_id = "api-gw"

#   depends_on = [google_api_gateway_api_config.api_cfg]

# }

## Service Accounts ##

resource "google_service_account" "package_rater_service_account" {
  account_id = "package-rater-sa"
  display_name = "Service account for package rater containers"
}

resource "google_service_account" "write_apis_service_account" {
  account_id = "write-apis-sa"
  display_name = "Service account for delete/write apis containers"
}

resource "google_service_account" "read_apis_service_account" {
  account_id = "read-apis-sa"
  display_name = "Service account for read apis containers"
}
