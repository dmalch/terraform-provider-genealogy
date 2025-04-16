package geni

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"

	"github.com/allegro/bigcache/v3"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (c *Client) GetDocumentFromCache(ctx context.Context, documentId string) (*DocumentResponse, error) {
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
			document, err := c.GetDocument(ctx, documentId)
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

	var document DocumentResponse
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
		documents, err := c.GetUploadedDocuments(ctx, 1)
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
			documents, err = c.GetUploadedDocuments(ctx, documents.Page+1)
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
