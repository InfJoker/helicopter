data "yandex_compute_image" "coi" {
  family = "container-optimized-image"
}

resource "yandex_compute_instance" "server" {
  name               = "server"
  platform_id        = "standard-v2"
  service_account_id = yandex_iam_service_account.helicopter_node_sa.id
  boot_disk {
    initialize_params {
      image_id = data.yandex_compute_image.coi.id
      size     = 30
    }
  }
  network_interface {
    subnet_id = yandex_vpc_subnet.helicopter-subnet-1.id
    nat       = true
  }
  resources {
    cores         = 2
    memory        = 1
    core_fraction = 20
  }
  scheduling_policy {
    preemptible = true
  }
  metadata = {
    docker-compose = templatefile("${path.module}/files/spec.yaml", {
      openai_api_key = var.openai_api_key,
      registry_id    = var.registry_id
    })
    user-data = file("${path.module}/files/cloud_config.yaml")
  }
}
