package main

import (
	"log"
	"os"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/aes"
	sopsYAML "github.com/getsops/sops/v3/stores/yaml"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
	"strconv"
)

func main() {
	filePath := os.Args[1]

	encryptedData, err := os.ReadFile(filePath)
	if err != nil {
		log.Panic(filePath, ": ", err)
	}

	noDecryptEnv := os.Getenv("SOPS_NO_DECRYPT")

	noDecrypt := false
	if noDecryptEnv != "" {
		v, err := strconv.ParseBool(noDecryptEnv)
		if err != nil {
			log.Panicf("invalid SOPS_NO_DECRYPT value %q: %v", noDecryptEnv, err)
		}
		noDecrypt = v
	}

	if !noDecrypt {
		runDecrypt(filePath, encryptedData)
	} else {
		runNoDecrypt(filePath, encryptedData)
	}
}

func runDecrypt(filePath string, encryptedData []byte) {
	decryptedData, err := decrypt(encryptedData)
	if err != nil {
		log.Panic(filePath, ": ", err)
	}

	secret, err := makeSecret(decryptedData)
	if err != nil {
		log.Panic(filePath, ": ", err)
	}

	if _, err := os.Stdout.Write(secret); err != nil {
		log.Panic(filePath, ": ", err)
	}
}

func runNoDecrypt(filePath string, encryptedData []byte) {
	secret, err := makeEncryptedSecret(encryptedData)
	if err != nil {
		log.Panic(filePath, ": ", err)
	}

	if _, err := os.Stdout.Write(secret); err != nil {
		log.Panic(filePath, ": ", err)
	}
}

func decrypt(data []byte) ([]byte, error) {
	// Initialize a Sops JSON store
	store := &sopsYAML.Store{}

	// Load SOPS file and access the data key
	tree, err := store.LoadEncryptedFile(data)
	if err != nil {
		return nil, err
	}
	key, err := tree.Metadata.GetDataKey()
	if err != nil {
		return nil, err
	}

	// Decrypt the tree
	if _, err := tree.Decrypt(key, aes.NewCipher()); err != nil {
		return nil, err
	}

	return store.EmitPlainFile(tree.Branches)
}

func makeSecret(data []byte) ([]byte, error) {
	secret := coreV1.Secret{}
	if err := yaml.Unmarshal(data, &secret); err != nil {
		return nil, err
	}
	secret.TypeMeta = metaV1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Secret",
	}

	return yaml.Marshal(secret)
}

func makeEncryptedSecret(data []byte) ([]byte, error) {
	store := &sopsYAML.Store{}
	tree, err := store.LoadEncryptedFile(data)
	if err != nil {
		return nil, err
	}

	secret := coreV1.Secret{}
	secret.TypeMeta = metaV1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Secret",
	}

	secret.Data = getDataNoDecrypt(tree.Branches)
	secret.StringData = getStringDataNoDecrypt(tree.Branches)

	metadata := getMetadataNoDecrypt(tree.Branches)
	secret.ObjectMeta.Name = metadata["name"]
	secret.ObjectMeta.Namespace = metadata["namespace"]

	return yaml.Marshal(secret)
}

func getDataNoDecrypt(branches sops.TreeBranches) map[string][]byte {
	for _, item := range branches[0] {
		if item.Key == "data" {
			var result = make(map[string][]byte)
			dataFields := item.Value.(sops.TreeBranch)
			for _, df := range dataFields {
				result[df.Key.(string)] = []byte("no-decrypt")
			}
			return result
		}
	}
	return nil
}

func getStringDataNoDecrypt(branches sops.TreeBranches) map[string]string {
	for _, item := range branches[0] {
		if item.Key == "stringData" {
			var result = make(map[string]string)
			stringData := item.Value.(sops.TreeBranch)
			for _, df := range stringData {
				result[df.Key.(string)] = "no-decrypt"
			}
			return result
		}
	}
	return nil
}

func getMetadataNoDecrypt(branches sops.TreeBranches) map[string]string {
	for _, item := range branches[0] {
		if item.Key == "metadata" {
			var result = make(map[string]string)

			metadata := item.Value.(sops.TreeBranch)
			for _, mdta := range metadata {
				if mdta.Key.(string) == "name" {
					result["name"] = mdta.Value.(string)
				} else if mdta.Key.(string) == "namespace" {
					result["namespace"] = mdta.Value.(string)
				}
			}
			return result
		}
	}
	return nil
}
