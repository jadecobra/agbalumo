# Phase 11: Docker Litestream Optimization (Eliminating Compiler Bottleneck)

## Objective
Accelerate the CI/CD deployment pipeline by eliminating the massive "compile Litestream from source" Go build stage in the `Dockerfile`. We will replace it with a direct binary copy from the official lightweight Litestream image or static release.

## Context
Currently, the `litestream-builder` stage pulls the full Go compiler, clones github.com/benbjohnson/litestream, and downloads hundreds of heavy dependencies (`grpc`, `oauth2`, `otel`) just to compile a binary. This wastes incredibly valuable CI time and adds vulnerability surface area. 

## Steps for Execution
1. Open the `Dockerfile`.
2. Delete the entire `litestream-builder` stage:
   ```dockerfile
   # DELETE THIS:
   FROM golang:1.25-bookworm AS litestream-builder
   ... all the way down to ...
   CGO_ENABLED=1 go install -ldflags '-extldflags "-static"' ./cmd/litestream
   ```
3. Modify the final `Runtime Stage` copy step. Instead of copying from `litestream-builder`, we natively copy the static binary from the official Litestream Docker image (which requires absolutely zero build time):
   ```dockerfile
   # Replace the old copy command with this:
   # (Note: Use the appropriate litestream version tag, e.g., latest or 0.3.13)
   COPY --from=litestream/litestream:latest /usr/local/bin/litestream /usr/local/bin/litestream
   ```
   *(Alternative fallback if the user requires a highly specific custom fork `v0.5.10` not published on Docker Hub: Create a minimal alpine downloader stage that runs `wget -O litestream.tar.gz https://github.com/.../releases/download/...` instead of compiling the source).*

4. Run `docker build -t agbalumo-test .` to verify that the build time is now drastically cut down.
5. Commit natively: `perf(infra): drastically reduce docker build time by using pre-compiled litestream binary`.

## Verification
- Docker build completes without pulling the `go-jose`, `grpc`, or `crypto` libraries for Litestream.
- Deployment pipeline latency is reduced by several minutes.
