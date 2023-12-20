resource "google_compute_network" "vpc_network" {
  name = "c8s-network"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "c8s_subnet" {
  name          = "c8s"
  network       = google_compute_network.vpc_network.name
  ip_cidr_range = "10.240.0.0/24"
}

resource "google_compute_firewall" "allow_internal" {
  name    = "c8s-allow-internal"
  network = google_compute_network.vpc_network.name

  allow {
    protocol = "tcp"
  }

  allow {
    protocol = "udp"
  }

  allow {
    protocol = "icmp"
  }

  source_ranges = ["10.240.0.0/24", "10.200.0.0/16"]
}

resource "google_compute_firewall" "allow_external" {
  name    = "c8s-allow-external"
  network = google_compute_network.vpc_network.name

  allow {
    protocol = "tcp"
    ports    = ["22", "6443"]
  }

  allow {
    protocol = "icmp"
  }

  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_address" "c8s_address" {
  name   = "c8s"
  region = "asia-northeast1" # Replace with your region
}
