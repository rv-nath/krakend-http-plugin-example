package main

import (
	"context"
	"errors"
	"fmt"
	"html"
	"net/http"
)

var (
	pluginName        = "krakend-plugin"
	HandlerRegisterer = registerer(pluginName)
)

type registerer string

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
),
) {
	logger.Debug("RegisterHandlers called.")
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(_ context.Context, extra map[string]interface{}, h http.Handler) (http.Handler, error) {
	logger.Debug("registerHandlers called with extra config:", extra)
	// If the plugin requires some configuration, it should be under the name of the plugin. E.g.:
	/*
	   "extra_config":{
	       "plugin/http-server":{
	           "name":["krakend-server-example"],
	           "krakend-server-example":{
	               "path": "/some-path"
	           }
	       }
	   }
	*/

	// The config variable contains all the keys defined in the configuration.
	// If the key doesn't exist or is not a map, the plugin returns an error and the default handler.
	config, ok := extra[pluginName].(map[string]interface{})
	if !ok {
		return h, errors.New("configuration not found")
	}

	path, _ := config["path"].(string)
	logger.Debug("Hijack path configured as:", path)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Incoming request path:", r.URL.Path)

		// If hijack path matches the URL in the request, then hijack it.
		if r.URL.Path != path {
			h.ServeHTTP(w, r)
			return
		}
		logger.Info("Hijacking request:", r.URL.Path)
		fmt.Fprintf(w, "Lowjack %q", html.EscapeString((r.URL.Path)))
		logger.Debug("request:", html.EscapeString(r.URL.Path))
	}), nil
}

func main() {}

// This logger is replaced by the RegisterLogger method to load the one from krakenD.
var logger Logger = noopLogger{}

func (registerer) RegisterLogger(v interface{}) {
	l, ok := v.(Logger)
	if !ok {
		fmt.Println("Failed to load logger.")
		return
	}
	logger = l
	logger.Debug(fmt.Sprintf("[PLUGIN: %s] Logger loaded", HandlerRegisterer))
}

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Critical(v ...interface{})
	Fatal(v ...interface{})
}

// Empty logger implementation
type noopLogger struct{}

func (n noopLogger) Debug(_ ...interface{})    {}
func (n noopLogger) Info(_ ...interface{})     {}
func (n noopLogger) Warning(_ ...interface{})  {}
func (n noopLogger) Error(_ ...interface{})    {}
func (n noopLogger) Critical(_ ...interface{}) {}
func (n noopLogger) Fatal(_ ...interface{})    {}
