package service

import (
	localAwsPkg "LoginAPI/aws"
	"LoginAPI/config"
	customerrors "LoginAPI/custom_errors"
	"LoginAPI/logger"
	"LoginAPI/models"
	"LoginAPI/repository"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	cognitoInterface "github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
)

type IAuthService interface {
	CreateCognitoUser(user *models.LoginDetails) (*models.Token, error)
	LoginUser(user *models.LoginDetails) (*models.Token, error)
}

type authService struct {
	repo    repository.IAuthRepository
	env     config.IConfig
	cognito cognitoInterface.CognitoIdentityProviderAPI
	logObj  logger.ILogger
}

var authServiceObj IAuthService

func InitAuthService(repoObj repository.IAuthRepository,
	env config.IConfig,
	cognito cognitoInterface.CognitoIdentityProviderAPI,
	logObj logger.ILogger) IAuthService {
	if authServiceObj == nil {
		authServiceObj = &authService{
			repo:    repoObj,
			env:     env,
			cognito: cognito,
			logObj:  logObj,
		}
	}
	return authServiceObj
}

func (s *authService) CreateCognitoUser(user *models.LoginDetails) (*models.Token, error) {

	//todo:  add statergies to use login type
	// Create a new user
	createUserInput := &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId:        aws.String(s.env.GetUserPoolId()),
		Username:          aws.String(user.UserName),
		MessageAction:     aws.String(localAwsPkg.MESSAGE_ACTION_SUPRESS),
		TemporaryPassword: aws.String(user.Password), // Provide a temporary password
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(user.Email),
			},
		},
	}

	userData, err := s.cognito.AdminCreateUser(createUserInput)
	if err != nil {
		s.logObj.Printf("error in creating user (%s), error: %+v", user.UserName, err)
		return nil, s.processCreateUserError(err)
	}
	s.logObj.Printf("User (%s) created successfully \n", *userData.User.Username)

	if *userData.User.UserStatus == localAwsPkg.USER_STATUS_FORCE_CHANGE_PASSWORD {
		auth, err := s.LoginUser(user)
		if auth == nil && err != nil {
			return nil, err
		}
		updatePasswordError := s.respondToNewPasswordChallenge(auth.Session, user.UserName, user.Password)
		if updatePasswordError != nil {
			return nil, &customerrors.UpdatePasswordException{}
		}
	}

	addUserToGroupInput := &cognitoidentityprovider.AdminAddUserToGroupInput{
		GroupName:  aws.String(user.Role),
		UserPoolId: aws.String(s.env.GetUserPoolId()),
		Username:   aws.String(*userData.User.Username),
	}

	_, err = s.cognito.AdminAddUserToGroup(addUserToGroupInput)
	if err != nil {
		s.logObj.Printf("error in adding user (%s) to group, error: %+v", user.UserName, err)
		return nil, s.processCreateUserError(err)
	}

	return s.LoginUser(user)
}

func (s *authService) LoginUser(loginDetails *models.LoginDetails) (*models.Token, error) {
	// Authenticate the user
	authInput := &cognitoidentityprovider.AdminInitiateAuthInput{
		UserPoolId: aws.String(s.env.GetUserPoolId()),
		ClientId:   aws.String(s.env.GetClientId()),
		AuthFlow:   aws.String(localAwsPkg.AUTH_FLOW_ADMIN_USER_PASSWORD_AUTH),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(loginDetails.UserName),
			"PASSWORD": aws.String(loginDetails.Password), // Provide the temporary password here
		},
	}

	authOutput, err := s.cognito.AdminInitiateAuth(authInput)
	if err != nil {
		s.logObj.Printf("error in login user %+v\n", err)
		return nil, &customerrors.InvalidCredentialException{}
	}

	response := new(models.Token)
	if authOutput.Session != nil {
		response.Session = *authOutput.Session
	}

	if authOutput.AuthenticationResult != nil {
		response.Token = *authOutput.AuthenticationResult.AccessToken
		response.RefreshToken = *authOutput.AuthenticationResult.RefreshToken
		response.TokenType = *authOutput.AuthenticationResult.TokenType
		response.Expires = *authOutput.AuthenticationResult.ExpiresIn
	}
	if authOutput.ChallengeName != nil && *authOutput.ChallengeName == localAwsPkg.COGNITO_CHALLANGE_NAME_NEW_PASSWORD_REQUIRED {
		return response, &customerrors.PasswordExpireException{}
	}

	return response, nil
}

func (s *authService) processCreateUserError(err error) error {

	if _, ok := err.(*cognitoidentityprovider.UsernameExistsException); ok {
		return &customerrors.UsernameExistsException{}
	}
	if _, ok := err.(*cognitoidentityprovider.InvalidPasswordException); ok {
		return &customerrors.InvalidPasswordException{}
	}
	return &customerrors.CreateUserException{}
}

func (s *authService) respondToNewPasswordChallenge(session, username, newPassword string) error {
	challengeResponse := &cognitoidentityprovider.AdminRespondToAuthChallengeInput{
		UserPoolId:    aws.String(s.env.GetUserPoolId()),
		ClientId:      aws.String(s.env.GetClientId()),
		ChallengeName: aws.String(localAwsPkg.COGNITO_CHALLANGE_NAME_NEW_PASSWORD_REQUIRED),
		Session:       aws.String(session),
		ChallengeResponses: map[string]*string{
			"USERNAME":     aws.String(username),
			"NEW_PASSWORD": aws.String(newPassword),
		},
	}

	_, err := s.cognito.AdminRespondToAuthChallenge(challengeResponse)
	if err != nil {
		s.logObj.Printf("Error when updating password %+v\n", err)
		return err
	}

	s.logObj.Printf("New password set successfully\n")
	return nil
}
