#!/usr/bin/env bash
set -euo pipefail

LOG_FILE="install-tools.log"
SUMMARY_FILE="${SUMMARY_FILE:-install-summary.json}"
VERSIONS_FILE="${VERSIONS_FILE:-.tool-versions.json}"
DRY_RUN=false
INSTALL_TOOLS=(all)

for arg in "$@"; do
  case $arg in
    --dry-run)
      DRY_RUN=true
      echo "[Dry Run] No changes will be made. Commands will be printed only."
      ;;
    --tools=*)
      IFS=',' read -ra INSTALL_TOOLS <<< "${arg#*=}"
      ;;
    --summary-path=*)
      SUMMARY_FILE="${arg#*=}"
      ;;
  esac
done

exec > >(tee -a "$LOG_FILE") 2>&1

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

SUMMARY_JSON="{}"
EXPECTED_JSON="{}"

if [[ -f "$VERSIONS_FILE" ]]; then
  EXPECTED_JSON=$(<"$VERSIONS_FILE")
fi

log_step() {
    echo -e "\n${YELLOW}🔧 $(date '+%Y-%m-%d %H:%M:%S') - $1${NC}"
}

run_cmd() {
    log_step "$1"
    shift
    if $DRY_RUN; then
        echo "[Dry Run] $*"
    else
        if "$@"; then
            echo -e "${GREEN}✅ Success: $1${NC}"
        else
            echo -e "${RED}❌ Failed: $1${NC}"
            exit 1
        fi
    fi
}

add_summary() {
  local name=$1
  local version=$2
  SUMMARY_JSON=$(echo "$SUMMARY_JSON" | jq --arg name "$name" --arg ver "$version" '. + {($name): $ver}')

  local expected_version
  expected_version=$(echo "$EXPECTED_JSON" | jq -r --arg name "$name" '.[$name] // empty')

  if [[ -n "$expected_version" && "$version" != "$expected_version" ]]; then
    echo -e "${RED}⚠️ Version mismatch for $name: expected $expected_version, got $version${NC}"
  fi
}

should_run() {
    [[ " ${INSTALL_TOOLS[*]} " =~ " all " || " ${INSTALL_TOOLS[*]} " =~ " $1 " ]]
}

log_step "Installing OS dependencies"
run_cmd "Install OS dependencies" sudo apt-get update -y && sudo apt-get install -y \
  curl unzip git jq gnupg software-properties-common ca-certificates lsb-release tar

# Terraform
if should_run terraform; then
  log_step "Installing Terraform"
  if ! $DRY_RUN; then
    curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
    echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" \
      | sudo tee /etc/apt/sources.list.d/hashicorp.list
  fi
  run_cmd "Update apt & install terraform" sudo apt-get update && sudo apt-get install -y terraform
  if ! $DRY_RUN; then
    TERRAFORM_VERSION=$(terraform version -json | jq -r .terraform_version)
    add_summary terraform "$TERRAFORM_VERSION"
  fi
fi

# AWS CLI
if should_run awscli; then
  log_step "Installing AWS CLI"
  run_cmd "Download AWS CLI" curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
  run_cmd "Unzip AWS CLI" unzip awscliv2.zip
  run_cmd "Install AWS CLI" sudo ./aws/install
  rm -rf awscliv2.zip aws
  AWS_VERSION=$(aws --version 2>&1 | awk '{print $1}' | cut -d/ -f2)
  add_summary awscli "$AWS_VERSION"
fi

# Terraform Docs
if should_run terraform-docs; then
  log_step "Installing terraform-docs"
  run_cmd "Download terraform-docs" curl -sLo terraform-docs.tar.gz https://github.com/terraform-docs/terraform-docs/releases/download/v0.12.0/terraform-docs-v0.12.0-$(uname)-amd64.tar.gz
  run_cmd "Extract terraform-docs" tar -xzf terraform-docs.tar.gz
  run_cmd "Move terraform-docs" sudo mv terraform-docs /usr/local/bin/
  rm terraform-docs.tar.gz
  TERRADOCS_VERSION=$(terraform-docs --version | awk '{print $2}')
  add_summary terraform-docs "$TERRADOCS_VERSION"
