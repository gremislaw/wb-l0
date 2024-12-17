package rest

import (
	"go.uber.org/zap"
	"html/template"
	"net/http"
	. "order_service/internal/logger"
	"order_service/internal/models"
)

func renderOrder(w http.ResponseWriter, order *models.Order) {
	tmpl, err := template.ParseFiles("/Users/gremislaw/Projects/Golang/wb-l0/templates/index.html")
	if err != nil {
		errMsg := "templade load error"
		Logger.Error(errMsg, zap.Error(err))
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, order); err != nil {
		errMsg := "template execute error"
		Logger.Error(errMsg, zap.Error(err))
		http.Error(w, errMsg, http.StatusInternalServerError)
	}
}
