variable "yc_cloud" {
  type = string
  description = "Yandex Cloud ID"
}

variable "yc_folder" {
  type = string
  description = "Yandex Cloud folder"
}

variable "yc_token" {
  type = string
  description = "Yandex Cloud OAuth token"
}

variable "db_password" {
  description = "PostgreSQL user pasword"
}

variable "yc_file" {
  type = string
  description = "Filename for Yandex service account key"
}

variable "yc_host_number" {
  description = "Used hosts number"
  type = number
}

