terraform {
  required_version = "1.7.5"
}

provider "google" {
  project     = "my-project-id"
  region      = "us-central1"
}

//Deploy a Google Cloud source repositry
resource "google_sourcerepo_repository" "repo" {
  name = var.repository_name
}

//Deploy Cloud run Services 
resource "google_cloud_run_service" "service" {
  name = var.service_name
  location = var.Location

  template {
    metadata {
      annotations = {
        "client.knative.dev/user-image" = local.image_name
      }
    }
    spec {
      containers {
        image = local.image_name
      }
    }
  }
  traffic {
    percent = 100
    latest_revision = true
  }
}

//Expose the services publicty 
resource "google_cloud_run_service_iam_member" "allUsers" {
  service  = google_cloud_run_service.service.name
  location = google_cloud_run_service.service.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

locals {
  image_name = var.image_name == "" ? "${var.gcr_region}.gcr.io/${var.Project}/${var.service_name}" : var.image_name
  # uncomment the following line to connect to the cloud sql database instance
  #instance_connection_name = "${var.project}:${var.location}:${google_sql_database_instance.master[0].name}"
}
// Create a Cloud build trigger 
resource "google_cloudbuild_trigger" "cloud_build_trigger" {
  description = "Cloud Source Repository Trigger ${var.repository_name} (${var.branch_name})"
  trigger_template {
    branch_name = var.branch_name
    repo_name   = var.repository_name
  }///
  substitutions = {
    _LOCATION     = var.Location
    _GCR_REGION   = var.gcr_region
    _SERVICE_NAME = var.service_name
  }
 # The filename argument instructs Cloud Build to look for a file in the root of the repository.
 filename = "cloudbuild.yaml"
 depends_on = [google_sourcerepo_repository.repo]

}

///filename = "cloudbuild.yaml"

resource "google_sql_database_instance" "master" {
  count            = var.deploy_db ? 1 : 0
  name             = var.db_instance_name
  region           = var.Location
  database_version = "MYSQL_5_7"

  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_database" "default" {
  count = var.deploy_db ? 1 : 0

  name     = var.db_name
  project  = var.Project
  instance = google_sql_database_instance.master[0].name

  depends_on = [google_sql_database_instance.master]
}

resource "google_sql_user" "default" {
  count = var.deploy_db ? 1 : 0

  project  = var.Project
  name     = var.db_username
  instance = google_sql_database_instance.master[0].name

  host     = var.db_user_host
  password = var.db_password

  depends_on = [google_sql_database.default]
}

