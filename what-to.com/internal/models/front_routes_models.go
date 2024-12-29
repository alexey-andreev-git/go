package models

type (
	RouteData struct {
		RouterLink string `json:"routerLink"`
		IconName   string `json:"iconName"`
		Title      string `json:"title"`
	}
	Route struct {
		Path       string     `json:"path"`
		Component  string     `json:"component,omitempty"`
		Data       *RouteData `json:"data,omitempty"`
		RedirectTo string     `json:"redirectTo,omitempty"`
		PathMatch  string     `json:"pathMatch,omitempty"`
	}
	Routes []Route
)
