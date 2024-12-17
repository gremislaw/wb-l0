package cache

import (
	"errors"
	"order_service/internal/models"
	"sync"
)

type CacheMap struct {
	m  map[string]*models.Order
	mu sync.Mutex
}

func NewCacheMap() *CacheMap {

	return &CacheMap{
		m: make(map[string]*models.Order),
	}
}

func (c *CacheMap) Set(order *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.m[order.OrderUID] = order
}

func (c *CacheMap) Get(key string) (*models.Order, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.m[key]
	if !ok {
		return nil, errors.New("order not found")
	}

	return val, nil
}

func (c *CacheMap) GetAll() ([]*models.Order, error) {
	ordersArr := make([]*models.Order, 0, len(c.m))
	for _, order := range c.m {
		ordersArr = append(ordersArr, order)
	}

	return ordersArr, nil
}

/*func NewCacheWrapper(dbRepo *db_repository.Repository, ordersMap map[string]*models.Order, mu *sync.Mutex) *CacheWrapper {
	return &CacheWrapper{dbRepo: dbRepo, orders: ordersMap, mu: mu}
}

func CreateCache(db *sql.DB) *CacheWrapper {
	// создание обертки для кэша
	ctx := context.Background()
	dbRepo := db_repository.NewRepository(db, ctx)
	orders, err := dbRepo.GetOrders()
	if err != nil {
		panic(err)
	}
	ordersMap := make(map[string]*models.Order)

	for _, order := range orders {
		ordersMap[order.OrderUID] = &order
	}

	cacheWrapper := NewCacheWrapper(dbRepo, ordersMap, &sync.Mutex{})

	return cacheWrapper
}*/
