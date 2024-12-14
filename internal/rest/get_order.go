package rest

import (
	"net/http"
	. "order_service/internal/logger"
	"order_service/internal/utils"
	"strconv"

	"go.uber.org/zap"
)

// endpoint для получения order по id
func (s *OrderService) GetOrder(w http.ResponseWriter, r *http.Request) {
	Logger.Info("got request: GetRelease")
	// получение id
	idStr := r.URL.Path[len("/order/"):]
	Logger.Info("got id", zap.String("id", idStr))

	// преобразование полученного id в int, т.к. id в формате string, проверка валидности
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		Logger.Warn("Invalid ID", zap.Int("id", id))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// запрос к бд для получения order
	order, err := s.Queries.GetOrder(s.Ctx, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.Respond(w, map[string]interface{}{"order": order})
	Logger.Info("release got", zap.Int("id", id))
}
