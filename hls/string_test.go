package hls

import (
	"os"
	"strconv"
	"testing"
)

// readManifest reads a manifest file and returns the parsed manifest
func readManifest(t *testing.T, path string) Manifest {
	rawData, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	data := string(rawData)

	manifest, err := ParseHlsManifest(data)
	if err != nil {
		t.Fatal(err)
	}

	return manifest
}

// testHeader tests the header of a manifest
func testHeader(t *testing.T, manifest Manifest, targetDuration uint8) {
	if manifest.Version != 3 {
		t.Errorf("expected version 3, got %d", manifest.Version)
	}

	if manifest.TargetDuration != targetDuration {
		t.Errorf("expected target duration %d, got %d", targetDuration, manifest.TargetDuration)
	}

	if manifest.MediaSequence != 0 {
		t.Errorf("expected media sequence 0, got %d", manifest.MediaSequence)
	}

	if manifest.DiscontinuitySequence != 0 {
		t.Errorf("expected discontinuity sequence 0, got %d", manifest.DiscontinuitySequence)
	}

	if !manifest.HasEndList {
		t.Errorf("expected end list to be true, got false")
	}

	if !manifest.HasValidTargetDuration() {
		t.Errorf("expected valid target duration, got invalid")
	}

	if len(manifest.SegmentGroups) != 1 {
		t.Errorf("expected 1 segment group, got %d", len(manifest.SegmentGroups))
	}
}

// testSegments tests each segment in a segment group against the expected duration
func testSegments(t *testing.T, segments []Segment, durations []float32) {
	if len(segments) != len(durations) {
		t.Errorf("expected %d segments, got %d", len(durations), len(segments))
	}

	for i, s := range segments {
		path := strconv.Itoa(i) + ".ts"
		if s.Path != path {
			t.Errorf("expected path %s, got %s", path, s.Path)
		}

		duration := durations[i]
		if s.Duration != duration {
			t.Errorf("expected duration %f, got %f", duration, s.Duration)
		}

		if s.Title != "" {
			t.Errorf("expected title \"\", got %s", s.Title)
		}
	}
}

// testToString tests that a manifest can be converted to a string and back
func testToString(t *testing.T, manifest Manifest) {
	manifestString := manifest.String()

	otherManifest, err := ParseHlsManifest(manifestString)
	if err != nil {
		t.Fatal(err)
	}

	otherManifestString := otherManifest.String()

	if manifestString != otherManifestString {
		t.Errorf("expected manifest to be the same, got different")
	}
}

func TestStream0(t *testing.T) {
	m := readManifest(t, "../testdata/stream0.m3u8")

	testHeader(t, m, 4)

	segmentDuration := []float32{4.166667, 3.483333}
	testSegments(t, m.SegmentGroups[0].Segments, segmentDuration)

	testToString(t, m)
}

func TestStream1(t *testing.T) {
	m := readManifest(t, "../testdata/stream1.m3u8")

	testHeader(t, m, 8)

	segmentsDuration := []float32{8.333333, 8.333333, 5.233333}
	testSegments(t, m.SegmentGroups[0].Segments, segmentsDuration)

	testToString(t, m)
}

func TestStream2(t *testing.T) {
	m := readManifest(t, "../testdata/stream2.m3u8")

	testHeader(t, m, 4)

	segmentsDuration := []float32{3.916667, 4.166667, 1.166667, 2.166667, 0.900000, 0.433333, 1.300000, 1.966667, 2.166667, 2.166667, 1.733333, 2.166667, 1.933333, 2.066667, 2.100000, 1.733333, 2.333333}
	testSegments(t, m.SegmentGroups[0].Segments, segmentsDuration)

	testToString(t, m)
}
