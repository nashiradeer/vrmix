package hls

// Merge merges two manifests.
func (m *Manifest) Merge(m2 Manifest) bool {
	hasBreakingChange := false
	if m.TargetDuration < m2.TargetDuration {
		m.TargetDuration = m2.TargetDuration
		hasBreakingChange = true
	}

	if m.Version < m2.Version {
		m.Version = m2.Version
		hasBreakingChange = true
	}

	m.SegmentGroups = append(m.SegmentGroups, m2.SegmentGroups...)
	return hasBreakingChange
}

// IsCompatible checks if two manifests are compatible.
func (m *Manifest) IsCompatible(m2 Manifest) bool {
	return m.Version >= m2.Version && m.TargetDuration >= m2.TargetDuration
}
