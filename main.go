package main

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protowire"
)

const (
	META = "./dji_meta.bin"
	DBGI = "./dji_dbgi.bin"
)

var tagNames = map[string]map[int32]string{
	META: {
		1: "description",
		2: "video info",
		3: "sample data",
	},
	DBGI: {
		1: "description",
		2: "sample data",
	},
}

func main() {
	ConfigureLogging(true, zerolog.InfoLevel, "Dev", "DJI Osmo3 Metadata Experiment")

	Read(META)
	Read(DBGI)
}

func Read(filename string) {
	binaryData := DumpFile(filename)

	fields := ParseUnknown(binaryData)
	log.Info().Int("fields", len(fields)).Msg("Fields read and parsed")

	tagNumbers := map[int32]int{}

	for _, field := range fields {
		/*
			if (filename == META && i < 3) ||  (filename == DBGI && i < 2) {
				log.Info().
					Int("fieldNumber", i).
					Int32("tagNum", field.Tag.Num).
					Str("tagType", ProtoTypeString(field.Tag.Type)).
					Interface("valPayload", field.Val.Payload).
					Int("valLength", field.Val.Length).
					Send()
			}
			//*/
		if _, ok := tagNumbers[field.Tag.Num]; ok {
			tagNumbers[field.Tag.Num]++
		} else {
			tagNumbers[field.Tag.Num] = 1
		}
	}

	for tag, count := range tagNumbers {
		log.Info().Str("filename", filename).Str("tag", tagNames[filename][tag]).Int("count", count).Msg("Top level tag")
	}
}

func ProtoTypeString(t protowire.Type) string {
	switch t {
	case protowire.VarintType:
		return "int"
	case protowire.Fixed32Type:
		return "int32"
	case protowire.Fixed64Type:
		return "int64"
	case protowire.BytesType:
		return "[]byte"
	case protowire.StartGroupType:
		return "startGroup"
	case protowire.EndGroupType:
		return "endGroup"
	default:
		return "unknown"
	}
}

func DumpFile(filename string) []byte {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal().Err(err).Str("filename", filename).Msg("Failed to open file")
	}

	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal().Err(err).Str("filename", filename).Msg("Failed to read entire file")
	}
	log.Info().Str("filename", filename).Int("bytesRead", len(b)).Msg("Read file")
	return b
}

func ConfigureLogging(interactive bool, logLevel zerolog.Level, versionStr, versionMsgStr string) {
	if interactive {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			},
		)
	}
	zerolog.SetGlobalLevel(logLevel)
}

type Tag struct {
	Num     int32
	Type    protowire.Type
	TypeStr string
}

type Val struct {
	Payload interface{}
	// PayloadAsString string `json:",omitempty"`
	Length int
}

type Field struct {
	Tag Tag
	Val Val
}

// From https://stackoverflow.com/questions/41348512/protobuf-unmarshal-unknown-message/69510141#69510141
func ParseUnknown(b []byte) []Field {
	fields := make([]Field, 0)
	for len(b) > 0 {
		n, t, fieldlen := protowire.ConsumeField(b)
		if fieldlen < 1 {
			return nil
		}
		field := Field{
			Tag: Tag{Num: int32(n), Type: t, TypeStr: ProtoTypeString(t)},
		}

		_, _, taglen := protowire.ConsumeTag(b[:fieldlen])
		if taglen < 1 {
			return nil
		}

		var (
			v    interface{}
			vlen int
		)
		switch t {
		case protowire.VarintType:
			v, vlen = protowire.ConsumeVarint(b[taglen:fieldlen])

		case protowire.Fixed64Type:
			v, vlen = protowire.ConsumeFixed64(b[taglen:fieldlen])

		case protowire.BytesType:
			v, vlen = protowire.ConsumeBytes(b[taglen:fieldlen])
			sub := ParseUnknown(v.([]byte))
			if sub != nil {
				v = sub
			}

		case protowire.StartGroupType:
			v, vlen = protowire.ConsumeGroup(n, b[taglen:fieldlen])
			sub := ParseUnknown(v.([]byte))
			if sub != nil {
				v = sub
			}

		case protowire.Fixed32Type:
			v, vlen = protowire.ConsumeFixed32(b[taglen:fieldlen])
		}

		if vlen < 1 {
			return nil
		}

		field.Val = Val{Payload: v, Length: vlen - taglen}
		/*
		        if t == protowire.BytesType{
					str, ok := v.([]byte)
					if ok {
						field.Val.PayloadAsString = string(str)
					}
				}
		*/

		fields = append(fields, field)
		b = b[fieldlen:]
	}
	return fields
}
