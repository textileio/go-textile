package controllers_test

import (
	"net/http"
	"os"
)

var client = &http.Client{}
var apiURL string
var refKey = os.Getenv("REF_KEY")
