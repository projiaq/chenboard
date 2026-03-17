package storage

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/klauspost/compress/zstd"
)

type DataType uint8

const (
	TypeTimestamp DataType = iota
	TypeInt
	TypeFloat
	TypeString
	TypeBool
)

type Column struct {
	Name string
	Type DataType
	Data []interface{}
}

type Table struct {
	Name    string
	Columns []*Column
	mu      sync.RWMutex
}

type Engine struct {
	dataDir string
	tables  map[string]*Table
	mu      sync.RWMutex
	encoder *zstd.Encoder
	decoder *zstd.Decoder
}

func NewEngine(dataDir string) (*Engine, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	encoder, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedFastest))
	if err != nil {
		return nil, err
	}

	decoder, err := zstd.NewReader(nil)
	if err != nil {
		return nil, err
	}

	e := &Engine{
		dataDir: dataDir,
		tables:  make(map[string]*Table),
		encoder: encoder,
		decoder: decoder,
	}

	return e, nil
}

func (e *Engine) CreateTable(name string, columns []string, types []DataType) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.tables[name]; exists {
		return fmt.Errorf("table %s already exists", name)
	}

	cols := make([]*Column, len(columns))
	for i, colName := range columns {
		cols[i] = &Column{
			Name: colName,
			Type: types[i],
			Data: make([]interface{}, 0, 1024),
		}
	}

	e.tables[name] = &Table{
		Name:    name,
		Columns: cols,
	}

	return nil
}

func (e *Engine) Insert(tableName string, values []interface{}) error {
	e.mu.RLock()
	table, exists := e.tables[tableName]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("table %s not found", tableName)
	}

	table.mu.Lock()
	defer table.mu.Unlock()

	if len(values) != len(table.Columns) {
		return fmt.Errorf("column count mismatch: expected %d, got %d", len(table.Columns), len(values))
	}

	for i, col := range table.Columns {
		col.Data = append(col.Data, values[i])
	}

	return nil
}

func (e *Engine) Query(tableName string, startTime, endTime time.Time, columns []string) ([]map[string]interface{}, error) {
	e.mu.RLock()
	table, exists := e.tables[tableName]
	e.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("table %s not found", tableName)
	}

	table.mu.RLock()
	defer table.mu.RUnlock()

	if len(table.Columns) == 0 || len(table.Columns[0].Data) == 0 {
		return []map[string]interface{}{}, nil
	}

	results := make([]map[string]interface{}, 0)
	rowCount := len(table.Columns[0].Data)

	for i := 0; i < rowCount; i++ {
		ts, ok := table.Columns[0].Data[i].(time.Time)
		if !ok {
			continue
		}

		if ts.Before(startTime) || ts.After(endTime) {
			continue
		}

		row := make(map[string]interface{})
		for _, col := range table.Columns {
			if len(columns) == 0 || contains(columns, col.Name) {
				row[col.Name] = col.Data[i]
			}
		}
		results = append(results, row)
	}

	return results, nil
}

func (e *Engine) Flush(tableName string) error {
	e.mu.RLock()
	table, exists := e.tables[tableName]
	e.mu.RUnlock()

	if !exists {
		return fmt.Errorf("table %s not found", tableName)
	}

	table.mu.RLock()
	defer table.mu.RUnlock()

	filePath := filepath.Join(e.dataDir, tableName+".tsdb")
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	binary.Write(f, binary.LittleEndian, uint32(len(table.Columns)))

	for _, col := range table.Columns {
		nameBytes := []byte(col.Name)
		binary.Write(f, binary.LittleEndian, uint32(len(nameBytes)))
		f.Write(nameBytes)
		binary.Write(f, binary.LittleEndian, col.Type)
		binary.Write(f, binary.LittleEndian, uint32(len(col.Data)))
	}

	return nil
}

func (e *Engine) Close() error {
	e.encoder.Close()
	e.decoder.Close()
	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
