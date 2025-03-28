package hls

import "testing"

func TestMerger(t *testing.T) {
	manifest0 := readManifest(t, "../testdata/stream0.m3u8")
	manifest1 := readManifest(t, "../testdata/stream1.m3u8")
	manifest2 := readManifest(t, "../testdata/stream2.m3u8")

	if !manifest0.IsCompatible(manifest2) {
		t.Errorf("expected manifest0 be compatible with manifest2")
	}

	if !manifest2.IsCompatible(manifest0) {
		t.Errorf("expected manifest2 be compatible with manifest0")
	}

	if manifest0.IsCompatible(manifest1) {
		t.Errorf("expected manifest0 not be compatible with manifest1")
	}

	if !manifest1.IsCompatible(manifest0) {
		t.Errorf("expected manifest1 be compatible with manifest0")
	}

	if !manifest1.IsCompatible(manifest2) {
		t.Errorf("expected manifest1 be compatible with manifest2")
	}

	if manifest2.IsCompatible(manifest1) {
		t.Errorf("expected manifest2 not be compatible with manifest1")
	}

	if !manifest0.Merge(manifest1) {
		t.Errorf("expected merge to have a breaking change")
	}

	if manifest0.Merge(manifest2) {
		t.Errorf("expected merge to not have a breaking change")
	}

	if len(manifest0.SegmentGroups) != 3 {
		t.Errorf("expected 3 segment groups, got %d", len(manifest0.SegmentGroups))
	}

	manifestString := manifest0.String()

	otherManifest, err := ParseHlsManifest(manifestString)
	if err != nil {
		t.Fatal(err)
	}

	otherManifestString := otherManifest.String()

	if manifestString != otherManifestString {
		t.Errorf("expected manifest to be the same, got different")
	}
}
