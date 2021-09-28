package main

import (
	controllers "goapm/app"

	_ "fmt"
	_ "net/http"
	_ "strconv"
	_ "strings"
)

var appName = "goapm"

var (
	Build        string
	Commit       string
	BuildTime    string
	Version      string
)

func main() {

	controllers.InitControllers(appName, Version, Build, Commit, BuildTime)

}
