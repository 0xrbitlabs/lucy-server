package middlewares

import "lucy/repositories"

type Middleware struct {
	userRepo repositories.UserRepo
}
