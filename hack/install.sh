#!/bin/sh
set -e

OS_ARCH=$(uname -m | sed 's/x86_64/amd64/g')
OS_NAME=$(uname -s | tr '[:upper:]' '[:lower:]')
PLACEMENT=${XDG_CONFIG_HOME:-$HOME/.config}/kustomize/plugin/inloco.com.br/v1/sops
PLUGIN=$PLACEMENT/SOPS
RELEASE_URL=https://github.com/inloco/sops-kustomize-generator-plugin/releases/download/v0.0.0

mkdir -p $PLACEMENT
wget -O $PLUGIN ${RELEASE_URL}/plugin-${OS_NAME}-${OS_ARCH}
chmod +x $PLUGIN
