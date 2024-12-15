package rest

import (
	"net/http"
	. "order_service/internal/logger"
	"order_service/internal/utils"

	"go.uber.org/zap"
)

// endpoint для получения order по id
func (s *OrderService) GetOrder(w http.ResponseWriter, r *http.Request) {
	Logger.Info("got request: GetRelease")
	// получение id
	idStr := r.URL.Path[len("/order/"):]
	Logger.Info("got id", zap.String("id", idStr))

	// запрос к бд для получения order
	order, ok := s.Cache.FindOrder(idStr)

	if !ok {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.Respond(w, map[string]interface{}{"order": order})
	Logger.Info("release got", zap.String("id", idStr))
}
