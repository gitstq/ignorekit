# Contributing to ignorekit

Thanks for taking the time to improve **ignorekit**! Every kind of help is welcome:
bug reports, feature requests, documentation improvements, and pull requests.

## Quick setup

```bash
git clone https://github.com/gitstq/ignorekit.git
cd ignorekit
go test ./...     # run the test suite
go build .        # build the binary
```

ignorekit intentionally uses **only the Go standard library**. Please do not add
third-party module dependencies unless there is a compelling reason discussed
in an issue first.

## Reporting bugs

Open an issue and include:

- The output of `ignorekit version`.
- Your operating system and the command you ran.
- A minimal reproduction (a small fixture directory is ideal).
- The expected versus actual behavior.

## Pull request process

1. Fork the repository and create a topic branch from `main`.
2. Follow the existing code style and add comments for exported identifiers.
3. Add tests for new behavior and make sure `go test ./...` passes.
4. Run `gofmt -s -w .` and `go vet ./...`.
5. Use [Conventional Commits](https://www.conventionalcommits.org/), for example:
   - `feat: add Bun lockfile detection`
   - `fix: correct middle double-star matching`
   - `docs: clarify MCP template guidance`
6. Open the pull request and describe what changed and why.

## Adding a detection signature / template

- Add a `Signature` in `internal/detect/signatures.go` with high-confidence markers.
- Add the matching `Template` in `internal/templates/templates.go`.
- Keep markers specific to avoid false positives; distinguish files that should be
  committed (`Shared`) from those that should be ignored.
- Add a test fixture in the relevant `*_test.go`.

## Code of conduct

Please be respectful and constructive in all interactions. We want this to be a
welcoming space for contributors of every background.

## License

By contributing, you agree that your contributions will be licensed under the
[MIT License](LICENSE).
