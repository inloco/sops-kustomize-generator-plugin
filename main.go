package main

import (
	"log"
	"os"

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
