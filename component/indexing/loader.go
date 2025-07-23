package indexing

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
)

// NewLoader 支持目录和单文件
func NewLoader(ctx context.Context) (document.Loader, error) {
	fileLoader, err := file.NewFileLoader(ctx, &file.FileLoaderConfig{})
	if err != nil {
		return nil, err
	}
	return &BatchFileLoader{
		Single: fileLoader,
	}, nil
}

// BatchFileLoader 支持批量加载目录下所有文件
type BatchFileLoader struct {
	Single *file.FileLoader
}

// Load 实现 document.Loader 接口，支持批量加载文件
func (b *BatchFileLoader) Load(ctx context.Context, src document.Source, opts ...document.LoaderOption) ([]*schema.Document, error) {
	var docs []*schema.Document
	fileInfo, err := os.Stat(src.URI)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		// 遍历目录下所有文件
		err := filepath.WalkDir(src.URI, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				subDocs, err := b.Single.Load(ctx, document.Source{URI: path})
				if err != nil {
					return err
				}
				docs = append(docs, subDocs...)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		return b.Single.Load(ctx, src)
	}
	return docs, nil
}

func (b *BatchFileLoader) GetType() string {
	return "BatchFileLoader"
}

func (b *BatchFileLoader) IsCallbacksEnabled() bool {
	return true
}
