package service

import (
	"LoginAPI/models"
	"LoginAPI/repository"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

const (
	region         = "ap-south-1" // Update with your desired region
	userPoolID     = "ap-south-1_PsMRSTJ4p"
	clientID       = "1bj3i8ln4lnlei1ldg0dcg0efi"
	clientSecret   = "dhncbj3hti4j6gd4a9iil239lb2mj2495p9im6r7f3ekovifj52"
	identityPoolID = "ap-south-1:b92ae7d3-c02a-4988-b56d-804c3c5df777"
)

type IAuthService interface {
	CreateCognitoUser(user *models.LoginDetails) (string, error)
	LoginUser(user *models.LoginDetails) (string, error)
}

type authService struct {
	repo repository.IAuthRepository
}

var authServiceObj IAuthService

func InitAuthService(repoObj repository.IAuthRepository) IAuthService {
	if authServiceObj == nil {
		authServiceObj = &authService{
			repo: repoObj,
		}
	}
	return authServiceObj
}

func (s *authService) CreateCognitoUser(user *models.LoginDetails) (string, error) {
	// Create a session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	// Create a Cognito Identity Provider client
	cognitoClient := cognitoidentityprovider.New(sess)

	// Create a new user
	createUserInput := &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId:        aws.String(userPoolID),
		Username:          aws.String(user.UserName),
		TemporaryPassword: aws.String(user.Password), // Provide a temporary password
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(user.Email),
			},
		},
	}

	_, err = cognitoClient.AdminCreateUser(createUserInput)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Println("User created successfully")

	secretHash := calculateSecretHash(user.UserName, clientID, clientSecret)

	// Authenticate the user
	authInput := &cognitoidentityprovider.AdminInitiateAuthInput{
		UserPoolId: aws.String(userPoolID),
		ClientId:   aws.String(clientID),
		AuthFlow:   aws.String("ADMIN_USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME":    aws.String(user.UserName),
			"PASSWORD":    aws.String(user.Password), // Provide the temporary password here
			"SECRET_HASH": aws.String(secretHash),
		},
	}

	authOutput, err := cognitoClient.AdminInitiateAuth(authInput)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if authOutput.ChallengeName != nil {
		fmt.Println("Challenge received:", *authOutput.ChallengeName)

		if *authOutput.ChallengeName == "NEW_PASSWORD_REQUIRED" {
			// Respond to the NEW_PASSWORD_REQUIRED challenge
			err := respondToNewPasswordChallenge(cognitoClient, *authOutput.Session, user.UserName, user.Password, secretHash)
			if err != nil {
				fmt.Println(err)
				return "", err
			}
		} else {
			log.Fatalf("Unhandled challenge: %s", *authOutput.ChallengeName)
		}
	} else {
		fmt.Println("User authenticated successfully")
		fmt.Println("Access Token:", *authOutput.AuthenticationResult.AccessToken)
		fmt.Println("Refresh Token:", *authOutput.AuthenticationResult.RefreshToken)
	}

	fmt.Println("User authenticated successfully")
	fmt.Println("Access Token:", *authOutput.AuthenticationResult.AccessToken)
	token := *authOutput.AuthenticationResult.AccessToken
	return token, nil
}

func (s *authService) LoginUser(user *models.LoginDetails) (string, error) {
	// Create a session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	// Create a Cognito Identity Provider client
	cognitoClient := cognitoidentityprovider.New(sess)

	secretHash := calculateSecretHash(user.UserName, clientID, clientSecret)

	// Authenticate the user
	authInput := &cognitoidentityprovider.AdminInitiateAuthInput{
		UserPoolId: aws.String(userPoolID),
		ClientId:   aws.String(clientID),
		AuthFlow:   aws.String("ADMIN_USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME":    aws.String(user.UserName),
			"PASSWORD":    aws.String(user.Password), // Provide the temporary password here
			"SECRET_HASH": aws.String(secretHash),
		},
	}

	authOutput, err := cognitoClient.AdminInitiateAuth(authInput)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if authOutput.ChallengeName != nil {
		fmt.Println("Challenge received:", *authOutput.ChallengeName)

		if *authOutput.ChallengeName == "NEW_PASSWORD_REQUIRED" {
			// Respond to the NEW_PASSWORD_REQUIRED challenge
			err := respondToNewPasswordChallenge(cognitoClient, *authOutput.Session, user.UserName, user.Password, secretHash)
			if err != nil {
				fmt.Println(err)
				return "", err
			}
		} else {
			log.Fatalf("Unhandled challenge: %s", *authOutput.ChallengeName)
		}
	} else {
		fmt.Println("User authenticated successfully")
		fmt.Println("Access Token:", *authOutput.AuthenticationResult.AccessToken)
		fmt.Println("Refresh Token:", *authOutput.AuthenticationResult.RefreshToken)
	}

	fmt.Println("User authenticated successfully")
	fmt.Println("Access Token:", *authOutput.AuthenticationResult.AccessToken)
	token := *authOutput.AuthenticationResult.AccessToken
	return token, nil
}

func calculateSecretHash(username, clientID, clientSecret string) string {
	msg := username + clientID
	hmac := hmac.New(sha256.New, []byte(clientSecret))
	hmac.Write([]byte(msg))
	secretHash := base64.StdEncoding.EncodeToString(hmac.Sum(nil))
	return secretHash
}

func respondToNewPasswordChallenge(client *cognitoidentityprovider.CognitoIdentityProvider, session, username, newPassword, secretHash string) error {
	challengeResponse := &cognitoidentityprovider.AdminRespondToAuthChallengeInput{
		UserPoolId:    aws.String(userPoolID),
		ClientId:      aws.String(clientID),
		ChallengeName: aws.String("NEW_PASSWORD_REQUIRED"),
		Session:       aws.String(session),
		ChallengeResponses: map[string]*string{
			"USERNAME":     aws.String(username),
			"NEW_PASSWORD": aws.String(newPassword),
			"SECRET_HASH":  aws.String(secretHash),
		},
	}

	_, err := client.AdminRespondToAuthChallenge(challengeResponse)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("User authentication completed successfully")
	fmt.Println("New password set successfully")
	return nil
}
