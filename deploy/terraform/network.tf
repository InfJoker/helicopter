resource "yandex_vpc_network" "helicopter-network" {
  name = "helicopter-network"
}

resource "yandex_vpc_subnet" "helicopter-subnet-1" {
  name           = "helicopter-subnet-1"
  v4_cidr_blocks = ["10.2.0.0/16"]
  zone           = "ru-central1-a"
  network_id     = yandex_vpc_network.helicopter-network.id
}
