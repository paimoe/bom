terraform {
  required_providers {
    bunnynet = {
      source = "BunnyWay/bunnynet"
    }
  }
}

variable "github_read_token" {
  description = "Token used to read images for bom"
  type        = string
  sensitive   = true
}

# Image registry is already manually made
resource "bunnynet_compute_container_imageregistry" "github" {
  registry = "GitHub"
  username = "paimoe"
  token    = var.github_read_token
}

resource "bunnynet_compute_container_app" "app_bom" {
  name    = "bom"
  version = 2

  autoscaling_min = 1
  autoscaling_max = 5

  regions_allowed  = ["SYD"]
  regions_required = ["SYD"]

  container {
    name            = "app"
    image_registry  = bunnynet_compute_container_imageregistry.github.id
    image_namespace = "paimoe"
    image_name      = "bom"
    image_tag       = "latest"

    endpoint {
      name = "app-endpoint"
      type = "CDN"

      cdn {
        origin_ssl = false

        sticky_sessions {
          headers = ["X-Forwarded-For"]
        }
      }

      port {
        container = 80
      }
    }

    # env {
    #   name  = "APP_ENV"
    #   value = "prod"
    # }

    # env {
    #   name  = "LISTEN_PORT"
    #   value = "3000"
    # }
  }
}
