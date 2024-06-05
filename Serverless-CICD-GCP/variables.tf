variable "Project" {
    description = "the project ID where all resources will be launched "
    type = string
}

variable "Location" {
  description = "The location (region or zone) to deploy the cloud Run Services."
  type = string
}

variable "gcr_region" {
  description = "name of teh GCP region where the GCR registry is located . e.g 'US', 'UE'"
  type = string
}

// OPTIONAL PARAMETERS


variable "service_name" {
  description = "The name of the Cloud Run service to deploy."
  type        = string
  default     = "sample-docker-service"
}

variable "repository_name" {
  description = "Name of the Google Cloud Source Repository to create."
  type        = string
  default     = "sample-docker-app"
}

variable "image_name" {
  description = "The name of the image to deploy. Defaults to a publically available image."
  type        = string
  default     = "gcr.io/cloudrun/hello"
}

variable "branch_name" {
  description = "Example branch name used to trigger builds."
  type        = string
  default     = "master"
}

variable "digest" {
  description = "The docker image digest or tag to deploy."
  type        = string
  default     = "latest"
}

variable "deploy_db" {
  description = "Whether to deploy a Cloud SQL database or not."
  type        = bool
  default     = false
}

variable "db_instance_name" {
  description = "The name of the Cloud SQL database instance."
  type        = string
  default     = "master-mysql-instance"
}

variable "db_name" {
  description = "The name of the Cloud SQL database."
  type        = string
  default     = "exampledb"
}

variable "db_username" {
  description = "The name of the Cloud SQL database user."
  type        = string
  default     = "testuser"
}

variable "db_password" {
  description = "The password of the Cloud SQL database user."
  type        = string
  default     = "testpassword"
}

variable "db_user_host" {
  description = "The host of the Cloud SQL database user. Used by MySQL."
  type        = string
  default     = "%"
}