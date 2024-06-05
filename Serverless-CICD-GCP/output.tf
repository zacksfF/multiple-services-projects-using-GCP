output "url" {
  description = "The URL where the cloud Run services can be access"
  value = google_cloud_run_service.service.status[0].url
}

output "trigger_id" {
  description = "The unique ideantifier for the cloud build trigger"
  value = google_cloudbuild_trigger.google_cloudbuild_trigger.trigger_id
}

output "repositry_http_url" {
  description = "HTTTP URL of the repository in CLoud Source Repositories."
  value = google_sourcerepo_repository.repo.url
}