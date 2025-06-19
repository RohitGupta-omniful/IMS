package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/RohitGupta-omniful/IMS/models"
)

// CacheHub caches a hub object
func CacheHub(ctx context.Context, hub models.Hub, ttl time.Duration) error {
	key := fmt.Sprintf("hub:%s", hub.ID.String())
	return SetJSON(ctx, key, hub, ttl)
}

// GetCachedHub fetches hub from Redis
func GetCachedHub(ctx context.Context, hubID string) (*models.Hub, error) {
	key := fmt.Sprintf("hub:%s", hubID)
	var hub models.Hub
	if err := GetJSON(ctx, key, &hub); err != nil {
		return nil, err
	}
	return &hub, nil
}

// DeleteHubCache removes cached hub
func DeleteHubCache(ctx context.Context, hubID string) error {
	key := fmt.Sprintf("hub:%s", hubID)
	return Del(ctx, key)
}
