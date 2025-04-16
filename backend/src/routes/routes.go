package routes

import (
	"interrupted-export/src/controllers"
)

var ApiRoutes = []Route{
	{Method: "POST", Path: "/login", Handler: controllers.Login},
}
