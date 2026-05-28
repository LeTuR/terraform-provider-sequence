# terraform-provider-sequence

A Terraform provider that generates **zero-padded sequential numbers**
(`"001"`, `"002"`, …) gated by a trigger. Conceptually a relative of
[`random_id`](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/id),
but producing stable, monotonic values instead of random ones.

## Usage

```hcl
terraform {
  required_providers {
    sequence = {
      source  = "LeTuR/sequence"
      version = "~> 0.1"
    }
  }
}

resource "sequence_number" "vm" {
  start  = 1       # first value, default 1
  width  = 3       # zero-padding width, default 3 → "001"
  prefix = "vm-"   # optional
  suffix = ""      # optional

  keepers = {
    region = var.region
  }
}

output "vm_name" {
  value = sequence_number.vm.formatted   # "vm-001"
}
```

Changing any value in `keepers` increments `number` by 1 on the next
`terraform apply`. The counter is per-resource and lives entirely in
Terraform state — no external backend.

## Attributes

| Attribute   | Type        | Required | Description                                     |
|-------------|-------------|----------|-------------------------------------------------|
| `start`     | number      | no       | First value on create (default `1`)             |
| `width`     | number      | no       | Zero-padding width (default `3`, `>= 0`)        |
| `prefix`    | string      | no       | Prepended to the padded number                  |
| `suffix`    | string      | no       | Appended to the padded number                   |
| `keepers`   | map(string) | no       | Trigger — any value change increments `number`  |
| `number`    | number      | computed | Current sequence value                          |
| `formatted` | string      | computed | `prefix + zero-pad(number, width) + suffix`     |
| `id`        | string      | computed | Same as `formatted`                             |

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.25 (for development only)

## Development

```bash
go build .              # build the provider binary
make generate           # regenerate docs/ from examples/
make testacc            # run acceptance tests (TF_ACC=1 go test ./...)
```

### Local provider override

To try the provider against real Terraform without publishing, add a
`dev_overrides` block to `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "LeTuR/sequence" = "/path/to/your/$GOPATH/bin"
  }
  direct {}
}
```

Then `terraform plan` in `examples/resources/sequence_number/` will use
your local build.

## Release process

Releases are **fully automated**. Push a conventional commit (`feat:`,
`fix:`, `perf:`) to `main`, and the `Release` workflow:

1. Runs `cog bump --auto --dry-run` to compute the next semver.
2. Creates a lightweight `vX.Y.Z` git tag (no commit on `main`).
3. Runs GoReleaser to build, sign, and publish a GitHub Release with
   binaries for `linux`, `darwin`, `windows`, and `freebsd` × `amd64`,
   `arm64`, `arm`, `386`.
4. The Terraform Registry picks up the new tag via webhook.

### One-time setup

1. Generate a GPG signing key:
   ```bash
   gpg --full-generate-key   # RSA 4096
   gpg --armor --export-secret-keys <key-id>   # for GitHub secret
   gpg --armor --export <key-id>               # for Terraform Registry
   ```
2. Add GitHub repo secrets:
   - `GPG_PRIVATE_KEY` — armored private key
   - `PASSPHRASE` — key passphrase
3. Add the public key to <https://registry.terraform.io> under
   **User Settings → Signing Keys**, then **Publish → Provider** and
   select `LeTuR/terraform-provider-sequence`.

## License

[MIT](LICENSE)
