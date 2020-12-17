package AppleProductMonitor

import (
	"encoding/json"
	"io/ioutil"
)

type DataService interface {
	SaveData(data []Product)
	LoadData() []Product
}

type SimpleDataAccess struct {
	savedData []Product
}

func NewSimpleDataAccess() *SimpleDataAccess {
	return &SimpleDataAccess{savedData: make([]Product, 0)}
}

func (d *SimpleDataAccess) SaveData(data []Product) {
	d.savedData = data
}

func (d *SimpleDataAccess) LoadData() []Product {
	return d.savedData
}

type SimpleFileDataAccess struct {
	Filename string
}

func (s SimpleFileDataAccess) SaveData(data []Product) {
	b, _ := json.Marshal(data)
	ioutil.WriteFile(s.Filename, b, 0644)
}

func (s SimpleFileDataAccess) LoadData() []Product {
	b, err := ioutil.ReadFile(s.Filename)
	var ret []Product
	if err != nil {
		return []Product{}
	}
	json.Unmarshal(b, &ret)
	return ret
}

func NewSimpleFileDataAccess() *SimpleFileDataAccess {
	return &SimpleFileDataAccess{Filename: "./vault/status.json"}
}
