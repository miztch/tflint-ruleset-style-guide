variable "config" {
  type        = any
  description = "A config map"
}

variable "items" {
  type        = list(any)
  description = "A list of items"
}

variable "record" {
  type = object({
    name  = string
    value = any
  })
  description = "A record object"
}
