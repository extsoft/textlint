#!/usr/bin/env bash
#MISE description="build release artifacts"
set -o pipefail -o errexit -o nounset

# Define platforms: format is "GOOS-GOARCH" (variants set internally for compatibility)
PLATFORMS=(
  "linux-amd64"
  "linux-arm64"
  "linux-386"
  "windows-amd64"
  "windows-arm64"
  "windows-386"
  "darwin-amd64"
  "darwin-arm64"
)

# Get build variant environment variables for a platform
get_build_variants() {
  local GOOS=$1
  local GOARCH=$2

  local GOAMD64=""
  local GO386=""
  local GOARM64=""

  case "$GOARCH" in
    amd64)
      GOAMD64="v1"
      ;;
    386)
      GO386="sse2"
      ;;
    arm64)
      GOARM64="v8.0"
      ;;
  esac

  echo "$GOAMD64|$GO386|$GOARM64"
}

# Build and archive a platform
build_and_archive_platform() {
  local dist_dir=$1
  local version=$2
  local platform=$3
  IFS='-' read -r GOOS GOARCH <<< "$platform"
  local archive_name="prosefmt-${version}-${platform}"
  local staging_dir="${dist_dir}/${archive_name}"
  mkdir -p "$staging_dir"

  IFS='|' read -r GOAMD64 GO386 GOARM64 <<< "$(get_build_variants "$GOOS" "$GOARCH")"

  local binary_name="prosefmt"
  if [ "$GOOS" = "windows" ]; then
    binary_name="prosefmt.exe"
  fi

  echo "Building $platform..."

  # Set build environment
  export CGO_ENABLED=0
  export GOOS="$GOOS"
  export GOARCH="$GOARCH"
  [ -n "$GOAMD64" ] && export GOAMD64="$GOAMD64"
  [ -n "$GO386" ] && export GO386="$GO386"
  [ -n "$GOARM64" ] && export GOARM64="$GOARM64"

  # Build binary to staging directory
  go build -ldflags="-s -w" -o "$staging_dir/$binary_name" .

  # Copy LICENSE and README.md to staging directory
  cp LICENSE "$staging_dir/"
  cp README.md "$staging_dir/"

  # Use absolute path to ensure it works after pushd
  local archive_path
  archive_path=$(cd "$dist_dir" && pwd)/"${archive_name}"
  pushd "$staging_dir" > /dev/null
  if [ "$GOOS" = "windows" ]; then
    # -r zips the current directory contents into the destination
    zip -r "${archive_path}.zip" . > /dev/null
    echo "Created ${archive_path}.zip"
    generate_checksum "${archive_path}.zip"
  else
    # -czf creates the tarball from the current directory contents
    tar -czf "${archive_path}.tar.gz" .
    echo "Created ${archive_path}.tar.gz"
    generate_checksum "${archive_path}.tar.gz"
  fi
  popd > /dev/null

  # Clean up the staging folder
  rm -rf "$staging_dir"
}

# Generate checksum for an archive
generate_checksum() {
  local archive_path=$1
  local checksum_file="${archive_path}.sha256"

  if command -v sha256sum > /dev/null; then
    sha256sum "$archive_path" | cut -d' ' -f1 > "$checksum_file"
  elif command -v shasum > /dev/null; then
    shasum -a 256 "$archive_path" | cut -d' ' -f1 > "$checksum_file"
  else
    echo "Error: Neither sha256sum nor shasum found" >&2
    exit 1
  fi

  echo "Generated checksum: $checksum_file"
}

# Main execution
## Get abbreviated version (tag or commit)
VERSION=$(git describe --tags --always --abbrev=7 2>/dev/null || git rev-parse --short=7 HEAD)
echo "Building release artifacts for version: $VERSION"
DIST_DIR="${DIST_DIR:-dist}"
DIST_VERSION_DIR="$DIST_DIR/$VERSION"
echo "Output directory: ${DIST_VERSION_DIR}"
## Remove version directory if present
if [ -d "${DIST_VERSION_DIR}" ]; then
  echo "Removing directory: ${DIST_VERSION_DIR}"
  rm -rf "${DIST_VERSION_DIR}"
fi
## Create version directory
mkdir -p "${DIST_VERSION_DIR}"
# Build and archive each platform
for platform in "${PLATFORMS[@]}"; do
  build_and_archive_platform "${DIST_VERSION_DIR}" "${VERSION}" "$platform"
done

echo "Release artifacts built successfully in $DIST_DIR/$VERSION"
