variable "project_id" {
  description = "GCP Project ID"
  type        = string
  sensitive   = true
}

variable "github_token" {
  description = "Github API Token"
  type        = string
  sensitive   = true
}