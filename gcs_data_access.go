package AppleProductMonitor

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
	"time"
)

type GcsDataAccess struct {
	client *storage.Client
	Bucket string
	Prefix string
}

func NewGcsDataAccess(ctx context.Context, bucket string, prefix string) (*GcsDataAccess, error) {
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

func (g *GcsDataAccess) SaveData(ctx context.Context, data []Product) error {
	fileStream, _ := yaml.Marshal(data)
	w := g.client.Bucket(g.Bucket).Object(strings.Join([]string{g.Prefix, "status.yaml"}, "/")).NewWriter(ctx)
	defer func() {
		_ = w.Close()
	}()
	_, err := w.Write(fileStream)
	return err
}

func (g *GcsDataAccess) LoadData(ctx context.Context) ([]Product, error) {
	var ret []Product
	r, err := g.client.Bucket(g.Bucket).Object(strings.Join([]string{g.Prefix, "status.yaml"}, "/")).NewReader(ctx)
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

	err = yaml.Unmarshal(c, &ret)
	return ret, err
}

func (g *GcsDataAccess) SaveDiff(ctx context.Context, event Event) error {
	dateTag := time.Now().Format(time.Stamp)
	diffPath := strings.Join([]string{g.Prefix, fmt.Sprintf("diff-%s.yaml", dateTag)}, "/")
	w := g.client.Bucket(g.Bucket).Object(diffPath).NewWriter(ctx)
	defer func() {
		_ = w.Close()
	}()
	b, err := yaml.Marshal(event)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
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
