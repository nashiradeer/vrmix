package hls

import "testing"

func TestRemoveEnd(t *testing.T) {
	manifest0 := readManifest(t, "../testdata/stream0.m3u8")
	manifest1 := readManifest(t, "../testdata/stream1.m3u8")

	manifest1.Merge(manifest0)

	oldSegmentCount := manifest1.SegmentCount()

	manifest1.RemoveFromEnd(1)

	if oldSegmentCount != manifest1.SegmentCount()+1 {
		t.Errorf("expected segment count to be %d, got %d", oldSegmentCount-1, manifest1.SegmentCount())
	}

	oldSegmentCount = manifest1.SegmentCount()

	manifest1.RemoveFromEnd(3)

	if oldSegmentCount != manifest1.SegmentCount()+3 {
		t.Errorf("expected segment count to be %d, got %d", oldSegmentCount-3, manifest1.SegmentCount())
	}
}

func TestRemoveStart(t *testing.T) {
	manifest0 := readManifest(t, "../testdata/stream0.m3u8")
	manifest1 := readManifest(t, "../testdata/stream1.m3u8")

	manifest0.Merge(manifest1)

	oldSegmentCount := manifest0.SegmentCount()

	manifest0.RemoveFromStart(1)

	if oldSegmentCount != manifest0.SegmentCount()+1 {
		t.Errorf("expected segment count to be %d, got %d", oldSegmentCount-1, manifest0.SegmentCount())
	}

	if manifest0.MediaSequence != 1 {
		t.Errorf("expected media sequence to be 1, got %d", manifest0.MediaSequence)
	}

	if manifest0.DiscontinuitySequence != 0 {
		t.Errorf("expected discontinuity sequence to be 0, got %d", manifest0.DiscontinuitySequence)
	}

	oldSegmentCount = manifest0.SegmentCount()

	manifest0.RemoveFromStart(3)

	if oldSegmentCount != manifest0.SegmentCount()+3 {
		t.Errorf("expected segment count to be %d, got %d", oldSegmentCount-3, manifest0.SegmentCount())
	}

	if manifest0.MediaSequence != 4 {
		t.Errorf("expected media sequence to be 4, got %d", manifest0.MediaSequence)
	}

	if manifest0.DiscontinuitySequence != 1 {
		t.Errorf("expected discontinuity sequence to be 1, got %d", manifest0.DiscontinuitySequence)
	}
}
