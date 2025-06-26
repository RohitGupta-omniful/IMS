package cache

import (
	"context"
	"time"

	"github.com/RohitGupta-omniful/IMS/models"
)

// CacheHub caches a hub object
func CacheHub(ctx context.Context, hub models.Hub, ttl time.Duration) error {
	key := "hub:" + hub.ID.String()
	return SetJSON(ctx, key, hub, ttl)
}

// GetCachedHub fetches hub from Redis
func GetCachedHub(ctx context.Context, hubID string) (*models.Hub, error) {
	key := "hub:" + hubID
	var hub models.Hub
	if err := GetJSON(ctx, key, &hub); err != nil {
		return nil, err
	}
	return &hub, nil
}

// DeleteHubCache removes cached hub
func DeleteHubCache(ctx context.Context, hubID string) error {
	key := "hub:" + hubID
	return Del(ctx, key)
}
