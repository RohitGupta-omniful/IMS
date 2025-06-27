package cache

import (
	"context"
	"time"

	"github.com/RohitGupta-omniful/IMS/models"
)

func CacheHub(ctx context.Context, hub models.Hub, ttl time.Duration) error {
	key := "hub:" + hub.ID.String()
	return SetJSON(ctx, key, hub, ttl)
}

func GetCachedHub(ctx context.Context, hubID string) (*models.Hub, error) {
	key := "hub:" + hubID
	var hub models.Hub
	if err := GetJSON(ctx, key, &hub); err != nil {
		return nil, err
	}
	return &hub, nil
}

func DeleteHubCache(ctx context.Context, hubID string) error {
	key := "hub:" + hubID
	return Del(ctx, key)
}
