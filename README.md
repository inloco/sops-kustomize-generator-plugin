# SOPS Kustomize Generator Plugin

It is a plugin for [Kustomize](https://github.com/kubernetes-sigs/kustomize) that allows you to use Kubernetes Secrets encrypted with [SOPS](https://github.com/mozilla/sops) as a generator.

## Getting Started

### Install

To install this plugin on Kustomize, download the binary to Kustomize Plugin folder with `apiVersion: inloco.com.br/v1` and `kind: SOPS`. Then make it executable.

#### Linux 64-bits and/or macOS 64-bits

```bash
VERSION=$(wget -qO- https://api.github.com/repos/inloco/sops-kustomize-generator-plugin/releases/latest | jq -r '.tag_name')
wget -qO- https://github.com/inloco/sops-kustomize-generator-plugin/releases/download/${VERSION}/install.sh | sh
```

#### Manual Build and Install for Other Systems and/or Architectures

```bash
git clone https://github.com/inloco/sops-kustomize-generator-plugin

cd sops-kustomize-generator-plugin

go get -d -v ./...

go build -a -installsuffix cgo -ldflags '-extldflags "-static" -s -w' -tags netgo -v ./...

PLACEMENT=${XDG_CONFIG_HOME:-$HOME/.config}/kustomize/plugin/inloco.com.br/v1/sops

mkdir -p $PLACEMENT

mv ./sops-kustomize-generator-plugin $PLACEMENT/SOPS

cd ..

rm -fR sops-kustomize-generator-plugin
```

### Using

We can start with a regular Kubernetes Secret in its YAML format.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
data:
  username: YWRtaW4=
  password: MWYyZDFlMmU2N2Rm
```

To convert it to a file that will be processed by the plugin, we replace `apiVersion: v1` with `apiVersion: inloco.com.br/v1` and `kind: Secret` with `kind: SOPS`.

```yaml
apiVersion: inloco.com.br/v1
kind: SOPS
metadata:
  name: mysecret
type: Opaque
data:
  username: YWRtaW4=
  password: MWYyZDFlMmU2N2Rm
```

Finally we encrypt it using SOPS with the following command:

```bash
sops --encrypt --encrypted-regex '^(data|stringData)$' --in-place ./secret.yaml
```

Now we can specify `./secret.yaml` as a generator on `kustomization.yaml`:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
generators:
  - ./secret.yaml
```

## Notes

- Remember to use `--enable-alpha-plugins` flag when running `kustomize build`.
- You may need to use environment variables, such as `AWS_PROFILE`, to configure SOPS decryption when running Kustomize.
- Integrity checks are disabled on SOPS decryption, this is done to prevent integrity failures due to Kustomize sorting the keys of original YAML file.
- This documentation assumes that you are familiar with [Kustomize](https://github.com/kubernetes-sigs/kustomize) and [SOPS](https://github.com/mozilla/sops), read their documentation if necessary.
- To make the generator behave like a patch, you might want to set `kustomize.config.k8s.io/behavior` annotation to `"merge"`. The other internal annotations described on [Kustomize Plugins Guide](https://kubernetes-sigs.github.io/kustomize/guides/plugins/#generator-options) are also supported.
