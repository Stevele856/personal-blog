package view

import "github.com/letrongvu/blog/internal/model"

type PageData struct{
	CurrentUser *model.User
	Data any
}
