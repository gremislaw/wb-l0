package rest

import (
	"net/http"
	. "order_service/internal/logger"

	"go.uber.org/zap"
)

// endpoint для получения order по id
func (s *RestService) GetOrder(w http.ResponseWriter, r *http.Request) {
	Logger.Info("got request: GetRelease")
	// получение id
	idStr := r.URL.Path[len("/order/"):]
	Logger.Info("got id", zap.String("id", idStr))

	// запрос к бд для получения order
	order, err := s.cache.GetOrder(idStr)

	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	renderOrder(w, order)
	Logger.Info("release got", zap.String("id", idStr))
}
