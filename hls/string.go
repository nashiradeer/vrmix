package hls

import (
	"errors"
	"strconv"
	"strings"
)

const (
	// DeclarationField is the field that indicates the start of a HLS manifest.
	DeclarationField = "#EXTM3U"

	// VersionField is the field that indicates the version of the manifest.
	VersionField = "#EXT-X-VERSION"

	// TargetDurationField is the field that indicates the target duration of each segment.
	TargetDurationField = "#EXT-X-TARGETDURATION"

	// MediaSequenceField is the field that indicates the media sequence number.
	MediaSequenceField = "#EXT-X-MEDIA-SEQUENCE"

	// DiscontinuitySequenceField is the field that indicates the discontinuity sequence number.
	DiscontinuitySequenceField = "#EXT-X-DISCONTINUITY-SEQUENCE"

	// SegmentField is the field that indicates a segment in the manifest.
	SegmentField = "#EXTINF"

	// DiscontinuityField is the field that indicates a discontinuity in the manifest.
	DiscontinuityField = "#EXT-DISCONTINUITY"

	// EndListField is the field that indicates the end of the manifest.
	EndListField = "#EXT-X-ENDLIST"
)

var (
	// ErrRequiredFieldMissing indicates that a required field is missing.
	ErrRequiredFieldMissing = errors.New("missing required field")

	// ErrSegmentPathMissing indicates that the path of a segment is missing.
	ErrSegmentPathMissing = errors.New("missing segment path")

	// ErrFieldRequireValue indicates that a field requires a value.
	ErrFieldRequireValue = errors.New("field requires a value")

	// ErrInvalidField indicates that a field is invalid.
	ErrInvalidField = errors.New("invalid field")
)

// ParseError records a parsing error in a HLS manifest.
type ParseError struct {
	Field string // field that caused the error
	Line  int    // line number where the error occurred
	Err   error  // the reason for the error
}

func (e *ParseError) Error() string {
	return "failed to parse " + e.Field + " at line " + strconv.Itoa(e.Line) + ": " + e.Err.Error()
}

func (e *ParseError) Unwrap() error { return e.Err }

func declarationError() *ParseError {
	return &ParseError{Field: DeclarationField, Line: 1, Err: ErrRequiredFieldMissing}
}

func fieldError(field string, line int, err error) *ParseError {
	return &ParseError{Field: field, Line: line, Err: err}
}

func valueError(field string, line int) *ParseError {
	return &ParseError{Field: field, Line: line, Err: ErrFieldRequireValue}
}

func invalidFieldError(field string, line int) *ParseError {
	return &ParseError{Field: field, Line: line, Err: ErrInvalidField}
}

func segmentPathError(line int) *ParseError {
	return &ParseError{Field: SegmentField, Line: line, Err: ErrSegmentPathMissing}
}

// getValue returns the value of a field or an empty string if the field or value is missing.
func getValue(line string) string {
	_, value, found := strings.Cut(line, ":")
	if !found {
		return ""
	}

	return value
}

// parseUintValue parses an unsigned integer value from the value of a field in a line, returning the value or an error already wrapped on a ParseError.
func parseUintValue(field string, line string, lineNumber int, bitSize int) (uint64, error) {
	value := getValue(line)
	if value == "" {
		return 0, valueError(field, lineNumber)
	}

	result, err := strconv.ParseUint(value, 10, bitSize)
	if err != nil {
		return 0, fieldError(field, lineNumber, err)
	}

	return result, nil
}

