package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestBooleanQuestion(t *testing.T) {
	type args struct {
		reader       *bufio.Reader
		label        string
		defaultValue bool
	}
	tests := []struct {
		name       string
		args       args
		want       bool
		wantWriter string
		wantErr    bool
	}{
		{
			name: "answer yes",
			args: args{
				reader:       bufio.NewReader(bytes.NewBufferString("yes\n")),
				label:        "Do you like this?",
				defaultValue: true,
			},
			want:       true,
			wantWriter: "Do you like this? [Y/N]: ",
		},
		{
			name: "answer y",
			args: args{
				reader:       bufio.NewReader(bytes.NewBufferString("y\n")),
				label:        "Do you like this?",
				defaultValue: true,
			},
			want:       true,
			wantWriter: "Do you like this? [Y/N]: ",
		},
		{
			name: "answer no",
			args: args{
				reader:       bufio.NewReader(bytes.NewBufferString("no\n")),
				label:        "Do you like this?",
				defaultValue: true,
			},
			want:       false,
			wantWriter: "Do you like this? [Y/N]: ",
		},
		{
			name: "answer n",
			args: args{
				reader:       bufio.NewReader(bytes.NewBufferString("n\n")),
				label:        "Do you like this?",
				defaultValue: true,
			},
			want:       false,
			wantWriter: "Do you like this? [Y/N]: ",
		},
		{
			name: "invalid input",
			args: args{
				reader:       bufio.NewReader(bytes.NewBufferString("invalid\n")),
				label:        "Do you like this?",
				defaultValue: true,
			},
			want:       false,
			wantWriter: "Do you like this? [Y/N]: ",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			got, err := BooleanQuestion(writer, tt.args.reader, tt.args.label, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("BooleanQuestion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BooleanQuestion() = %v, want %v", got, tt.want)
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("BooleanQuestion() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestStringQuestion(t *testing.T) {
	type args struct {
		reader       *bufio.Reader
		label        string
		defaultValue string
		validations  []ValidationFunc[string]
	}
	tests := []struct {
		name       string
		args       args
		want       string
		wantWriter string
		wantErr    bool
	}{
		{
			name: "simple string question",
			args: args{
				reader:       bufio.NewReader(bytes.NewBufferString("answer\n")),
				label:        "What is your answer?",
				defaultValue: "default",
			},
			want:       "answer",
			wantWriter: "What is your answer? [default]: ",
		},
		{
			name: "a longer text",
			args: args{
				reader:       bufio.NewReader(bytes.NewBufferString("a longer text\n")),
				label:        "What is your answer?",
				defaultValue: "default",
			},
			want:       "a longer text",
			wantWriter: "What is your answer? [default]: ",
		},
		{
			name: "empty input",
			args: args{
				reader:       bufio.NewReader(bytes.NewBufferString("\n")),
				label:        "What is your answer?",
				defaultValue: "default",
			},
			want:       "default",
			wantWriter: "What is your answer? [default]: ",
		},
		{
			name: "empty input with default",
			args: args{
				reader:       bufio.NewReader(bytes.NewBufferString("\n")),
				label:        "What is your answer?",
				defaultValue: "default",
			},
			want:       "default",
			wantWriter: "What is your answer? [default]: ",
		},
		{
			name: "first invalid validation",
			args: args{
				reader:       bufio.NewReader(bytes.NewBufferString("invalid\n")),
				label:        "What is your answer?",
				defaultValue: "default",
				validations: []ValidationFunc[string]{
					func(s string) error {
						return fmt.Errorf("invalid input")
					},
				},
			},
			want:       "default",
			wantWriter: "What is your answer? [default]: ",
			wantErr:    true,
		},
		{
			name: "second invalid validation",
			args: args{
				reader:       bufio.NewReader(bytes.NewBufferString("invalid\n")),
				label:        "What is your answer?",
				defaultValue: "default",
				validations: []ValidationFunc[string]{
					func(s string) error {
						return nil
					},
					func(s string) error {
						return fmt.Errorf("invalid input")
					},
				},
			},
			want:       "default",
			wantErr:    true,
			wantWriter: "What is your answer? [default]: ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			got, err := StringQuestion(writer, tt.args.reader, tt.args.label, tt.args.defaultValue, tt.args.validations...)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringQuestion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StringQuestion() = %v, want %v", got, tt.want)
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("StringQuestion() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
