package middleware

import "net/http"

type FactoryFunc func(http.Handler) http.Handler
