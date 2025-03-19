package utils

import (
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
		privateKeyPath := os.Getenv("OCI_PRIVATE_KEY_PATH")
		tenancy := os.Getenv("OCI_TENANCY")
		region := os.Getenv("OCI_REGION")

		if user == "" || fingerprint == "" || privateKeyPath == "" || tenancy == "" || region == "" {
			log.Fatalln("não foi possível configurar o OCI")
		}

		privateKey, err := os.ReadFile(privateKeyPath)
		if err != nil {
			log.Fatalf("private key do OCI não foi encontrada: %v", err)
		}

		config := common.NewRawConfigurationProvider(tenancy, user, region, fingerprint, string(privateKey), nil)

		client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(config)
		if err != nil {
			log.Fatalf("não foi possível criar o cliente OCI: %v", err)
		}

		instance = &client
	})
	return instance
}
