package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/chenboard/tsdb/internal/storage"
	"github.com/gorilla/mux"
)

type Server struct {
	engine *storage.Engine
	router *mux.Router
	server *http.Server
}

type InsertRequest struct {
	Table  string                   `json:"table"`
	Values []map[string]interface{} `json:"values"`
}

type QueryRequest struct {
	Table     string    `json:"table"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Columns   []string  `json:"columns"`
}

type CreateTableRequest struct {
	Table   string   `json:"table"`
	Columns []string `json:"columns"`
	Types   []string `json:"types"`
}

func NewServer(engine *storage.Engine, port string) *Server {
	s := &Server{
		engine: engine,
		router: mux.NewRouter(),
	}

	s.router.HandleFunc("/api/create", s.handleCreate).Methods("POST")
	s.router.HandleFunc("/api/insert", s.handleInsert).Methods("POST")
	s.router.HandleFunc("/api/query", s.handleQuery).Methods("POST")
	s.router.HandleFunc("/health", s.handleHealth).Methods("GET")

	s.server = &http.Server{
		Addr:         ":" + port,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return s
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	types := make([]storage.DataType, len(req.Types))
	for i, t := range req.Types {
		switch t {
		case "timestamp":
			types[i] = storage.TypeTimestamp
		case "int":
			types[i] = storage.TypeInt
		case "float":
			types[i] = storage.TypeFloat
		case "string":
			types[i] = storage.TypeString
		case "bool":
			types[i] = storage.TypeBool
		default:
			http.Error(w, "invalid type: "+t, http.StatusBadRequest)
			return
		}
	}

	if err := s.engine.CreateTable(req.Table, req.Columns, types); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleInsert(w http.ResponseWriter, r *http.Request) {
	var req InsertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get table schema to know column types
	table := s.engine.GetTable(req.Table)
	if table == nil {
		http.Error(w, "table not found", http.StatusNotFound)
		return
	}

	for _, row := range req.Values {
		values := make([]interface{}, len(table.Columns))

		for i, col := range table.Columns {
			val, exists := row[col.Name]
			if !exists {
				http.Error(w, "missing column: "+col.Name, http.StatusBadRequest)
				return
			}

			// Convert timestamp strings to time.Time
			if col.Type == storage.TypeTimestamp {
				if strVal, ok := val.(string); ok {
					t, err := time.Parse(time.RFC3339, strVal)
					if err != nil {
						http.Error(w, "invalid timestamp format: "+err.Error(), http.StatusBadRequest)
						return
					}
					values[i] = t
				} else {
					http.Error(w, "timestamp must be a string", http.StatusBadRequest)
					return
				}
			} else {
				values[i] = val
			}
		}

		if err := s.engine.Insert(req.Table, values); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleQuery(w http.ResponseWriter, r *http.Request) {
	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results, err := s.engine.Query(req.Table, req.StartTime, req.EndTime, req.Columns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"data":   results,
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
