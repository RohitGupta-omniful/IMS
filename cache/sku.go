package cache

import (
	"context"
	"time"

	"github.com/RohitGupta-omniful/IMS/models"
)

func CacheSKU(ctx context.Context, sku models.SKU, ttl time.Duration) error {
	key := "sku:" + sku.ID.String()
	return SetJSON(ctx, key, sku, ttl)
}

func GetCachedSKU(ctx context.Context, skuID string) (*models.SKU, error) {
	key := "sku:" + skuID
	var sku models.SKU
	if err := GetJSON(ctx, key, &sku); err != nil {
		return nil, err
	}
	return &sku, nil
}

func DeleteSKUCache(ctx context.Context, skuID string) error {
	key := "sku:" + skuID
	return Del(ctx, key)
}
