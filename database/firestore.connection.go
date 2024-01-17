package database

import (
	"log"
	"context"

	"github.com/InnoFours/skin-savvy/config"

	"google.golang.org/api/option"
	"cloud.google.com/go/firestore"
)

func FirestoreConnection() (*firestore.Client, error) {
	opt := option.WithCredentialsFile("D:/Tech Projects/SkinSavvyApi/SkinSavvyAPI/service-account-key.json")
    client, err := firestore.NewClient(context.Background(), config.ConfigFirebaseProjectId(), opt)
    if err != nil {
        log.Println("error connecting to firestore", err.Error())
    }

	return client, nil
}	
