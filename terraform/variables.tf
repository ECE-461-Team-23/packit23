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

variable "jwt_secret" {
  description = "Secret value for signing JWTs"
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

variable "user_logins" {
  description = "JSON object holding logins to fill in user table"
  type = map(string)
  sensitive = true
}