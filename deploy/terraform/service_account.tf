resource "yandex_iam_service_account" "helicopter_node_sa" {
  name        = "helicopter-node-sa"
  description = "service account to pull images from the registry"
}

resource "yandex_resourcemanager_folder_iam_binding" "coi_puller" {
  role      = "container-registry.images.puller"
  folder_id = var.folder_id
  members = [
    "serviceAccount:${yandex_iam_service_account.helicopter_node_sa.id}",
  ]
}
