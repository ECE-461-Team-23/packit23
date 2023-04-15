variable "project_id" {
  description = "GCP Project ID"
  type        = string
  sensitive   = true
}

variable "jwt_secret" {
  description = "Secret value for signing JWTs"
  type        = string
  sensitive   = true
}