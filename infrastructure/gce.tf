resource "google_compute_instance" "controller" {
  count         = 3
  name          = "c8s-${count.index}"
  machine_type  = "e2-standard-2"
  zone          = "asia-northeast1-a" # Replace with your zone

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2004-lts"
      size  = 200
    }
  }

  network_interface {
    network    = google_compute_network.vpc_network.name
    subnetwork = google_compute_subnetwork.c8s_subnet.name

    access_config {
      // Ephemeral IP
    }
  }

  service_account {
    email  = google_service_account.default.email
    scopes = ["cloud-platform"]
  }

  tags = ["c8s"]
}