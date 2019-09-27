package main

import (
	"io/ioutil"
	"log"
	"os"

	"go.mozilla.org/sops/decrypt"
	"go.mozilla.org/sops/stores"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type document struct {
	Data map[string]string `json:"data,omitempty"`
	SOPS *stores.Metadata  `json:"sops,omitempty"`
}

type payload struct {
	Type     coreV1.SecretType `json:"type,omitempty"`
	Metadata metaV1.ObjectMeta `json:"metadata,omitempty"`
	document `json:",inline"`
}

func main() {
	filePath := os.Args[1]

	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panic(err)
	}

	var encryptedPayload payload
	if err := yaml.Unmarshal(input, &encryptedPayload); err != nil {
		log.Panic(err)
	}

	encryptedData, err := yaml.Marshal(encryptedPayload.document)
	if err != nil {
		log.Panic(err)
	}

	decryptedData, err := decrypt.Data(encryptedData, "yaml")
	if err != nil {
		log.Panic(err)
	}

	decryptedPayload := coreV1.Secret{
		TypeMeta: metaV1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		Type:       encryptedPayload.Type,
		ObjectMeta: encryptedPayload.Metadata,
	}
	if err := yaml.Unmarshal(decryptedData, &decryptedPayload); err != nil {
		log.Panic(err)
	}

	output, err := yaml.Marshal(decryptedPayload)
	if err != nil {
		log.Panic(err)
	}

	os.Stdout.Write(output)
}
