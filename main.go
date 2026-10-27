package main

import (
	"os"
	"time"
	"uncomplicated-docs/core"

	"github.com/Elagoht/collage/pkg/collage"
)

func main() {
	collage.Main(newApp)
}

func newApp(env collage.Env) (*collage.App, error) {
	store := core.NewFileStore("docs")
	err := store.Load()
	if err != nil {
		return nil, err
	}

	pluginConfig, err := env.PluginConfig(files)
	if err != nil {
		return nil, err
	}

	app, err := collage.New(&collage.Config{
		DevMode: env.Dev,
		BaseURL: os.Getenv("BASE_URL"),
		Server:  collage.ServerConfig{Host: env.Host, Port: env.Port},
		Template: collage.TemplateConfig{
			FS:        files,
			Root:      "templates",
			Extension: ".html",
		},
		Cache: collage.CacheConfig{
			Enabled:    true,
			Type:       "disk",
			Dir:        env.CacheDir,
			DefaultTTL: 5 * time.Minute,
		},
		PluginConfig: pluginConfig,
		Plugins:      plugins(),
		Security: collage.SecurityConfig{
			CSRFKey: []byte(os.Getenv("COLLAGE_CSRF_KEY")),
		},
	})
	if err != nil {
		return nil, err
	}

	if err := register(app, store); err != nil {
		return nil, err
	}

	return app, app.Mount("/static/", env.Static(files, "static"))
}
