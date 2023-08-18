package main

import (
	"go.mozilla.org/sops/v3"
	"io/ioutil"
	"log"
	"os"

	sopsYAML "go.mozilla.org/sops/v3/stores/yaml"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func main() {
	filePath := os.Args[1]

	encryptedData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panic(filePath, ": ", err)
	}

	secret, err := toSecret(encryptedData)
	if err != nil {
		log.Panic(filePath, ": ", err)
	}

	if _, err := os.Stdout.Write(secret); err != nil {
		log.Panic(filePath, ": ", err)
	}
}

func toSecret(data []byte) ([]byte, error) {
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
	secret.Data = getData(tree.Branches)
	secret.StringData = getStringData(tree.Branches)

	metadata := getMetadata(tree.Branches)
	secret.ObjectMeta.Name = metadata["name"]
	secret.ObjectMeta.Namespace = metadata["namespace"]

	return yaml.Marshal(secret)
}

func getData(branches sops.TreeBranches) map[string][]byte {
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

func getStringData(branches sops.TreeBranches) map[string]string {
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

func getMetadata(branches sops.TreeBranches) map[string]string {
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
