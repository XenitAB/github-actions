variable "location" {
  description = "The Azure Region in which all resources will be created."
  type        = string
  validation {
    condition     = contains([
      "northeurope",
      "westeurope",
      "swedencentral"
      ], var.location)
    error_message = "The location must be one of 'northeurope', 'westeurope', or 'swedencentral'."
  }
}

variable "subscription_id" {
  description = "SubscriptionId of context"
  type        = string
  default     = "935666ec-670a-4300-95bb-6dbe86bb61f7"
}

variable "environment" {
  description = "The environment in which the resources will be created."
  type        = string
  }

variable "suffix" {
  description = "The suffix to append to the resource names."
  type        = string
}

variable "resource_locks" {
  description = "Enable or disable resource locks for critical resources."
  type        = bool
  default     = false
}