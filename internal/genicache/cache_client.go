package genicache

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"sync"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/dmalch/terraform-provider-genealogy/internal/geni"
	"github.com/dmalch/terraform-provider-genealogy/internal/genibatch"
)

const (
	maxProfilesPerPage  = 50
	maxDocumentsPerPage = 50
)

type Client struct {
	client                   *geni.Client
	batchClient              *genibatch.Client
	cache                    *bigcache.BigCache
	profileCacheLock         sync.Mutex
	profileCacheInitialized  bool
	documentCacheLock        sync.Mutex
	documentCacheInitialized bool
}

func NewClient(client *geni.Client, batchClient *genibatch.Client) *Client {
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		tflog.Error(context.Background(), "Error creating cache", map[string]interface{}{"error": err})
		panic(err)
	}

	return &Client{
		client:      client,
		batchClient: batchClient,
		cache:       cache,
	}
}

func (c *Client) GetProfile(ctx context.Context, profileId string) (*geni.ProfileResponse, error) {
	if err := c.initProfileCache(ctx); err != nil {
		tflog.Error(ctx, "Error initializing cache", map[string]interface{}{"error": err})
		return nil, err
	}

	// Retrieve the profile from the cache
	data, err := c.cache.Get(profileId)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			tflog.Debug(ctx, "Profile not found in cache", map[string]interface{}{"profileId": profileId})

			// If the profile is not found in the cache, retrieve it using GetProfile method
			profile, err := c.batchClient.GetProfile(ctx, profileId)
			if err != nil {
				tflog.Error(ctx, "Error retrieving profile", map[string]interface{}{"error": err})
				return nil, err
			}

			// Store the retrieved profile in the cache
			if err := c.storeInCache(ctx, profile.Id, *profile); err != nil {
				tflog.Error(ctx, "Error storing profile in cache", map[string]interface{}{"error": err})
				return nil, err
			}

			return profile, nil
		}

		tflog.Error(ctx, "Error retrieving profile from cache", map[string]interface{}{"error": err})
		return nil, err
	}

	var profile geni.ProfileResponse
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&profile); err != nil {
		tflog.Error(ctx, "Error decoding profile", map[string]interface{}{"error": err})
		return nil, err
	}

	return &profile, nil
}

func (c *Client) initProfileCache(ctx context.Context) error {
	c.profileCacheLock.Lock()
	defer c.profileCacheLock.Unlock()

	// If the cache is empty, retrieve all managed profiles
	if !c.profileCacheInitialized {
		// Retrieve the first page of managed profiles using the API
		profiles, err := c.client.GetManagedProfiles(ctx, 1)
		if err != nil {
			tflog.Error(ctx, "Error retrieving managed profiles", map[string]interface{}{"error": err})
			return err
		}

		for _, profile := range profiles.Results {
			if err := c.storeInCache(ctx, profile.Id, profile); err != nil {
				tflog.Error(ctx, "Error storing profile in cache", map[string]interface{}{"error": err})
				return err
			}
		}

		// Retrieve all managed profiles using the API, run up to 200 times
		for i := 0; i < 200 && len(profiles.Results) == maxProfilesPerPage; i++ {
			profiles, err = c.client.GetManagedProfiles(ctx, profiles.Page+1)
			if err != nil {
				tflog.Error(ctx, "Error retrieving managed profiles", map[string]interface{}{"error": err})
				return err
			}

			for _, profile := range profiles.Results {
				if err := c.storeInCache(ctx, profile.Id, profile); err != nil {
					tflog.Error(ctx, "Error storing profile in cache", map[string]interface{}{"error": err})
					return err
				}
			}
		}

		c.profileCacheInitialized = true
	}
	return nil
}

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

func (c *Client) GetDocument(ctx context.Context, documentId string) (*geni.DocumentResponse, error) {
	if err := c.initDocumentCache(ctx); err != nil {
		tflog.Error(ctx, "Error initializing cache", map[string]interface{}{"error": err})
		return nil, err
	}

	// Retrieve the document from the cache
	data, err := c.cache.Get(documentId)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			tflog.Debug(ctx, "Document not found in cache", map[string]interface{}{"documentId": documentId})

			// If the document is not found in the cache, retrieve it using GetDocument method
			document, err := c.client.GetDocument(ctx, documentId)
			if err != nil {
				tflog.Error(ctx, "Error retrieving document", map[string]interface{}{"error": err})
				return nil, err
			}

			// Store the retrieved document in the cache
			if err := c.storeInCache(ctx, document.Id, *document); err != nil {
				tflog.Error(ctx, "Error storing document in cache", map[string]interface{}{"error": err})
				return nil, err
			}

			return document, nil
		}

		tflog.Error(ctx, "Error retrieving document from cache", map[string]interface{}{"error": err})
		return nil, err
	}

	var document geni.DocumentResponse
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&document); err != nil {
		tflog.Error(ctx, "Error decoding document", map[string]interface{}{"error": err})
		return nil, err
	}

	return &document, nil
}

func (c *Client) initDocumentCache(ctx context.Context) error {
	c.documentCacheLock.Lock()
	defer c.documentCacheLock.Unlock()

	// If the cache is empty, retrieve all managed documents
	if !c.documentCacheInitialized {
		// Retrieve the first page of managed documents using the API
		documents, err := c.client.GetUploadedDocuments(ctx, 1)
		if err != nil {
			tflog.Error(ctx, "Error retrieving managed documents", map[string]interface{}{"error": err})
			return err
		}

		for _, document := range documents.Results {
			if err := c.storeInCache(ctx, document.Id, document); err != nil {
				tflog.Error(ctx, "Error storing document in cache", map[string]interface{}{"error": err})
				return err
			}
		}

		// Retrieve all managed documents using the API, run up to 200 times
		// 50 is the maximum number of documents per page
		for i := 0; i < 200 && len(documents.Results) == maxDocumentsPerPage; i++ {
			documents, err = c.client.GetUploadedDocuments(ctx, documents.Page+1)
			if err != nil {
				tflog.Error(ctx, "Error retrieving managed documents", map[string]interface{}{"error": err})
				return err
			}

			for _, document := range documents.Results {
				if err := c.storeInCache(ctx, document.Id, document); err != nil {
					tflog.Error(ctx, "Error storing document in cache", map[string]interface{}{"error": err})
					return err
				}
			}
		}

		c.documentCacheInitialized = true
	}
	return nil
}
