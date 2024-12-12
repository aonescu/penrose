package utils

import (
	"bufio"
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/tuneinsight/lattigo/v5/utils/buffer"
)

// Parameters holds the configuration for Lattigo context
type Parameters struct {
	Scheme        SchemeParameters
	Bootstrapping *BootstrappingParameters
}

type SchemeParameters struct {
	LogN            int
	LogQ            []int
	LogP            []int
	LogDefaultScale int
}

type BootstrappingParameters struct {
	Enable    bool
	Threshold int
}

type EvaluationKeySet struct {
	// Add fields relevant for evaluation key (if necessary)
}

func Serialize(object interface{}, path string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("os.Create(%s): %w", path, err)
	}
	defer f.Close()

	switch object := object.(type) {
	case io.WriterTo:
		if _, err = object.WriteTo(f); err != nil {
			return fmt.Errorf("%T.WriteTo: %w", object, err)
		}
	case encoding.BinaryMarshaler:
		var data []byte
		if data, err = object.MarshalBinary(); err != nil {
			return fmt.Errorf("%T.MarshalBinary: %w", object, err)
		}
		if _, err = f.Write(data); err != nil {
			return fmt.Errorf("file.Write: %w", err)
		}
	default:
		return fmt.Errorf("%T does not implement io.WriterTo or encoding.BinaryMarshaler", object)
	}

	return nil
}

func Deserialize(object interface{}, path string) (err error) {
	switch object := object.(type) {
	case io.ReaderFrom:
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("os.Open(%s): %w", path, err)
		}
		defer f.Close()

		if _, err = object.ReadFrom(f); err != nil {
			return fmt.Errorf("%T.ReadFrom: %w", object, err)
		}
	case encoding.BinaryUnmarshaler:
		var data []byte
		if data, err = os.ReadFile(path); err != nil {
			return fmt.Errorf("os.ReadFile(%s): %w", path, err)
		}

		if err = object.UnmarshalBinary(data); err != nil {
			return fmt.Errorf("%T.UnmarshalBinary: %w", object, err)
		}

	default:
		return fmt.Errorf("%T does not implement io.ReaderFrom or encoding.BinaryUnmarshaler", object)
	}

	return nil
}

func (p Parameters) BinarySize() int {
	data, _ := p.MarshalJSON()
	return len(data) + 4
}

func (p Parameters) WriteTo(w io.Writer) (n int64, err error) {
	switch w := w.(type) {
	case buffer.Writer:
		var data []byte
		if data, err = p.MarshalJSON(); err != nil {
			return
		}

		if n, err = buffer.WriteAsUint32[int](w, len(data)); err != nil {
			return n, fmt.Errorf("buffer.WriteAsUint32[int]: %w", err)
		}

		var inc int
		if inc, err = w.Write(data); err != nil {
			return int64(n), fmt.Errorf("io.Write.Write: %w", err)
		}

		n += int64(inc)

		return n, w.Flush()
	default:
		return p.WriteTo(bufio.NewWriter(w))
	}
}

func (p *Parameters) ReadFrom(r io.Reader) (n int64, err error) {
	switch r := r.(type) {
	case buffer.Reader:
		var size int
		if n, err = buffer.ReadAsUint32[int](r, &size); err != nil {
			return int64(n), fmt.Errorf("buffer.ReadAsUint64[int]: %w", err)
		}

		bytes := make([]byte, size)

		var inc int
		if inc, err = r.Read(bytes); err != nil {
			return n + int64(inc), fmt.Errorf("io.Reader.Read: %w", err)
		}

		return n + int64(inc), p.UnmarshalJSON(bytes)

	default:
		return p.ReadFrom(bufio.NewReader(r))
	}
}

func (p Parameters) MarshalJSON() (data []byte, err error) {
	aux := struct {
		Scheme        SchemeParameters
		Bootstrapping BootstrappingParameters
	}{
		Scheme:        p.Scheme,
		Bootstrapping: *p.Bootstrapping,
	}

	return json.Marshal(aux)
}

func (p *Parameters) UnmarshalJSON(data []byte) (err error) {
	aux := struct {
		Scheme        SchemeParameters
		Bootstrapping BootstrappingParameters
	}{}

	if err = json.Unmarshal(data, &aux); err != nil {
		return
	}

	p.Scheme = aux.Scheme
	p.Bootstrapping = &aux.Bootstrapping
	return
}

// Utility functions for managing encryption keys or other operations
func DeserializeEvaluationKey(evk *EvaluationKeySet, path string) error {
	return Deserialize(evk, path)
}
