#!/bin/sh

OS_NAME_LOWERCASE=$(uname -s | tr '[:upper:]' '[:lower:]')
PLACEMENT=${XDG_CONFIG_HOME:-$HOME/.config}/kustomize/plugin/inloco.com.br/v1/sops
PLUGIN=$PLACEMENT/SOPS
RELEASE_URL=https://github.com/inloco/sops-kustomize-generator-plugin/releases/download/v1.1.2

mkdir -p $PLACEMENT
wget -O $PLUGIN ${RELEASE_URL}/plugin-${OS_NAME_LOWERCASE}-amd64
chmod +x $PLUGIN
