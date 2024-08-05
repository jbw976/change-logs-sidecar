package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/protobuf/encoding/protojson"

	changelogs "github.com/crossplane/crossplane-runtime/apis/changelogs/proto/v1alpha1"
)

type Server struct {
	changelogs.UnimplementedChangeLogServiceServer
}

func (s *Server) SendChangeLog(ctx context.Context, entry *changelogs.SendChangeLogRequest) (*changelogs.SendChangeLogResponse, error) {
	// Marshal the change log entry coming over the wire to JSON using the
	// protojson helper
	b, err := protojson.Marshal(entry)
	if err != nil {
		return &changelogs.SendChangeLogResponse{
			Success: false,
			Message: "Failed to marshall input entry",
		}, err
	}

	// Now unmarshal those bytes into a map so we can inject fields/data that
	// are server side only
	var entryData map[string]interface{}
	if err := json.Unmarshal(b, &entryData); err != nil {
		return &changelogs.SendChangeLogResponse{
			Success: false,
			Message: "Failed to unmarshal input entry",
		}, err
	}

	// Include a server side current timestamp into the data
	entryData["timestamp"] = time.Now().UTC().Format(time.RFC3339)

	// Marshal the data back to JSON bytes once again so we can write the final
	// product to the output
	b, err = json.Marshal(entryData)
	if err != nil {
		return &changelogs.SendChangeLogResponse{
			Success: false,
			Message: "Failed to marshal final entry",
		}, err
	}

	// write the final change log entry to stdout
	fmt.Printf("%s\n", string(b))

	return &changelogs.SendChangeLogResponse{
		Success: true,
		Message: "Change log entry received",
	}, nil
}
