package hls

import "math"

// Manifest represents a HLS manifest.
type Manifest struct {
	Version               uint8          // Version of the manifest
	TargetDuration        uint8          // Target duration of each segment
	MediaSequence         uint32         // Media sequence number
	DiscontinuitySequence uint32         // Discontinuity sequence number
	HasEndList            bool           // Indicates if the manifest has the #EXT-X-ENDLIST tag
	SegmentGroups         []SegmentGroup // List of segment groups
}

// SegmentCount returns the number of segments in the manifest, summing all segments from all segment groups.
func (m *Manifest) SegmentCount() int {
	var count = 0

	for _, segmentGroup := range m.SegmentGroups {
		count += len(segmentGroup.Segments)
	}

	return count
}

// Duration returns the total duration of the manifest, summing all durations from all segments.
func (m *Manifest) Duration() float64 {
	var duration = 0.0

	for _, segmentGroup := range m.SegmentGroups {
		duration += segmentGroup.Duration()
	}

	return duration
}

// HasValidTargetDuration returns true if the target duration of the manifest is valid, which is when the target duration is greater than or equal to the recommended target duration.
func (m *Manifest) HasValidTargetDuration() bool {
	return !m.ExceedsTargetDuration(m.TargetDuration)
}

// ExceedsTargetDuration returns true if any segment in the manifest has a target duration greater than the specified target duration.
func (m *Manifest) ExceedsTargetDuration(targetDuration uint8) bool {
	for _, segmentGroup := range m.SegmentGroups {
		if segmentGroup.ExceedsTargetDuration(targetDuration) {
			return true
		}
	}

	return false
}

// MaxTargetDuration returns the target duration of the longest segment in the manifest, which is the duration rounded to the nearest integer.
func (m *Manifest) MaxTargetDuration() uint8 {
	var maxTargetDuration uint8 = 0

	for _, segmentGroup := range m.SegmentGroups {
		targetDuration := segmentGroup.MaxTargetDuration()

		if targetDuration > maxTargetDuration {
			maxTargetDuration = targetDuration
		}
	}

	return maxTargetDuration
}

// MaxDuration returns the duration of the longest segment in the manifest.
func (m *Manifest) MaxDuration() float32 {
	var maxDuration float32 = 0.0

	for _, segmentGroup := range m.SegmentGroups {
		duration := segmentGroup.MaxDuration()

		if duration > maxDuration {
			maxDuration = duration
		}
	}

	return maxDuration
}

// SegmentGroup represents a group of segments in a HLS manifest, usually separated by the #EXT-DISCONTINUITY tag.
type SegmentGroup struct {
	// List of segments in the group
	Segments []Segment
}

// Duration returns the total duration of the group, summing all durations from all segments.
func (g *SegmentGroup) Duration() float64 {
	var duration = 0.0

	for _, segment := range g.Segments {
		duration += float64(segment.Duration)
	}

	return duration
}

// MaxTargetDuration returns the target duration of the longest segment in the group, which is the duration rounded to the nearest integer.
func (g *SegmentGroup) MaxTargetDuration() uint8 {
	var maxTargetDuration uint8 = 0

	for _, segment := range g.Segments {
		targetDuration := segment.TargetDuration()

		if targetDuration > maxTargetDuration {
			maxTargetDuration = targetDuration
		}
	}

	return maxTargetDuration
}

// ExceedsTargetDuration returns true if any segment in the group has a target duration greater than the specified target duration.
func (g *SegmentGroup) ExceedsTargetDuration(targetDuration uint8) bool {
	for _, segment := range g.Segments {
		if segment.TargetDuration() > targetDuration {
			return true
		}
	}

	return false
}

// MaxDuration returns the duration of the longest segment in the group.
func (g *SegmentGroup) MaxDuration() float32 {
	var maxDuration float32 = 0.0

	for _, segment := range g.Segments {
		duration := segment.Duration

		if duration > maxDuration {
			maxDuration = duration
		}
	}

	return maxDuration
}

// Segment represents a segment in a HLS manifest.
type Segment struct {
	Path     string  // Path to the segment
	Duration float32 // Duration of the segment
	Title    string  // Title of the segment
}

// TargetDuration returns the target duration of the segment, which is the duration rounded to the nearest integer.
func (s *Segment) TargetDuration() uint8 {
	return uint8(math.Round(float64(s.Duration)))
}
