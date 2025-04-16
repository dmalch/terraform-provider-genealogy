package geni

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"

	"github.com/allegro/bigcache/v3"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (c *Client) GetProfileFromCache(ctx context.Context, profileId string) (*ProfileResponse, error) {
	if err := c.initCache(ctx); err != nil {
		tflog.Error(ctx, "Error initializing cache", map[string]interface{}{"error": err})
		return nil, err
	}

	// Retrieve the profile from the cache
	data, err := c.cache.Get(profileId)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			tflog.Debug(ctx, "Profile not found in cache", map[string]interface{}{"profileId": profileId})

			// If the profile is not found in the cache, retrieve it using GetProfile method
			profile, err := c.GetProfile(ctx, profileId)
			if err != nil {
				tflog.Error(ctx, "Error retrieving profile", map[string]interface{}{"error": err})
				return nil, err
			}

			// Store the retrieved profile in the cache
			if err := c.storeInCache(ctx, *profile); err != nil {
				tflog.Error(ctx, "Error storing profile in cache", map[string]interface{}{"error": err})
				return nil, err
			}

			return profile, nil
		}

		tflog.Error(ctx, "Error retrieving profile from cache", map[string]interface{}{"error": err})
		return nil, err
	}

	var profile ProfileResponse
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&profile); err != nil {
		tflog.Error(ctx, "Error decoding profile", map[string]interface{}{"error": err})
		return nil, err
	}

	return &profile, nil
}

func (c *Client) initCache(ctx context.Context) error {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()

	// If the cache is empty, retrieve all managed profiles
	if c.cache.Len() == 0 {
		// Retrieve the first page of managed profiles using the API
		profiles, err := c.GetManagedProfiles(ctx, 1)
		if err != nil {
			tflog.Error(ctx, "Error retrieving managed profiles", map[string]interface{}{"error": err})
			return err
		}

		// Retrieve all managed profiles using the API, run up to 200 times
		for i := 0; i < 200 && len(profiles.Results) > 0; i++ {
			profiles, err = c.GetManagedProfiles(ctx, profiles.Page+1)
			if err != nil {
				tflog.Error(ctx, "Error retrieving managed profiles", map[string]interface{}{"error": err})
				return err
			}

			for _, profile := range profiles.Results {
				if err := c.storeInCache(ctx, profile); err != nil {
					tflog.Error(ctx, "Error storing profile in cache", map[string]interface{}{"error": err})
					return err
				}
			}
		}
	}
	return nil
}

func (c *Client) storeInCache(ctx context.Context, profile ProfileResponse) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(profile); err != nil {
		tflog.Error(ctx, "Error encoding profile", map[string]interface{}{"error": err})
		return err
	}
	if err := c.cache.Set(profile.Id, buf.Bytes()); err != nil {
		tflog.Error(ctx, "Error setting profile in cache", map[string]interface{}{"error": err})
		return err
	}
	return nil
}
