package main

import (
	"flag"

	"github.com/thospol/minio/config"
	"github.com/thospol/minio/minio"
)

func main() {
	environment := flag.String("environment", "local", "set working environment")
	configs := flag.String("config", "config", "set configs path, default as: 'configs'")

	flag.Parse()

	// Init configuration
	if err := config.InitConfig(*configs, *environment); err != nil {
		panic(err)
	}
	//=======================================================

	// Init client minio connection
	conf := minio.Configuration{
		Host:            config.CF.Minio.Host,
		AccessKeyID:     config.CF.Minio.AccessKey,
		SecretAccessKey: config.CF.Minio.SecretKey,
	}
	if err := minio.NewConnection(conf); err != nil {
		panic(err)
	}
	//========================================================

	client := minio.GetClient()
	bucketName := config.CF.Minio.BucketName
	objectName := "profile.jpeg"
	filepath := "./assets/images/profile.jpeg"
	err := client.UploadImage(bucketName, objectName, filepath)
	if err != nil {
		panic(err)
	}
}
