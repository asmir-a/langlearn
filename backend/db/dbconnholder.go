package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Conn *pgxpool.Pool

func getSecretFromAwsSecretsManager() {
	//secretName := "DB_STRING"
	//region := "ap-northeast-2"

	//config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	//if err != nil {
	//	httperrors.Fatal(err)
	//}

	//svc := secretsmanager.NewFromConfig(config)
	//input := &secretsmanager.GetSecretValueInput{
	//	SecretId: aws.String(secretName),
	//}

	//result, err := svc.GetSecretValue(context.TODO(), input)
	//if err != nil {
	//	httperrors.Fatal(err)
	//}

	//var secretString string = *result.SecretString
	//log.Println("the secret is: ", secretString)
	secretValue := os.Getenv("username_value")
	log.Println("the value of the username_value is: ", secretValue)
}

func init() {
	getSecretFromAwsSecretsManager()
	initDbConn()
}

func initDbConn() {
	godotenv.Load()
	_ = os.Getenv("DB_STRING")
	DB_STRING := "postgresql://postgres:qwertyuiop@langlearndb.cmhmoaojrw66.ap-northeast-2.rds.amazonaws.com:5432/langlearn"
	log.Println("DB STRING: ", DB_STRING)

	var err error
	Conn, err = pgxpool.New(context.Background(), DB_STRING) //why context background

	if err != nil {
		log.Fatal("could not connect to the database:\n", err)
	}
}
