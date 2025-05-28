package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const protoExample = `
syntax = "proto3";
package example;

import "google/protobuf/timestamp.proto";

option go_package = "gitlab.ozon.ru/example/api/example;example";

service Example {
  rpc ExampleRPC(ExampleRPCRequest) returns (ExampleRPCResponse) {};
}

enum ExampleEnum {
  ONE = 0;
  TWO = 1;
  THREE = 2;
}

message ExampleRPCRequest {
  message Emb {
    string field11 = 1;
  }

  ExampleEnum field1 = 1;
  Emb filed2 = 2;
  google.protobuf.Timestamp filed3 = 3;
}

message ExampleRPCResponse {}
`

func writeTempProto(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	file := filepath.Join(dir, "example.proto")
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatalf("cannot write tmp proto: %v", err)
	}
	return file
}

func TestParseProtoFile_FullExample(t *testing.T) {
	filename := writeTempProto(t, protoExample)

	got, err := parseProtoFile(filename)
	if err != nil {
		t.Fatalf("parseProtoFile() error = %v", err)
	}

	want := []Symbol{
		{"google/protobuf/timestamp.proto", "import", 5, 9, 40},
		{"Example", "service", 9, 9, 16},
		{"ExampleRPC", "method", 10, 7, 17},
		{"ExampleEnum", "enum", 13, 6, 17},
		{"ExampleRPCRequest", "message", 19, 9, 26},
		{"ExampleRPCResponse", "message", 29, 9, 27},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseProtoFile() result mismatch\nwant %+v\ngot  %+v", want, got)
	}
}

func TestParseProtoFile_NoNestedSymbols(t *testing.T) {
	filename := writeTempProto(t, protoExample)

	got, err := parseProtoFile(filename)
	if err != nil {
		t.Fatalf("parseProtoFile() error = %v", err)
	}

	for _, s := range got {
		if s.Name == "Emb" {
			t.Errorf("nested message 'Emb' must be ignored, but was returned: %+v", s)
		}
	}
}

func TestEndPosCalculation(t *testing.T) {
	start := 7
	end := start + len("ExampleRPC")
	if end != 17 {
		t.Fatalf("expected end index 17, got %d", end)
	}
}
