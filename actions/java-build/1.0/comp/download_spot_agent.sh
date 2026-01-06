#!/usr/bin/env bash

set -exo pipefail

# Define version array in chronological order, including latest
versions=(
    3.16
    3.17
    3.18
    3.19
    3.20
    3.21
    4.0
    1.1
    1.2
    1.3
    1.4
    1.5
    1.6
    2.0
    2.4
    latest
)

# Configuration
compDir="/opt/action/comp"
baseUrl="https://terminus-dice.oss-cn-hangzhou.aliyuncs.com/spot/java-agent/action/release"
maxRetries=3
retryDelay=5

# Logging function
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

# Download function with retry mechanism
download_with_retry() {
    local url="$1"
    local output_file="$2"
    local description="$3"
    local accept_404="${4:-false}"  # New parameter to accept 404 errors

    for attempt in $(seq 1 $maxRetries); do
        log "Downloading $description (attempt $attempt/$maxRetries)..."

        if curl -f -L -o "$output_file" "$url" --progress-bar; then
            log "✓ $description downloaded successfully"
            echo "0"  # Success
            return 0
        else
            local http_code=$(curl -s -o /dev/null -w "%{http_code}" "$url")
            if [ "$accept_404" = "true" ] && [ "$http_code" = "404" ]; then
                log "⚠ $description not available (HTTP 404) - skipping"
                echo "2"  # 404 - not available
                return 0  # Don't exit on 404
            else
                log "✗ $description download failed (attempt $attempt/$maxRetries)"
                if [ $attempt -lt $maxRetries ]; then
                    log "Waiting ${retryDelay} seconds before retry..."
                    sleep $retryDelay
                fi
            fi
        fi
    done

    log "✗ $description download failed, reached maximum retry attempts"
    echo "1"  # Failed
    return 0  # Don't exit on failure
}

# Verify tar.gz file integrity
verify_tar_file() {
    local file="$1"
    local description="$2"

    if tar -tzf "$file" >/dev/null 2>&1; then
        log "✓ $description verification successful"
        return 0
    else
        log "✗ $description verification failed, file may be corrupted"
        return 1
    fi
}

# Create directory structure
log "Creating directory structure..."
mkdir -p "$compDir"
spotAgentDir="$compDir/spot-agent"
mkdir -p "$spotAgentDir"

# Download spot-agent and spot-agent-jdk17 for all versions
log "Starting download of all versions of spot-agent..."
for version in "${versions[@]}"; do
    log "Processing version $version..."

    # Create version directory
    versionDir="$spotAgentDir/$version"
    mkdir -p "$versionDir"

    # Download spot-agent (required)
    spotAgentFile="$versionDir/spot-agent.tar.gz"
    spotAgentUrl="$baseUrl/$version/spot-agent.tar.gz"

    if download_with_retry "$spotAgentUrl" "$spotAgentFile" "spot-agent v$version"; then
        if ! verify_tar_file "$spotAgentFile" "spot-agent v$version"; then
            log "Warning: spot-agent v$version file verification failed, removing corrupted file"
            rm -f "$spotAgentFile"
        fi
    fi

    # Download spot-agent-jdk17 (optional - accept 404)
    spotAgentJdk17File="$versionDir/spot-agent-jdk17.tar.gz"
    spotAgentJdk17Url="$baseUrl/$version/spot-agent-jdk17.tar.gz"

    download_result=$(download_with_retry "$spotAgentJdk17Url" "$spotAgentJdk17File" "spot-agent-jdk17 v$version" "true")
    case $download_result in
        0)  # Success
            if ! verify_tar_file "$spotAgentJdk17File" "spot-agent-jdk17 v$version"; then
                log "Warning: spot-agent-jdk17 v$version file verification failed, removing corrupted file"
                rm -f "$spotAgentJdk17File"
            fi
            ;;
        2)  # 404 - not available
            log "Info: spot-agent-jdk17 v$version is not available for this version"
            ;;
        1)  # Other errors
            log "Warning: spot-agent-jdk17 v$version download failed with errors"
            ;;
    esac

    log "Version $version processing completed"
done

# Display download statistics
log "Download completed!"
log "Number of versions downloaded: ${#versions[@]}"
log "Each version contains: spot-agent.tar.gz (required) and spot-agent-jdk17.tar.gz (optional)"
log "Files saved in: $spotAgentDir"