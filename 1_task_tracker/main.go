package main

import (
	"MyMod/repository"
	"MyMod/service"
	"MyMod/ui"
)

func main() {
	repo := repository.RepostoryFile{}
	repo.Init("myTask.txt")
	serv := service.MainService{Repository: &repo}
	ui := ui.UI{Service: &serv}
	ui.Start()
}
