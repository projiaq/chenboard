package storage

import (
	"testing"
	"time"
)

func BenchmarkInsert(b *testing.B) {
	engine, err := NewEngine("./test_data")
	if err != nil {
		b.Fatal(err)
	}
	defer engine.Close()

	err = engine.CreateTable("bench", []string{"ts", "value"}, []DataType{TypeTimestamp, TypeFloat})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Insert("bench", []interface{}{time.Now(), float64(i)})
	}
}

func BenchmarkQuery(b *testing.B) {
	engine, err := NewEngine("./test_data")
	if err != nil {
		b.Fatal(err)
	}
	defer engine.Close()

	err = engine.CreateTable("bench", []string{"ts", "value"}, []DataType{TypeTimestamp, TypeFloat})
	if err != nil {
		b.Fatal(err)
	}

	now := time.Now()
	for i := 0; i < 10000; i++ {
		engine.Insert("bench", []interface{}{now.Add(time.Duration(i) * time.Second), float64(i)})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Query("bench", now, now.Add(10000*time.Second), []string{"ts", "value"})
	}
}

func TestCreateTable(t *testing.T) {
	engine, err := NewEngine("./test_data")
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	err = engine.CreateTable("test", []string{"ts", "value"}, []DataType{TypeTimestamp, TypeInt})
	if err != nil {
		t.Fatal(err)
	}

	err = engine.CreateTable("test", []string{"ts", "value"}, []DataType{TypeTimestamp, TypeInt})
	if err == nil {
		t.Fatal("expected error for duplicate table")
	}
}

func TestInsertAndQuery(t *testing.T) {
	engine, err := NewEngine("./test_data")
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	err = engine.CreateTable("test", []string{"ts", "temp", "device"}, []DataType{TypeTimestamp, TypeFloat, TypeString})
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	err = engine.Insert("test", []interface{}{now, 23.5, "sensor001"})
	if err != nil {
		t.Fatal(err)
	}

	results, err := engine.Query("test", now.Add(-time.Hour), now.Add(time.Hour), []string{})
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0]["device"] != "sensor001" {
		t.Fatalf("expected device sensor001, got %v", results[0]["device"])
	}
}
