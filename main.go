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
)

func main() {
	filePath := os.Args[1]

	encryptedData, err := os.ReadFile(filePath)
	if err != nil {
		log.Panic(filePath, ": ", err)
	}

	_, noDecryptEnv := os.LookupEnv("SOPS_NO_DECRYPT")
	shouldDecrypt := !noDecryptEnv

	var secret []byte
	if shouldDecrypt {
		log.Println("Running with decryption")
		secret, err = runDecrypt(encryptedData)
	} else {
		log.Println("[SOPS] Running in no-decrypt mode")
		secret, err = runNoDecrypt(encryptedData)
	}
	if err != nil {
		log.Panic(filePath, ": ", err)
	}
	if _, err := os.Stdout.Write(secret); err != nil {
		log.Panic(filePath, ": ", err)
	}
}

func runDecrypt(encryptedData []byte) ([]byte, error) {
	decryptedData, err := decrypt(encryptedData)
	if err != nil {
		return nil, err
	}

	return makeSecret(decryptedData)
}

func runNoDecrypt(encryptedData []byte) ([]byte, error) {
	return makeSecretNoDecrypt(encryptedData)
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

func makeSecretNoDecrypt(data []byte) ([]byte, error) {
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
