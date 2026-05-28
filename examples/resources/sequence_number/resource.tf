variable "region" {
  type    = string
  default = "eu-west-1"
}

resource "sequence_number" "vm" {
  start  = 1
  width  = 3
  prefix = "vm-"

  keepers = {
    region = var.region
  }
}

output "vm_name" {
  value = sequence_number.vm.formatted
}
