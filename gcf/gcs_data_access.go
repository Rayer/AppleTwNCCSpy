package gcf

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"github.com/Rayer/AppleTwNCCSpy"
	"io/ioutil"
	"strings"
)

type GcsDataAccess struct {
	client *storage.Client
	Bucket string
	Prefix string
}

func NewGcsDataAccess(ctx context.Context, bucket string, prefix string) (AppleProductMonitor.DataAccess, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &GcsDataAccess{
		client: client,
		Bucket: bucket,
		Prefix: prefix,
	}, nil
}

func (g *GcsDataAccess) SaveData(ctx context.Context, data []AppleProductMonitor.Product) error {
	fileStream, _ := json.Marshal(data)
	w := g.client.Bucket(g.Bucket).Object(strings.Join([]string{g.Prefix, "status.json"}, "/")).NewWriter(ctx)
	defer func() {
		_ = w.Close()
	}()
	_, err := w.Write(fileStream)
	return err
}

func (g *GcsDataAccess) LoadData(ctx context.Context) ([]AppleProductMonitor.Product, error) {
	var ret []AppleProductMonitor.Product
	r, err := g.client.Bucket(g.Bucket).Object(strings.Join([]string{g.Prefix, "status.json"}, "/")).NewReader(ctx)
	if err != nil {
		return ret, err
	}
	defer func() {
		_ = r.Close()
	}()

	c, err := ioutil.ReadAll(r)
	if err != nil {
		return ret, err
	}

	err = json.Unmarshal(c, &ret)
	return ret, err
}

//
//func (s SimpleFileDataAccess) SaveData(data []Product) {
//	b, _ := json.Marshal(data)
//	ioutil.WriteFile(s.Filename, b, 0644)
//}
//
//func (s SimpleFileDataAccess) LoadData() []Product {
//	b, err := ioutil.ReadFile(s.Filename)
//	var ret []Product
//	if err != nil {
//		return []Product{}
//	}
//	json.Unmarshal(b, &ret)
//	return ret
//}
