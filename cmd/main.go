package main

import (
	"context"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/qeunasd/coniven/server"
	"github.com/qeunasd/coniven/services/category_service"
	"github.com/qeunasd/coniven/storage"
)

func main() {
	println("hello, myself")

	ctx := context.Background()
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env: %v", err)
	}

	db, err := storage.NewPostgres(ctx, os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatalf("error connecting to postgres: %v", err)
	}

	repository := storage.NewRepository(db, nil)
	if err := repository.RunMigration(ctx); err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()
	templates, err := ParseTemplate("./templates")
	if err != nil {
		log.Fatalf("parsing template: %s", err)
	}

	categoryService := category_service.NewCategoryService(repository)

	log.Println("listening to server at localhost:8080")
	srv := server.NewServer(router, templates, categoryService)
	if err := srv.Run(); err != nil {
		log.Fatalf("error listening to server: %v", err)
	}
}

func ParseTemplate(dir string) (*template.Template, error) {
	cleanRoot := filepath.Clean(dir)
	tmpl := template.New("")

	tmpl.Funcs(template.FuncMap{
		"embed": func(name string, data any) template.HTML {
			var output strings.Builder
			if err := tmpl.ExecuteTemplate(&output, name, data); err != nil {
				log.Fatalf("embedding template: %s", err)
			}
			return template.HTML(output.String())
		},
	})

	err := filepath.Walk(cleanRoot, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && (strings.HasSuffix(path, ".tmpl")) {
			if err != nil {
				return err
			}

			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			name := path[len(cleanRoot)+1:]
			_, err = tmpl.New(name).Parse(string(b))
			if err != nil {
				return err
			}
		}
		return nil
	})

	return tmpl, err
}
