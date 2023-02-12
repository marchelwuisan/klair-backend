package configs

import (
	"context"
	"fmt"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func SetupFirebase() *auth.Client {
	serviceAccount, err := filepath.Abs("./klair-5dfd2-firebase-adminsdk-xwyv9-272493882e.json")
	if err != nil {
		panic("Unable to load serviceAccountKeys.json file")
	}
	// serviceAccount := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	opt := option.WithCredentialsFile(serviceAccount)
	//Firebase admin SDK initialization
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("Firebase load error")
	}
	//Firebase Auth
	auth, err := app.Auth(context.Background())
	if err != nil {
		fmt.Println(err)
		panic("Firebase load error")
	}
	return auth
}

// Client instance
var FB *auth.Client = SetupFirebase()
