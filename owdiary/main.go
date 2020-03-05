package main

import (
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	date := time.Now().Format("2006-01-02 Mon")
	t := template.Must(template.ParseFiles(filepath.Join(dir, "message.txt")))
	if err := t.Execute(os.Stdout, date); err != nil {
		log.Fatal(err)
	}
}
