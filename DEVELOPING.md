# Developing `vidx2pidx`

Follow these steps to set up and start developing for `vidx2pidx`:

## Prerequisites

Ensure you have the following installed:

- [GNU Make](https://www.gnu.org/software/make/)
- [Golang](https://golang.org/doc/install)
- [GolangCI-Lint](https://golangci-lint.run/usage/install/#local-installation)

## Setup

1. Clone the repository:

   ```sh
   git clone https://github.com/open-cmsis-pack/vidx2pidx.git
   ```

2. Navigate into the project directory:

   ```sh
   cd vidx2pidx
   ```

3. Configure your local environment:

   ```sh
   make config
   ```

4. Run tests to verify everything is working:

   ```sh
   make test-all
   ```

5. Build the project:

   ```sh
   make build/vidx2pidx
   ```

6. You're all set! 🎉 Start modifying the source code and refer to the [Contributing Guide](CONTRIBUTING.md)
for guidelines on contributing.

## Releasing

If you have push access to the `main` branch, you can create a new release by running:

```sh
make release
```

📌 **Note:** We follow [Semantic Versioning](https://semver.org/) for versioning `vidx2pidx`.
