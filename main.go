package main

import (
	"io/ioutil"
	"log"
	"os"

	"go.mozilla.org/sops/aes"
	sopsYAML "go.mozilla.org/sops/stores/yaml"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func main() {
	filePath := os.Args[1]

	encryptedData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panic(err)
	}

	decryptedData, err := decrypt(encryptedData)
	if err != nil {
		log.Panic(err)
	}

	secret := coreV1.Secret{}
	if err := yaml.Unmarshal(decryptedData, &secret); err != nil {
		log.Panic(err)
	}
	secret.TypeMeta = metaV1.TypeMeta{
		APIVersion: "v1",
		Kind:       "Secret",
	}

	output, err := yaml.Marshal(secret)
	if err != nil {
		log.Panic(err)
	}

	if _, err := os.Stdout.Write(output); err != nil {
		log.Panic(err)
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
