package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/internal/storage"
	"time"
	//"github.com/softlayer/softlayer-go/sl"
	"golang.org/x/crypto/bcrypt"

	//"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



type Auth struct {
	usrSaver UserSaver
	usrProvider UserProvider
	log *slog.Logger
	appProvider AppProvider
	tokenTTL time.Duration
}


type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}


type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}


var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func New(
	log *slog.Logger, 
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		usrSaver : userSaver,
		usrProvider : userProvider,
		log : log,
		appProvider : appProvider,
		tokenTTL : tokenTTL,
	}
}	


func (a *Auth) Login(
	ctx context.Context,
	email string, 
	password string, 
	appID int,
) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op: ", op),
		slog.String("email: ", email),
	)


	log.Info("attempting to login user")

	user, err := a.usrProvider.User(ctx, email)

	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn("user not found", 	)
			return "", fmt.Errorf("%s %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", status.Error(codes.Internal, ""))

		return "", fmt.Errorf("%s %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", status.Error(codes.Internal, ""))

		return "", fmt.Errorf("%s %w", op, ErrInvalidCredentials)
	}


	app, err := a.appProvider.App(ctx, appID)

	if err != nil {
		return "", fmt.Errorf("%s %w", op, err)
	}

	log.Info("user logged successfully")


	token, err := jwt.NewToken(user, app, a.tokenTTL)

	if err != nil {
		a.log.Error("failed to generate token", err) //status.Error(codes.Internal, ""))

		return "", fmt.Errorf("%s %w", op, err)
	}

	return token, nil

	

}


func (a *Auth) RegisterNewUser(ctx context.Context, email string, pass string) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op: ", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	if err != nil {
		log.Error("failed to generate password hash", status.Error(codes.Internal, ""))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)

	if err != nil {
		log.Error("failed to save user", err) //status.Errorf(codes.Internal, ""))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")
	
	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"


	log := a.log.With(
		slog.String("op: ", op),
		slog.Int64("userID: ", userID),
	)
	log.Info("starting checking")


	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)

	if err != nil {
		log.Error("failed to check", err) //status.Error(codes.Internal, ""))
		return false, fmt.Errorf("%s %w", op, err)
	}

	log.Info("checked is admin or no", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil

}

