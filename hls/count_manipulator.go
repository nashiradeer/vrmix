package hls

import "slices"

// RemoveFromStart removes n segments from the start of the manifest returning the count of segments group and segments removed and updating the media sequence and discontinuity sequence.
func (m *Manifest) RemoveFromStart(n int) (int, int) {
	var newGroups []SegmentGroup
	segmentsRemoved := 0
	segmentsGroupRemoved := 0

	for _, group := range m.SegmentGroups {
		if n <= 0 {
			newGroups = append(newGroups, group)
			continue
		}

		removed := group.RemoveFromStart(n)
		segmentsRemoved += removed
		n -= removed

		if len(group.Segments) == 0 {
			m.DiscontinuitySequence += 1
			m.MediaSequence += uint32(removed)
			segmentsGroupRemoved += 1
			continue
		}

		newGroups = append(newGroups, group)
		m.MediaSequence += uint32(removed)
	}

	m.SegmentGroups = newGroups
	return segmentsGroupRemoved, segmentsRemoved
}

// RemoveFromEnd removes n segments from the end of the manifest returning the count of segments group and segments removed.
func (m *Manifest) RemoveFromEnd(n int) (int, int) {
	var newGroups []SegmentGroup
	segmentsRemoved := 0
	segmentsGroupRemoved := 0

	for _, group := range slices.Backward(m.SegmentGroups) {
		if n <= 0 {
			newGroups = append([]SegmentGroup{group}, newGroups...)
			continue
		}

		removed := group.RemoveFromEnd(n)
		segmentsRemoved += removed
		n -= removed

		if len(group.Segments) == 0 {
			segmentsGroupRemoved += 1
			m.DiscontinuitySequence += 1
			continue
		}

		newGroups = append([]SegmentGroup{group}, newGroups...)
	}

	m.SegmentGroups = newGroups
	m.MediaSequence += uint32(segmentsRemoved)
	return segmentsGroupRemoved, segmentsRemoved
}

// RemoveFromStart removes n segments from the start of the segment group returning the count of segments removed, this will not update the media sequence or discontinuity sequence so it is not recommended to use this method directly.
func (g *SegmentGroup) RemoveFromStart(n int) int {
	if n <= 0 {
		return 0
	}
	oldLen := len(g.Segments)

	if n >= len(g.Segments) {
		g.Segments = []Segment{}
		return oldLen
	}

	g.Segments = g.Segments[n:]
	return oldLen - len(g.Segments)
}

// RemoveFromEnd removes n segments from the end of the segment group returning the count of segments removed.
func (g *SegmentGroup) RemoveFromEnd(n int) int {
	if n <= 0 {
		return 0
	}

	oldLen := len(g.Segments)

	if n >= len(g.Segments) {
		g.Segments = []Segment{}
		return oldLen
	}

	g.Segments = g.Segments[:len(g.Segments)-n]
	return oldLen - len(g.Segments)
}
