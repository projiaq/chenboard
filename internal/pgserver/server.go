package pgserver

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/chenboard/tsdb/internal/storage"
	"github.com/jackc/pgproto3/v2"
	"github.com/xwb1989/sqlparser"
)

type PgServer struct {
	engine   *storage.Engine
	listener net.Listener
	port     string
}

func NewPgServer(engine *storage.Engine, port string) *PgServer {
	return &PgServer{
		engine: engine,
		port:   port,
	}
}

func (s *PgServer) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", ":"+s.port)
	if err != nil {
		return err
	}

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		go s.handleConnection(conn)
	}
}

func (s *PgServer) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *PgServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	backend := pgproto3.NewBackend(pgproto3.NewChunkReader(conn), conn)

	// Startup
	startupMsg, err := backend.ReceiveStartupMessage()
	if err != nil {
		return
	}

	switch startupMsg.(type) {
	case *pgproto3.StartupMessage:
		// Send authentication OK
		buf, _ := (&pgproto3.AuthenticationOk{}).Encode(nil)
		buf, _ = (&pgproto3.ReadyForQuery{TxStatus: 'I'}).Encode(buf)
		_, err = conn.Write(buf)
		if err != nil {
			return
		}
	case *pgproto3.SSLRequest:
		// Reject SSL for simplicity
		_, err = conn.Write([]byte("N"))
		if err != nil {
			return
		}
		s.handleConnection(conn)
		return
	default:
		return
	}

	// Main query loop
	for {
		msg, err := backend.Receive()
		if err != nil {
			if err != io.EOF {
				s.sendError(backend, err.Error())
			}
			return
		}

		switch m := msg.(type) {
		case *pgproto3.Query:
			s.handleQuery(backend, m.String)
		case *pgproto3.Terminate:
			return
		}
	}
}

func (s *PgServer) handleQuery(backend *pgproto3.Backend, query string) {
	query = strings.TrimSpace(query)

	// Parse SQL
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		s.sendError(backend, fmt.Sprintf("SQL parse error: %v", err))
		return
	}

	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		s.handleSelect(backend, stmt)
	case *sqlparser.Insert:
		s.handleInsert(backend, stmt)
	case *sqlparser.DDL:
		s.handleDDL(backend, stmt)
	default:
		s.sendError(backend, "Unsupported query type")
	}
}

func (s *PgServer) handleSelect(backend *pgproto3.Backend, stmt *sqlparser.Select) {
	// Extract table name
	tableName := sqlparser.String(stmt.From[0])

	// Simple query: SELECT * FROM table
	results, err := s.engine.Query(tableName, time.Time{}, time.Now().Add(100*365*24*time.Hour), []string{})
	if err != nil {
		s.sendError(backend, err.Error())
		return
	}

	if len(results) == 0 {
		s.sendEmptyResult(backend)
		return
	}

	// Build column descriptions
	var columns []string
	for k := range results[0] {
		columns = append(columns, k)
	}

	// Send row description
	fields := make([]pgproto3.FieldDescription, len(columns))
	for i, col := range columns {
		fields[i] = pgproto3.FieldDescription{
			Name:                 []byte(col),
			TableOID:             0,
			TableAttributeNumber: 0,
			DataTypeOID:          25, // text
			DataTypeSize:         -1,
			TypeModifier:         -1,
			Format:               0,
		}
	}

	buf, _ := (&pgproto3.RowDescription{Fields: fields}).Encode(nil)

	// Send data rows
	for _, row := range results {
		values := make([][]byte, len(columns))
		for i, col := range columns {
			values[i] = []byte(fmt.Sprintf("%v", row[col]))
		}
		buf, _ = (&pgproto3.DataRow{Values: values}).Encode(buf)
	}

	// Send command complete
	buf, _ = (&pgproto3.CommandComplete{CommandTag: []byte(fmt.Sprintf("SELECT %d", len(results)))}).Encode(buf)
	buf, _ = (&pgproto3.ReadyForQuery{TxStatus: 'I'}).Encode(buf)

	backend.Send(buf)
}

func (s *PgServer) handleInsert(backend *pgproto3.Backend, stmt *sqlparser.Insert) {
	tableName := sqlparser.String(stmt.Table)

	// Extract values
	rows, ok := stmt.Rows.(sqlparser.Values)
	if !ok {
		s.sendError(backend, "Invalid INSERT syntax")
		return
	}

	insertCount := 0
	for _, row := range rows {
		values := make([]interface{}, len(row))
		for i, expr := range row {
			val := sqlparser.String(expr)
			// Remove quotes if string literal
			val = strings.Trim(val, "'\"")
			values[i] = val
		}

		if err := s.engine.Insert(tableName, values); err != nil {
			s.sendError(backend, err.Error())
			return
		}
		insertCount++
	}

	buf, _ := (&pgproto3.CommandComplete{CommandTag: []byte(fmt.Sprintf("INSERT 0 %d", insertCount))}).Encode(nil)
	buf, _ = (&pgproto3.ReadyForQuery{TxStatus: 'I'}).Encode(buf)
	backend.Send(buf)
}

func (s *PgServer) handleDDL(backend *pgproto3.Backend, stmt *sqlparser.DDL) {
	if stmt.Action != "create" {
		s.sendError(backend, "Only CREATE TABLE is supported")
		return
	}

	tableName := stmt.NewName.Name.String()
	var columns []string
	var types []storage.DataType

	for _, col := range stmt.TableSpec.Columns {
		columns = append(columns, col.Name.String())

		// Map SQL types to internal types
		sqlType := strings.ToLower(col.Type.Type)
		var dataType storage.DataType
		switch {
		case strings.Contains(sqlType, "timestamp"), strings.Contains(sqlType, "time"):
			dataType = storage.TypeTimestamp
		case strings.Contains(sqlType, "int"), strings.Contains(sqlType, "bigint"):
			dataType = storage.TypeInt
		case strings.Contains(sqlType, "float"), strings.Contains(sqlType, "double"), strings.Contains(sqlType, "real"):
			dataType = storage.TypeFloat
		case strings.Contains(sqlType, "bool"):
			dataType = storage.TypeBool
		default:
			dataType = storage.TypeString
		}
		types = append(types, dataType)
	}

	if err := s.engine.CreateTable(tableName, columns, types); err != nil {
		s.sendError(backend, err.Error())
		return
	}

	buf, _ := (&pgproto3.CommandComplete{CommandTag: []byte("CREATE TABLE")}).Encode(nil)
	buf, _ = (&pgproto3.ReadyForQuery{TxStatus: 'I'}).Encode(buf)
	backend.Send(buf)
}

func (s *PgServer) sendError(backend *pgproto3.Backend, message string) {
	buf, _ := (&pgproto3.ErrorResponse{
		Severity: "ERROR",
		Code:     "XX000",
		Message:  message,
	}).Encode(nil)
	buf, _ = (&pgproto3.ReadyForQuery{TxStatus: 'I'}).Encode(buf)
	backend.Send(buf)
}

func (s *PgServer) sendEmptyResult(backend *pgproto3.Backend) {
	buf, _ := (&pgproto3.CommandComplete{CommandTag: []byte("SELECT 0")}).Encode(nil)
	buf, _ = (&pgproto3.ReadyForQuery{TxStatus: 'I'}).Encode(buf)
	backend.Send(buf)
}
