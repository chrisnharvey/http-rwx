package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

var username, password, listen, templateFile, output, cmd *string
var templateString string

func main() {
	username = flag.String("username", "", "The username of the user used to authenticate against the web server")
	password = flag.String("password", "", "The password of the user used to authenticate against the web server")
	listen = flag.String("listen", ":80", "The host and port for the web server to listen on")
	templateFile = flag.String("template", "", "The template file that will be parsed when a request is received")
	output = flag.String("output", "", "The output file to write the parsed template to")
	cmd = flag.String("cmd", "", "The command/script to run after each request")

	flag.Parse()

	handler := http.HandlerFunc(handleRequest)
	http.Handle("/", handler)

	if *username == "" || *password == "" || *listen == "" || *templateFile == "" || *output == "" {
		fmt.Println("Missing required flag. See below for usage:")
		flag.Usage()
		os.Exit(1)
	}

	t, err := ioutil.ReadFile(*templateFile)

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	templateString = string(t)

	err = http.ListenAndServe(*listen, nil)

	fmt.Println(err)
	os.Exit(3)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	u, p, valid := r.BasicAuth()

	if !valid {
		fmt.Println("Invalid auth passed")
		w.WriteHeader(401)
		return
	}

	if u != *username || p != *password {
		fmt.Println("Invalid credentials provided:", u, p)
		w.WriteHeader(401)
		return
	}

	err := writeConfig(r.URL.Query())

	if err != nil {
		fmt.Println("Failed to write file")
		fmt.Println(err)
		w.WriteHeader(500)
	}

	fmt.Println("File written")

	if *cmd == "" {
		fmt.Println("No command to execute, skipping")
		w.WriteHeader(200)
		return
	}

	go executeCommand()

	fmt.Println("Responding with 200 OK")
	w.WriteHeader(200)
}

func executeCommand() {
	fmt.Println("Executing command: ", *cmd)

	err := exec.Command("/bin/sh", "-c", *cmd).Run()

	if err != nil {
		fmt.Println("Command execution failed")
		fmt.Println(err)

		return
	}

	fmt.Println("Command executed")
}

func writeConfig(values url.Values) error {
	tpl, err := template.New("template").Funcs(sprig.TxtFuncMap()).Parse(templateString)

	if err != nil {
		panic(err)
	}

	f, err := os.Create(*output)

	if err != nil {
		return err
	}

	err = tpl.Execute(f, values)

	f.Close()

	if err != nil {
		return err
	}

	return nil
}
