#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Ensure we're in the root directory of the project
cd "$(dirname "$0")/.."

# Configuration
VERSION="v0.2.1"
SHA256="2a61c7248203d38087688aae146483054fef8f33620c2131d9f0ee220686963c"
HOMEBREW_TAP_DIR="../homebrew-tap"  # Adjust to your actual location

# Check if homebrew-tap directory exists
if [ ! -d "$HOMEBREW_TAP_DIR" ]; then
  echo -e "${RED}Error: Homebrew tap directory not found at $HOMEBREW_TAP_DIR${NC}"
  echo "Please clone your homebrew-tap repository to this location or adjust the path in this script."
  exit 1
fi

# Navigate to the homebrew-tap directory
cd "$HOMEBREW_TAP_DIR"

# Update the formula
echo -e "${YELLOW}Updating formula with version ${VERSION} and SHA256 ${SHA256}...${NC}"
if [ -f "Formula/scribe.rb" ]; then
  sed -i.bak "s|url \".*\"|url \"https://github.com/dikaio/scribe/archive/refs/tags/${VERSION}.tar.gz\"|" Formula/scribe.rb
  sed -i.bak "s|sha256 \".*\"|sha256 \"${SHA256}\"|" Formula/scribe.rb
  rm -f Formula/scribe.rb.bak
else
  echo -e "${RED}Error: Formula/scribe.rb not found${NC}"
  exit 1
fi

# Verify the changes
echo -e "${YELLOW}Verifying changes...${NC}"
grep -E "(url|sha256)" Formula/scribe.rb

# Ask for confirmation
echo
echo -e "${YELLOW}Changes have been made to Formula/scribe.rb${NC}"
echo -e "Proceed with committing and pushing these changes? [y/N] "
read -r CONFIRM

if [[ "$CONFIRM" =~ ^[Yy]$ ]]; then
  # Commit and push
  git add Formula/scribe.rb
  git commit -m "Update scribe to ${VERSION}"
  git push
  echo -e "${GREEN}Changes committed and pushed successfully!${NC}"
  echo "Users can now install with: brew install dikaio/tap/scribe"
else
  echo -e "${YELLOW}No changes committed. Please review and commit manually.${NC}"
fi