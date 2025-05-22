package main

import (
	"context"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/qeunasd/coniven/server"
	"github.com/qeunasd/coniven/services"
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

	templates, err := ParseTemplate("./templates")
	if err != nil {
		log.Fatalf("parsing template: %s", err)
	}

	categoryService := services.NewCategoryService(repository)
	locationService := services.NewLocationService(repository)
	roomService := services.NewRoomService(repository)
	itemService := services.NewRoomService(repository)

	log.Println("listening to server at localhost:8080")
	srv := server.NewServer(templates, categoryService, locationService, roomService, itemService)

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
		"pageRange": func(current, total, max int) []int {
			if total <= max {
				r := make([]int, total)
				for i := range r {
					r[i] = i + 1
				}
				return r
			}

			var pages []int
			half := max / 2
			start := current - half
			end := current + half

			if start < 1 {
				start = 1
				end = max
			} else if end > total {
				end = total
				start = total - max + 1
			}

			for i := start; i <= end; i++ {
				pages = append(pages, i)
			}

			return pages
		},
		"sub": func(x, y int) int {
			return x - y
		},
		"parseTime": func(date time.Time) string {
			return date.Format("02-01-2006 15:04:05")
		},
		"uidStr": func(id uuid.UUID) string {
			return id.String()
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
