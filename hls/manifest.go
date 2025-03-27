// Package hls provides a simple HLS manifest parser and manipulator to operations like merging manifests or removing segments.
package hls

import (
	"errors"
	"strconv"
	"strings"
)

// ParseHlsManifest parses a HLS manifest from a string and returns a Manifest object.
func ParseHlsManifest(data string) (manifest Manifest, err error) {
	lines := strings.Split(data, "\n")

	declaration, lines := lines[0], lines[1:]
	if declaration != "#EXTM3U" {
		err = errors.New("data is not a valid HLS manifest, missing #EXTM3U declaration")
		return
	}

	var tempSegmentGroup *SegmentGroup = nil
	var tempSegment *Segment = nil

	for _, line := range lines {
		if strings.HasPrefix(line, "#EXT-X-VERSION:") {
			if version, parseErr := strconv.ParseUint(strings.TrimPrefix(line, "#EXT-X-VERSION:"), 10, 8); parseErr == nil {
				manifest.version = uint8(version)
			} else {
				err = errors.New("failed to parse #EXT-X-VERSION, invalid value")
				return
			}
		} else if strings.HasPrefix(line, "#EXT-X-TARGETDURATION:") {
			if targetDuration, parseErr := strconv.ParseFloat(strings.TrimPrefix(line, "#EXT-X-TARGETDURATION:"), 32); parseErr == nil {
				manifest.targetDuration = float32(targetDuration)
			} else {
				err = errors.New("failed to parse #EXT-X-TARGETDURATION, invalid value")
				return
			}
		} else if strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:") {
			if mediaSequence, parseErr := strconv.ParseUint(strings.TrimPrefix(line, "#EXT-X-MEDIA-SEQUENCE:"), 10, 32); parseErr == nil {
				manifest.mediaSequence = uint32(mediaSequence)
			} else {
				err = errors.New("failed to parse #EXT-X-MEDIA-SEQUENCE, invalid value")
				return
			}
		} else if strings.HasPrefix(line, "#EXT-X-DISCONTINUITY-SEQUENCE:") {
			if discontinuitySequence, parseErr := strconv.ParseUint(strings.TrimPrefix(line, "#EXT-X-DISCONTINUITY-SEQUENCE:"), 10, 32); parseErr == nil {
				manifest.discontinuitySequence = uint32(discontinuitySequence)
			} else {
				err = errors.New("failed to parse #EXT-X-DISCONTINUITY-SEQUENCE, invalid value")
				return
			}
		} else if strings.HasPrefix(line, "#EXTINF:") {
			if tempSegmentGroup == nil {
				tempSegmentGroup = &SegmentGroup{}
			}

			tempSegment = &Segment{}

			if duration, parseErr := strconv.ParseFloat(strings.Split(strings.TrimPrefix(line, "#EXTINF:"), ",")[0], 32); parseErr == nil {
				tempSegment.duration = float32(duration)
			} else {
				err = errors.New("failed to parse #EXTINF, invalid value")
				return
			}
		} else if strings.HasPrefix(line, "#EXT-DISCONTINUITY") {
			if tempSegmentGroup != nil {
				manifest.segmentGroups = append(manifest.segmentGroups, *tempSegmentGroup)
			}

			tempSegmentGroup = nil
		} else if strings.HasPrefix(line, "#EXT-X-ENDLIST") {
			manifest.hasEndList = true
			break
		} else {
			if tempSegment != nil && tempSegmentGroup != nil {
				tempSegment.path = line
				tempSegmentGroup.segments = append(tempSegmentGroup.segments, *tempSegment)
				tempSegment = nil
			}
		}
	}

	if tempSegment != nil {
		return manifest, errors.New("failed to parse HLS manifest, incomplete segment")
	}

	if tempSegmentGroup != nil {
		manifest.segmentGroups = append(manifest.segmentGroups, *tempSegmentGroup)
	}

	return manifest, nil
}

// Manifest represents a HLS manifest.
type Manifest struct {
	version               uint8          // Version of the manifest
	targetDuration        float32        // Target duration of each segment
	mediaSequence         uint32         // Media sequence number
	discontinuitySequence uint32         // Discontinuity sequence number
	hasEndList            bool           // Indicates if the manifest has the #EXT-X-ENDLIST tag
	segmentGroups         []SegmentGroup // List of segment groups
}

