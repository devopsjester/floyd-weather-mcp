// Package mcp provides functionality for the Model Context Protocol
package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// Server handles MCP request/response I/O
type Server struct {
	handler     Handler
	encoder     *json.Encoder
	decoder     *json.Decoder
	scanner     *bufio.Scanner
	logFile     *os.File
	interactive bool
}

// NewServer creates a new MCP server
func NewServer(handler Handler) (*Server, error) {
	// Setup logging to a file for debugging
	logFile, err := os.OpenFile("/tmp/floyd-weather-server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %v", err)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.Println("Floyd weather server started")

	return &Server{
		handler:     handler,
		encoder:     json.NewEncoder(os.Stdout),
		decoder:     json.NewDecoder(os.Stdin),
		scanner:     bufio.NewScanner(os.Stdin),
		logFile:     logFile,
		interactive: isInteractiveMode(),
	}, nil
}

// isInteractiveMode checks if the program is running in interactive mode (e.g., VSCode Chat)
func isInteractiveMode() bool {
	// Check if stdin is a terminal
	fileInfo, _ := os.Stdin.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// Serve starts handling requests
func (s *Server) Serve() {
	defer s.logFile.Close()

	log.Printf("Interactive mode: %v", s.interactive)

	if s.interactive {
		s.handleInteractiveMode()
	} else {
		s.handlePipedMode()
	}
}

// handleInteractiveMode handles interactive mode (e.g., VSCode Chat)
func (s *Server) handleInteractiveMode() {
	log.Println("Entering interactive mode handling")
	scanSuccessful := s.scanner.Scan()
	log.Printf("Scanner.Scan() result: %v", scanSuccessful)

	if scanSuccessful {
		input := s.scanner.Text()
		log.Printf("Received input: %s", input)

		var request Request
		if err := json.Unmarshal([]byte(input), &request); err != nil {
			log.Printf("Error unmarshaling request: %v", err)
			s.sendErrorResponse("Error parsing request: " + err.Error())
		} else {
			log.Printf("Unmarshaled request successfully: Method=%s", request.Method)
			response := s.handler.ProcessRequest(request)
			log.Printf("Processed request, encoding response")
			s.encoder.Encode(response)
			log.Println("Request processing complete")
		}
	} else if err := s.scanner.Err(); err != nil {
		log.Printf("Scanner error: %v", err)
		s.sendErrorResponse("Error reading input: " + err.Error())
	} else {
		log.Println("No input received from scanner")
		s.sendErrorResponse("No input received")
	}
}

// handlePipedMode handles non-interactive mode (e.g., piped input)
func (s *Server) handlePipedMode() {
	log.Println("Entering non-interactive mode (piped input)")
	requestCounter := 0

	for {
		var request Request
		log.Printf("Waiting for request #%d", requestCounter+1)

		if err := s.decoder.Decode(&request); err != nil {
			if err == io.EOF {
				log.Println("Reached EOF, exiting")
				break
			}

			log.Printf("Error decoding request: %v", err)
			s.sendErrorResponse("Error decoding request: " + err.Error())
			continue
		}

		requestCounter++
		log.Printf("Processing request #%d: Method=%s", requestCounter, request.Method)

		response := s.handler.ProcessRequest(request)
		s.encoder.Encode(response)

		log.Printf("Completed request #%d", requestCounter)
	}

	log.Printf("Exiting after processing %d requests", requestCounter)
}

// sendErrorResponse sends an error response in MCP format
func (s *Server) sendErrorResponse(message string) {
	s.encoder.Encode(Response{
		Type: "error",
		Content: map[string]string{
			"message": message,
		},
	})
}
