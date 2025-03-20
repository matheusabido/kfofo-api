package utils

import (
	"encoding/base64"
	"log"
	"os"
	"sync"

	"github.com/oracle/oci-go-sdk/v49/common"
	"github.com/oracle/oci-go-sdk/v49/objectstorage"
)

var instance *objectstorage.ObjectStorageClient
var once sync.Once

func GetClient() *objectstorage.ObjectStorageClient {
	once.Do(func() {
		user := os.Getenv("OCI_USER")
		fingerprint := os.Getenv("OCI_FINGERPRINT")
		privateKey64 := os.Getenv("OCI_PRIVATE_KEY_BASE64")
		tenancy := os.Getenv("OCI_TENANCY")
		region := os.Getenv("OCI_REGION")

		if user == "" || fingerprint == "" || privateKey64 == "" || tenancy == "" || region == "" {
			log.Fatalln("não foi possível configurar o OCI")
		}

		var err error
		privateKeyDecoded, err := base64.StdEncoding.DecodeString(privateKey64)
		if err != nil {
			log.Fatalf("não foi possível decodificar a chave privada do OCI: %v", err)
		}

		config := common.NewRawConfigurationProvider(tenancy, user, region, fingerprint, string(privateKeyDecoded), nil)

		client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(config)
		if err != nil {
			log.Fatalf("não foi possível criar o cliente OCI: %v", err)
		}

		instance = &client
	})
	return instance
}