// Version returns the version of the manifest.
func (m *Manifest) Version() uint8 {
	return m.version
}

// TargetDuration returns the target duration of each segment.
func (m *Manifest) TargetDuration() float32 {
	return m.targetDuration
}

// MediaSequence returns the media sequence number.
func (m *Manifest) MediaSequence() uint32 {
	return m.mediaSequence
}

// DiscontinuitySequence returns the discontinuity sequence number.
func (m *Manifest) DiscontinuitySequence() uint32 {
	return m.discontinuitySequence
}

// HasEndList returns true if the manifest has the #EXT-X-ENDLIST tag.
func (m *Manifest) HasEndList() bool {
	return m.hasEndList
}

// SegmentGroupCount returns the number of segment groups in the manifest.
func (m *Manifest) SegmentGroupCount() int {
	return len(m.segmentGroups)
}

// SegmentGroups returns the list of segment groups in the manifest.
func (m *Manifest) SegmentGroups() []SegmentGroup {
	return m.segmentGroups
}

// SegmentCount returns the number of segments in the manifest, summing all segments from all segment groups.
func (m *Manifest) SegmentCount() int {
	var count int = 0

	for _, segmentGroup := range m.segmentGroups {
		count += segmentGroup.SegmentCount()
	}

	return count
}

// Duration returns the total duration of the manifest, summing all durations from all segments.
func (m *Manifest) Duration() float32 {
	var duration float32 = 0.0

	for _, segmentGroup := range m.segmentGroups {
		duration += segmentGroup.Duration()
	}

	return duration
}

// ToString returns the manifest as a string.
func (m *Manifest) ToString() string {
	var builder strings.Builder

	builder.WriteString("#EXTM3U\n")
	builder.WriteString("#EXT-X-VERSION:" + strconv.FormatUint(uint64(m.version), 10) + "\n")
	builder.WriteString("#EXT-X-TARGETDURATION:" + strconv.FormatFloat(float64(m.targetDuration), 'f', -1, 32) + "\n")
	builder.WriteString("#EXT-X-MEDIA-SEQUENCE:" + strconv.FormatUint(uint64(m.mediaSequence), 10) + "\n")
	builder.WriteString("#EXT-X-DISCONTINUITY-SEQUENCE:" + strconv.FormatUint(uint64(m.discontinuitySequence), 10) + "\n")

	for i, segmentGroup := range m.segmentGroups {
		for _, segment := range segmentGroup.Segments() {
			builder.WriteString("#EXTINF:" + strconv.FormatFloat(float64(segment.Duration()), 'f', -1, 32) + "\n")
			builder.WriteString(segment.Path() + "\n")
		}

		if i < len(m.segmentGroups)-1 {
			builder.WriteString("#EXT-DISCONTINUITY\n")
		}
	}

	if m.hasEndList {
		builder.WriteString("#EXT-X-ENDLIST\n")
	}

	return builder.String()
}

// SegmentGroup represents a group of segments in a HLS manifest, usually separated by the #EXT-DISCONTINUITY tag.
type SegmentGroup struct {
	// List of segments in the group
	segments []Segment
}

// SegmentCount returns the number of segments in the group.
func (s *SegmentGroup) SegmentCount() int {
	return len(s.segments)
}

// Duration returns the total duration of the group, summing all durations from all segments.
func (s *SegmentGroup) Duration() float32 {
	var duration float32 = 0.0

	for _, segment := range s.segments {
		duration += segment.duration
	}

	return duration
}

// Segments returns the list of segments in the group.
func (s *SegmentGroup) Segments() []Segment {
	return s.segments
}

// Segment represents a segment in a HLS manifest.
type Segment struct {
	path     string  // Path to the segment
	duration float32 // Duration of the segment
}

// Path returns the path to the segment.
func (s *Segment) Path() string {
	return s.path
}

// Duration returns the duration of the segment.
func (s *Segment) Duration() float32 {
	return s.duration
}
