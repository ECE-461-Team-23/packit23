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
                google_secret_manager_secret_iam_member.package_rater_access] # Make sure service account is attached to policy to give access to secret token
}

# Run containers for read apis
resource "google_cloud_run_service" "read_apis_run_service" {
  name = local.package_rater_app_cloud_run_name
  location = local.region

  template {
    spec {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/placeholder:latest" # Placeholder
        # image = "us-central1-docker.pkg.dev/${var.project_id}/${local.artifact_registry_repo_name}/${local.package_rater_app_image_name}:latest"
        env {
          name = "PROJECT_ID"
          value = var.project_id
        }
        env {
          name = "REGION"
          value = local.region
        }
        env {
          name = "INSTANCE_NAME"
          value = local.mysql_db_instance_name
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
          name = "DB_PASSWORD"
          value_from {
            secret_key_ref {
              name = "READ_USER_PASSWORD"
              key  = "latest"
            }
          }
        }
      }

      timeout_seconds = 90
      container_concurrency = 5
      service_account_name = google_service_account.read_apis_service_account.email
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "10"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [google_project_service.cloud_run_api,
                google_secret_manager_secret_iam_member.read_apis_access]
}

# Run containers for write apis
resource "google_cloud_run_service" "write_apis_run_service" {
  name = local.write_apis_app_cloud_run_name
  location = local.region

  template {
    spec {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/placeholder:latest" # Placeholder
        # image = "us-central1-docker.pkg.dev/${var.project_id}/${local.artifact_registry_repo_name}/${local.package_rater_app_image_name}:latest"
        env {
          name = "PROJECT_ID"
          value = var.project_id
        }
        env {
          name = "REGION"
          value = local.region
        }
        env {
          name = "INSTANCE_NAME"
          value = local.mysql_db_instance_name
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
          name = "DB_PASSWORD"
          value_from {
            secret_key_ref {
              name = "WRITE_USER_PASSWORD"
              key  = "latest"
            }
          }
        }
      }

      timeout_seconds = 90
      container_concurrency = 5
      service_account_name = google_service_account.read_apis_service_account.email
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "10"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [google_project_service.cloud_run_api,
                google_secret_manager_secret_iam_member.write_apis_access]
}

# Run cloud run service from api-spec template file
data "template_file" "cloud_run" {
  template = "${file("${path.module}/api_spec.tpl")}"
  vars = {
    read_url = "${google_cloud_run_service.run_service.status[1].url}"
    write_url = "${google_cloud_run_service.run_service.status[2].url}"
  }
}