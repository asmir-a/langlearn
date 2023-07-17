package db

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Conn *pgxpool.Pool

func getSecretFromAwsSecretsManagerUsingSdk() { //not needed, gettings from the environment variables is the same; and it wont work now; changes to the taskRun roles are required
	secretName := "username_value"
	region := "ap-northeast-2"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Println("there was an error loading the config: ", err)
		return
	}

	svc := secretsmanager.NewFromConfig(config)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Println("there was an error getting the secret: ", err)
		return
	}

	var secretString string = *result.SecretString
	log.Println("the secret is: ", secretString)
}

func init() {
	initDbConn()
}

func initDbConn() {
	godotenv.Load() //need to refactor so that this is only used when we are in dev mode
	DB_STRING := os.Getenv("DB_STRING")
	log.Println("new revision")
	log.Println("the db string is: ", DB_STRING)

	var err error
	Conn, err = pgxpool.New(context.Background(), DB_STRING) //why context background

	if err != nil {
		log.Fatal("could not connect to the database:\n", err)
	}
}
