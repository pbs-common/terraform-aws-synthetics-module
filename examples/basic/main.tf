locals {
  rendered_file_content = templatefile("${path.module}/canary.js.tpl", {
    name            = var.product,
    take_screenshot = true,
    urls = jsonencode([
      "https://www.pbs.org",
    ]),
    region = data.aws_region.current.region
  })

  zip_file = "lambda-canary-${sha256(local.rendered_file_content)}.zip"
}

data "archive_file" "lambda_canary_zip" {
  type        = "zip"
  output_path = local.zip_file
  source {
    content  = local.rendered_file_content
    filename = "nodejs/node_modules/canary.js"
  }
}

module "synthetics" {
  source = "../.."

  zip_file = data.archive_file.lambda_canary_zip.output_path

  # Makes the bucket easier to cleanup in testing.
  # Don't do this in production.
  force_destroy = true

  organization = var.organization
  environment  = var.environment
  product      = var.product
  repo         = var.repo
  owner        = var.owner
}
