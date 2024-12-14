package utils

import (
	"encoding/json"
	"net/http"
)

// Шаблон ответа сервера на запросы
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Шаблон сообщения в ответе сервера на запросы
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}
