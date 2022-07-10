// Package filestore provides storage functions from/to repository form/into file.
package filestore

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/andrei-cloud/go-devops/internal/model"
	"github.com/andrei-cloud/go-devops/internal/repo"
)

// Filestore defines interface for storing and resporing information
// from repository provideing two methods:
//   Store - stores data from repository, returns error if failed
//   Restore - retrived data form file into repository, returns error if  failed.
type Filestore interface {
	Store(repo.Repository) error
	Restore(repo.Repository) error
}

// Base fileSorage structure.
type FileStorage struct {
	// File name for storage.
	filename string
}

var _ Filestore = (*FileStorage)(nil)

// NewFileStorage - creates and instance of FileStorage on file with filename.
func NewFileStorage(filename string) *FileStorage {
	return &FileStorage{
		filename: filename,
	}
}

func (s *FileStorage) openOW() (*os.File, error) {
	return os.OpenFile(s.filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
}

func (s *FileStorage) open() (*os.File, error) {
	return os.OpenFile(s.filename, os.O_RDONLY|os.O_CREATE, 0777)
}

// Store - stores the date from repository to file.
func (s *FileStorage) Store(repo repo.Repository) error {
	file, err := s.openOW()
	if err != nil {
		log.Fatal(err)
	}
	defer s.close(file)

	writer := bufio.NewWriter(file)
	metric := model.Metric{}
	{
		metric.MType = "gauge"
		gauges, err := repo.GetGaugeAll(context.Background())
		if err != nil {
			return err
		}
		for k, v := range gauges {
			metric.ID = k
			metric.Value = &v
			json.NewEncoder(writer).Encode(&metric)
		}
	}

	metric.Delta = nil
	metric.Value = nil

	{
		metric.MType = "counter"
		counters, err := repo.GetCounterAll(context.Background())
		if err != nil {
			return err
		}
		for k, v := range counters {
			metric.ID = k
			metric.Delta = &v
			json.NewEncoder(writer).Encode(&metric)
		}
	}

	return writer.Flush()
}

// Restore - respores the data from file to repository.
func (s *FileStorage) Restore(repo repo.Repository) error {
	file, err := s.open()
	if err != nil {
		log.Fatal(err)
	}
	defer s.close(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Bytes()
		// fmt.Printf("restore: %s\n", string(data))
		metric := model.Metric{}
		err = json.Unmarshal(data, &metric)
		if err != nil {
			return err
		}

		switch metric.MType {
		case "gauge":
			if metric.Value != nil {
				if err := repo.UpdateGauge(context.Background(), metric.ID, *metric.Value); err != nil {
					fmt.Println(err)
				}
			}
		case "counter":
			if metric.Delta != nil {
				if err := repo.UpdateCounter(context.Background(), metric.ID, *metric.Delta); err != nil {
					fmt.Println(err)
				}
			}
		default:
		}
	}

	if err := scanner.Err(); !errors.Is(err, io.EOF) && err != nil {
		return err
	}

	return nil
}

func (s *FileStorage) close(f *os.File) error {
	return f.Close()
}
