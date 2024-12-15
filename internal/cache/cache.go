package cache

import "order_service/internal/models"

func (c *CacheWrapper) FindAll() []*models.Order {
	ordersArr := make([]*models.Order, 0, len(c.orders))

	for _, order := range c.orders {
		ordersArr = append(ordersArr, order)
	}

	return ordersArr
}

func (c *CacheWrapper) FindOrder(uid string) (*models.Order, bool) {
	order, ok := c.orders[uid]
	return order, ok
}

func (c *CacheWrapper) AddOrder(order *models.Order) {
	c.mu.Lock()
	err := c.Queries.CreateOrder(c.Ctx, *order)
	if err != nil {
		panic(err)
	}
	c.orders[order.OrderUID] = order
	c.mu.Unlock()
}