fi

# Terragrunt
if should_run terragrunt; then
  log_step "Installing Terragrunt"
  TG_VERSION=$(curl -s https://api.github.com/repos/gruntwork-io/terragrunt/releases/latest | jq -r .tag_name)
  run_cmd "Download Terragrunt" curl -Lo terragrunt https://github.com/gruntwork-io/terragrunt/releases/download/${TG_VERSION}/terragrunt_$(uname -s)_amd64
  run_cmd "Move Terragrunt" chmod +x terragrunt && sudo mv terragrunt /usr/local/bin/
  add_summary terragrunt "$TG_VERSION"
fi

# Terrascan
if should_run terrascan; then
  log_step "Installing Terrascan"
  run_cmd "Install Terrascan" curl -s https://runterrascan.io/install.sh | bash
  run_cmd "Move Terrascan" sudo mv terrascan /usr/local/bin/
  TERRASCAN_VERSION=$(terrascan version | head -n1 | awk '{print $3}')
  add_summary terrascan "$TERRASCAN_VERSION"
fi

# TFLint
if should_run tflint; then
  log_step "Installing TFLint"
  run_cmd "Install TFLint" curl -s https://raw.githubusercontent.com/terraform-linters/tflint/master/install_linux.sh | bash
  run_cmd "Move TFLint" sudo mv tflint /usr/local/bin/
  TFLINT_VERSION=$(tflint --version | head -n1 | awk '{print $2}')
  add_summary tflint "$TFLINT_VERSION"
fi

# TFSec
if should_run tfsec; then
  log_step "Installing TFSec"
  run_cmd "Download tfsec" curl -sLo tfsec https://github.com/aquasecurity/tfsec/releases/latest/download/tfsec-$(uname)-amd64
  run_cmd "Move tfsec" chmod +x tfsec && sudo mv tfsec /usr/local/bin/
  TFSEC_VERSION=$(tfsec --version | awk '{print $3}')
  add_summary tfsec "$TFSEC_VERSION"
fi

# Trivy
if should_run trivy; then
  log_step "Installing Trivy"
  TRIVY_VERSION=$(curl -s https://api.github.com/repos/aquasecurity/trivy/releases/latest | jq -r .tag_name)
  run_cmd "Install Trivy" curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin ${TRIVY_VERSION}
  TRIVY_VER=$(trivy --version | head -n1 | awk '{print $2}')
  add_summary trivy "$TRIVY_VER"
fi

# Infracost
if should_run infracost; then
  log_step "Installing Infracost"
  run_cmd "Install Infracost" curl -s https://raw.githubusercontent.com/infracost/infracost/master/scripts/install.sh | sh
  run_cmd "Move Infracost" sudo mv infracost /usr/local/bin/
  INFRACOST_VERSION=$(infracost --version | awk '{print $3}')
  add_summary infracost "$INFRACOST_VERSION"
fi

# tfupdate
if should_run tfupdate; then
  log_step "Installing tfupdate"
  run_cmd "Install tfupdate" go install github.com/minamijoyo/tfupdate@latest
  run_cmd "Move tfupdate" sudo mv ~/go/bin/tfupdate /usr/local/bin/
  TFUPDATE_VERSION=$(tfupdate --version | awk '{print $3}')
  add_summary tfupdate "$TFUPDATE_VERSION"
fi

# hcledit
if should_run hcledit; then
  log_step "Installing hcledit"
  run_cmd "Install hcledit" go install github.com/minamijoyo/hcledit@latest
  run_cmd "Move hcledit" sudo mv ~/go/bin/hcledit /usr/local/bin/
  HCLEDIT_VERSION=$(hcledit --version | awk '{print $3}')
  add_summary hcledit "$HCLEDIT_VERSION"
fi

if ! $DRY_RUN; then
  echo "$SUMMARY_JSON" | jq . > "$SUMMARY_FILE"
  echo -e "\n${GREEN}📦 Tool summary written to $SUMMARY_FILE${NC}"
fi

echo -e "\n${GREEN}✅ All tools installed successfully at $(date '+%Y-%m-%d %H:%M:%S')${NC}"
