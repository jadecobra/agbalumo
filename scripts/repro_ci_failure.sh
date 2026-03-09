#!/bin/bash
# scripts/repro_ci_failure.sh
# Recreates the security failure seen in CI by building the Docker image
# and running a Trivy scan.

set -e

RED=$(printf '\033[0;31m')
GREEN=$(printf '\033[0;32m')
YELLOW=$(printf '\033[1;33m')
NC=$(printf '\033[0m')

IMAGE_NAME="agbalumo-repro"
TAG="latest"

echo "${YELLOW}Building Docker image $IMAGE_NAME:$TAG...${NC}"
docker build -t "$IMAGE_NAME:$TAG" .

echo ""
echo "${YELLOW}Running Trivy scan (via Docker) to match CI failure...${NC}"

# We use the official Trivy Docker image to ensure it works even if Trivy isn't installed locally.
# This matches the configuration in .github/workflows/ci.yml
# - exit-code: 1 (fail if vulnerabilities found)
# - ignore-unfixed: true
# - vuln-type: os,library
# - severity: CRITICAL,HIGH

docker run --rm \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v "$(pwd)/.cache/trivy:/root/.cache/" \
  aquasec/trivy:latest \
  image \
  --exit-code 1 \
  --ignore-unfixed \
  --vuln-type os,library \
  --severity CRITICAL,HIGH \
  "$IMAGE_NAME:$TAG"

RESULT=$?

if [ $RESULT -eq 0 ]; then
    echo ""
    echo "${GREEN}✅ Reproduction Passed (No vulnerabilities found). Did you already fix it?${NC}"
else
    echo ""
    echo "${RED}❌ Reproduction Successful! CI failure recreated locally.${NC}"
fi

exit $RESULT
