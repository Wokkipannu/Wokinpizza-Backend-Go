package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"wp-backend/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

var conf = &oauth2.Config{
	RedirectURL:  config.Config("API_URL") + "auth/callback",
	ClientID:     config.Config("DISCORD_CLIENT_ID"),
	ClientSecret: config.Config("DISCORD_CLIENT_SECRET"),
	Scopes:       []string{discord.ScopeIdentify, discord.ScopeEmail, discord.ScopeGuilds},
	Endpoint:     discord.Endpoint,
}

func Callback(c *fiber.Ctx) error {
	type DiscordUser struct {
		Id            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Avatar        string `json:"avatar"`
		Public_flags  int    `json:"public_flags"`
		Flags         int    `json:"flags"`
		Banner        string `json:"banner"`
		Banner_color  string `json:"banner_color"`
		Accent_color  string `json:"accent_color"`
		Locale        string `json:"locale"`
		Mfa_enabled   bool   `json:"mfa_enabled"`
		Premium_type  int    `json:"premium_type"`
		Email         string `json:"email"`
		Verified      bool   `json:"verified"`
	}
	var user DiscordUser

	if c.Query("state") != config.Config("STATE") {
		return c.SendStatus(fiber.StatusForbidden)
	}

	token, err := conf.Exchange(context.Background(), c.Query("code"))
	if err != nil {
		log.Printf("Error when authenticating with discord")
		return c.Status(500).JSON(fiber.Map{"message": "Error when authenticating with discord"})
	}

	// Get user info
	res, err := conf.Client(context.Background(), token).Get("https://discord.com/api/users/@me")
	if err != nil {
		log.Printf("Error when fetching user data from discord")
		return c.Status(500).JSON(fiber.Map{"message": "Error when fetching user data from discord"})
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&user)

	jwtToken := jwt.New(jwt.SigningMethodHS256)
	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["id"] = user.Id
	claims["username"] = user.Username
	claims["email"] = user.Email
	// For now I only need one admin account, so this solution works
	if user.Email == config.Config("ADMIN_EMAIL") {
		claims["access_level"] = "admin"
	} else {
		claims["access_level"] = "user"
	}
	claims["avatar"] = "https://cdn.discordapp.com/avatars/" + user.Id + "/" + user.Avatar + ".png"
	claims["token"] = token
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := jwtToken.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = t
	cookie.HTTPOnly = true
	cookie.Expires = time.Now().Add(time.Hour * 72)

	c.Cookie(cookie)
	return c.Redirect("/login")
	// return c.JSON(fiber.Map{"message": "Login success"})
}

func Login(c *fiber.Ctx) error {
	return c.Redirect(conf.AuthCodeURL(config.Config("STATE")))
}

func Logout(c *fiber.Ctx) error {
	// This currently does not work
	// c.ClearCookie("auth_token")

	cookie := new(fiber.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = "deleted"
	cookie.HTTPOnly = true
	cookie.Expires = time.Now().Add(-3 * time.Second)
	c.Cookie(cookie)

	return c.JSON(fiber.Map{"message": "Logged out"})
}

func GetUser(c *fiber.Ctx) error {
	cookie := c.Cookies("auth_token")

	claims, err := ValidateToken(cookie)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Token validation error"})
	}

	return c.JSON(fiber.Map{"message": "User valid", "data": claims})
}

// func GetGuilds(c *fiber.Ctx) error {
// 	type DiscordGuild struct {
// 		Id          string            `json:"id"`
// 		Name        string            `json:"name"`
// 		Icon        string            `json:"icon"`
// 		Owner       bool              `json:"owner"`
// 		Permissions string            `json:"permissions"`
// 		Features    map[string]string `json:"features"`
// 	}
// 	var guilds []DiscordGuild

// 	cookie := c.Cookies("auth_token")
// 	claims, err := ValidateToken(cookie)
// 	if err != nil {
// 		return c.JSON(fiber.Map{"status": "error", "message": err})
// 	}

// 	// Probably exists a better way of doing this, but for now this works
// 	t, err := json.Marshal(claims["token"])
// 	if err != nil {
// 		return c.JSON(fiber.Map{"status": "error", "message": "Error when converting tokens"})
// 	}

// 	var token *oauth2.Token
// 	json.Unmarshal([]byte(t), &token)

// 	res, err := conf.Client(context.Background(), token).Get("https://discord.com/api/users/@me/guilds")
// 	if err != nil {
// 		log.Printf("Error when fetching guild data from user")
// 		return c.JSON(fiber.Map{"status": "error", "message": "Error when fetching guild data from user"})
// 	}
// 	defer res.Body.Close()

// 	json.NewDecoder(res.Body).Decode(&guilds)

// 	return c.JSON(fiber.Map{"status": "success", "message": "Fetched guilds", "data": guilds})
// }

func ValidateToken(t string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Config("SECRET")), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error when validating token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid or expired token")
	}
}

func AuthLevelAdmin(c *fiber.Ctx) error {
	cookie := c.Cookies("auth_token")

	claims, err := ValidateToken(cookie)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Token validation error"})
	}

	if claims["access_level"] == "admin" {
		return c.Next()
	} else {
		return c.SendStatus(fiber.StatusForbidden)
	}
}