// ParseHlsManifest parses a HLS manifest from a string and returns a Manifest object.
func ParseHlsManifest(data string) (Manifest, error) {
	lines := strings.Split(data, "\n")
	manifest := Manifest{}

	declaration, lines := lines[0], lines[1:]
	if declaration != DeclarationField {
		return manifest, declarationError()
	}

	var tempSegmentGroup *SegmentGroup = nil
	var tempSegment *Segment = nil

	for i, line := range lines {
		lineNumber := i + 2

		if strings.HasPrefix(line, VersionField) {
			version, err := parseUintValue(VersionField, line, lineNumber, 8)
			if err != nil {
				return manifest, err
			}

			manifest.Version = uint8(version)
		} else if strings.HasPrefix(line, TargetDurationField) {
			duration, err := parseUintValue(TargetDurationField, line, lineNumber, 8)
			if err != nil {
				return manifest, err
			}

			manifest.TargetDuration = uint8(duration)
		} else if strings.HasPrefix(line, MediaSequenceField) {
			mediaSequence, err := parseUintValue(MediaSequenceField, line, lineNumber, 32)
			if err != nil {
				return manifest, err
			}

			manifest.MediaSequence = uint32(mediaSequence)
		} else if strings.HasPrefix(line, DiscontinuitySequenceField) {
			discontinuitySequence, err := parseUintValue(DiscontinuitySequenceField, line, lineNumber, 32)
			if err != nil {
				return manifest, err
			}

			manifest.DiscontinuitySequence = uint32(discontinuitySequence)
		} else if strings.HasPrefix(line, SegmentField) {
			if tempSegmentGroup == nil {
				tempSegmentGroup = &SegmentGroup{}
			}

			if tempSegment != nil {
				return manifest, segmentPathError(lineNumber)
			}
			tempSegment = &Segment{}

			value := getValue(line)
			durationValue, title, found := strings.Cut(value, ",")
			if !found {
				return manifest, valueError(SegmentField, lineNumber)
			}

			duration, err := strconv.ParseFloat(durationValue, 32)
			if err != nil {
				return manifest, fieldError(SegmentField, lineNumber, err)
			}

			tempSegment.Duration = float32(duration)
			tempSegment.Title = title
		} else if strings.HasPrefix(line, DiscontinuityField) {
			if tempSegmentGroup != nil {
				manifest.SegmentGroups = append(manifest.SegmentGroups, *tempSegmentGroup)
			}

			tempSegmentGroup = nil
		} else if strings.HasPrefix(line, EndListField) {
			manifest.HasEndList = true
			break
		} else {
			if tempSegment != nil && tempSegmentGroup != nil {
				tempSegment.Path = line
				tempSegmentGroup.Segments = append(tempSegmentGroup.Segments, *tempSegment)
				tempSegment = nil
			} else {
				return manifest, invalidFieldError(line, lineNumber)
			}
		}
	}

	if tempSegment != nil {
		return manifest, segmentPathError(len(lines) + 1)
	}

	if tempSegmentGroup != nil {
		manifest.SegmentGroups = append(manifest.SegmentGroups, *tempSegmentGroup)
	}

	return manifest, nil
}

// ToString returns the manifest as a string.
func (m *Manifest) String() string {
	var builder strings.Builder

	builder.WriteString(DeclarationField + "\n")
	builder.WriteString(VersionField + ":" + strconv.FormatUint(uint64(m.Version), 10) + "\n")
	builder.WriteString(TargetDurationField + ":" + strconv.FormatFloat(float64(m.TargetDuration), 'f', -1, 32) + "\n")
	builder.WriteString(MediaSequenceField + ":" + strconv.FormatUint(uint64(m.MediaSequence), 10) + "\n")
	builder.WriteString(DiscontinuitySequenceField + ":" + strconv.FormatUint(uint64(m.DiscontinuitySequence), 10) + "\n")

	for i, segmentGroup := range m.SegmentGroups {
		for _, segment := range segmentGroup.Segments {
			builder.WriteString(SegmentField + ":" + strconv.FormatFloat(float64(segment.Duration), 'f', -1, 32) + "," + segment.Title + "\n")
			builder.WriteString(segment.Path + "\n")
		}

		if i < len(m.SegmentGroups)-1 {
			builder.WriteString(DiscontinuityField + "\n")
		}
	}

	if m.HasEndList {
		builder.WriteString(EndListField + "\n")
	}

	return builder.String()
}
