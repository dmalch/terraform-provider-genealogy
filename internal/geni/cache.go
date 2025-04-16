package geni

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (c *Client) storeInCache(ctx context.Context, id string, object any) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(object); err != nil {
		tflog.Error(ctx, "Error encoding object", map[string]interface{}{"error": err})
		return err
	}
	if err := c.cache.Set(id, buf.Bytes()); err != nil {
		tflog.Error(ctx, "Error setting object in cache", map[string]interface{}{"error": err})
		return err
	}
	return nil
}
