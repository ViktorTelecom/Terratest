output "load_balancer_public_ip" {
  description = "Public IP address of load balancer"
  value = yandex_lb_network_load_balancer.wp_lb.listener.*.external_address_spec[0].*.address
}

output "server_1_public_ip" {
  description = "Public IP address of server #1"
  value = yandex_compute_instance.wp-app-1.network_interface.0.nat_ip_address
}

output "server_2_public_ip" {
  description = "Public IP address of server #2"
  value = yandex_compute_instance.wp-app-2.network_interface.0.nat_ip_address
}