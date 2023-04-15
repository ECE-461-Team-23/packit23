variable "project_id"{
  description = "GCP Project ID"
  type        = string
  sensitive   = false
}

variable "github_token" {
  description = "Github API Token"
  type        = string
  sensitive   = true
}

variable "read_user_password" {
  description = "Read User Password"
  type = string
  sensitive = true
}

variable "write_user_password" {
  description = "Write User Password"
  type = string
  sensitive = true
}