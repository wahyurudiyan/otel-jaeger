package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/wahyurudiyan/otel-jaeger/pkg/random"
	"go.opentelemetry.io/otel"
)

var (
	tracer = otel.Tracer("WebServer-Otel-Jaeger")
)

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "GetUser")
	defer span.End()

	user := struct {
		Name     string
		Email    string
		Password string
	}{
		Name:     "John Doe",
		Email:    "john@email.com",
		Password: "Super5ecr3t!",
	}
	blob, _ := json.Marshal(&user)

	sleepDuration := time.Duration(time.Millisecond * time.Duration(random.GenerateRandNum()))
	time.Sleep(sleepDuration)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(blob)
}

func Router(router chi.Router) {
	router.Get("/user", getUserHandler)
}
