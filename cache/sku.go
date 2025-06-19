package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/RohitGupta-omniful/IMS/models"
)

func CacheSKU(ctx context.Context, sku models.SKU, ttl time.Duration) error {
	key := fmt.Sprintf("sku:%s", sku.ID.String())
	return SetJSON(ctx, key, sku, ttl)
}

func GetCachedSKU(ctx context.Context, skuID string) (*models.SKU, error) {
	key := fmt.Sprintf("sku:%s", skuID)
	var sku models.SKU
	if err := GetJSON(ctx, key, &sku); err != nil {
		return nil, err
	}
	return &sku, nil
}

func DeleteSKUCache(ctx context.Context, skuID string) error {
	key := fmt.Sprintf("sku:%s", skuID)
	return Del(ctx, key)
}
