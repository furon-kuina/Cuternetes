terraform {
  required_version = ">= 1.6.6" 
  
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "5.10.0"
    }
  }
}

provider "google" {
  credentials = file("./secrets/playground-408516-7e0853935baa.json")

  project = "playground-408516"
  region = "asia-northeast1"
  zone = "asia-northeaset1-a"
}

resource "google_service_account" "default" {
  account_id   = "compute-sa"
  display_name = "Custom SA for VM Instance"
}