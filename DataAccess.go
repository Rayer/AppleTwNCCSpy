package AppleProductMonitor

import (
	"context"
	"encoding/json"
	"io/ioutil"
)

type DataAccess interface {
	SaveData(ctx context.Context, data []Product) (err error)
	LoadData(ctx context.Context) ([]Product, error)
}

type SimpleDataAccess struct {
	savedData []Product
}

func NewSimpleDataAccess() *SimpleDataAccess {
	return &SimpleDataAccess{savedData: make([]Product, 0)}
}

func (d *SimpleDataAccess) SaveData(ctx context.Context, data []Product) error {
	d.savedData = data
	return nil
}

func (d *SimpleDataAccess) LoadData(context.Context) ([]Product, error) {
	return d.savedData, nil
}

type SimpleFileDataAccess struct {
	Filename string
}

func (s SimpleFileDataAccess) SaveData(ctx context.Context, data []Product) error {
	b, _ := json.Marshal(data)
	ioutil.WriteFile(s.Filename, b, 0644)
	return nil
}

func (s SimpleFileDataAccess) LoadData(context.Context) ([]Product, error) {
	b, err := ioutil.ReadFile(s.Filename)
	var ret []Product
	if err != nil {
		return []Product{}, nil
	}
	json.Unmarshal(b, &ret)
	return ret, nil
}

func NewSimpleFileDataAccess() *SimpleFileDataAccess {
	return &SimpleFileDataAccess{Filename: "./vault/status.json"}
}
