package common

import "github.com/snowlyg/go-tenancy/models"

// Menus 菜单
type Menus struct {
	HomeInfo struct {
		Title string `json:"title"`
		Href  string `json:"href"`
	} `json:"homeInfo"`
	LogoInfo struct {
		Title string `json:"title"`
		Href  string `json:"href"`
		Image string `json:"image"`
	} `json:"logoInfo"`

	MenuInfo []*models.Perm `json:"menuInfo"`
}
