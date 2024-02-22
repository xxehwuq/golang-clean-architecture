package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xxehwuq/go-clean-architecture/internal/usecase"
)

type usersRoutes struct {
	usecase usecase.UserUsecase
}

func initUsersRoutes(handler fiber.Router, usecase usecase.UserUsecase) {
	r := &usersRoutes{
		usecase: usecase,
	}

	users := handler.Group("/users")
	{
		users.Post("/sign-up", r.SignUp)
		users.Post("/sign-in", r.SignIn)
		users.Post("/refresh", r.Refresh)
	}
}

type usersSignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type usersSignUpResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r *usersRoutes) SignUp(ctx *fiber.Ctx) error {
	var req usersSignUpRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	userTokens, err := r.usecase.SignUp(ctx.Context(), usecase.UserSignUpInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return ctx.SendString(err.Error()) // TODO: error handling
	}

	return ctx.Status(fiber.StatusOK).JSON(usersSignUpResponse{
		AccessToken:  userTokens.AccessToken,
		RefreshToken: userTokens.RefreshToken,
	})
}

type usersSignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type usersSignInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r *usersRoutes) SignIn(ctx *fiber.Ctx) error {
	var req usersSignInRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	userTokens, err := r.usecase.SignIn(ctx.Context(), usecase.UserSignInInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return ctx.SendString(err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(usersSignInResponse{
		AccessToken:  userTokens.AccessToken,
		RefreshToken: userTokens.RefreshToken,
	})
}

type usersRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type usersRefreshResponse struct {
	Token string `json:"token"`
}

func (r *usersRoutes) Refresh(ctx *fiber.Ctx) error {
	var req usersRefreshRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	userTokens, err := r.usecase.RefreshTokens(ctx.Context(), req.RefreshToken)
	if err != nil {
		return ctx.SendString(err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(usersSignInResponse{
		AccessToken:  userTokens.AccessToken,
		RefreshToken: userTokens.RefreshToken,
	})
}
