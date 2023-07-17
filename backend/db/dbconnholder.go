package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/asmir-a/langlearn/backend/httperrors"
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

type dbInfo struct {
	Engine   string `json:"engine"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Dbname   string `json:"dbname"`
}

func turnDbInfoFromJsonToString(dbInfoJson []byte) string {
	var dbInfoStruct dbInfo
	json.Unmarshal(dbInfoJson, &dbInfoStruct)
	dbInfoString := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		dbInfoStruct.Engine,
		dbInfoStruct.Username,
		dbInfoStruct.Password,
		dbInfoStruct.Host,
		dbInfoStruct.Port,
		dbInfoStruct.Dbname,
	)
	return dbInfoString
}

func init() {
	initDbConn()
}

func initDbConn() {
	godotenv.Load() //need to refactor so that this is only used when we are in dev mode
	DB_STRING := os.Getenv("DB_STRING")
	if DB_STRING == "" {
		httperrors.Fatal(errors.New("DB_STRING is empty"))
	}
	DB_STRING = turnDbInfoFromJsonToString([]byte(DB_STRING))

	var err error
	Conn, err = pgxpool.New(context.Background(), DB_STRING) //why context background

	if err != nil {
		log.Fatal("could not connect to the database:\n", err)
	}
}
