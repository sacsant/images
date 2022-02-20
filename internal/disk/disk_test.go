package disk

import (
	"math/rand"
	"testing"

	"github.com/osbuild/osbuild-composer/internal/blueprint"
	"github.com/stretchr/testify/assert"
)

func TestDisk_AlignUp(t *testing.T) {

	pt := PartitionTable{}
	firstAligned := DefaultGrainBytes

	tests := []struct {
		size uint64
		want uint64
	}{
		{0, 0},
		{1, firstAligned},
		{firstAligned - 1, firstAligned},
		{firstAligned, firstAligned}, // grain is already aligned => no change
		{firstAligned / 2, firstAligned},
		{firstAligned + 1, firstAligned * 2},
	}

	for _, tt := range tests {
		got := pt.AlignUp(tt.size)
		assert.Equal(t, tt.want, got, "Expected %d, got %d", tt.want, got)
	}
}

func TestDisk_DynamicallyResizePartitionTable(t *testing.T) {
	mountpoints := []blueprint.FilesystemCustomization{
		{
			MinSize:    2147483648,
			Mountpoint: "/usr",
		},
	}
	pt := PartitionTable{
		UUID: "D209C89E-EA5E-4FBD-B161-B461CCE297E0",
		Type: "gpt",
		Partitions: []Partition{
			{
				Size:     2048,
				Bootable: true,
				Type:     BIOSBootPartitionGUID,
				UUID:     BIOSBootPartitionUUID,
			},
			{
				Type: FilesystemDataGUID,
				UUID: RootPartitionUUID,
				Payload: &Filesystem{
					Type:         "xfs",
					Label:        "root",
					Mountpoint:   "/",
					FSTabOptions: "defaults",
					FSTabFreq:    0,
					FSTabPassNo:  0,
				},
			},
		},
	}
	var expectedSize uint64 = 2147483648
	// math/rand is good enough in this case
	/* #nosec G404 */
	rng := rand.New(rand.NewSource(0))
	pt, err := CreatePartitionTable(mountpoints, 1024, &pt, rng)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, pt.SectorsToBytes(pt.Size), expectedSize)
}

// common partition table that use used by tests
var canonicalPartitionTable = PartitionTable{
	UUID: "D209C89E-EA5E-4FBD-B161-B461CCE297E0",
	Type: "gpt",
	Partitions: []Partition{
		{
			Size:     2048,
			Bootable: true,
			Type:     BIOSBootPartitionGUID,
			UUID:     BIOSBootPartitionUUID,
		},
		{
			Size: 204800,
			Type: EFISystemPartitionGUID,
			UUID: EFISystemPartitionUUID,
			Payload: &Filesystem{
				Type:         "vfat",
				UUID:         EFIFilesystemUUID,
				Mountpoint:   "/boot/efi",
				FSTabOptions: "defaults,uid=0,gid=0,umask=077,shortname=winnt",
				FSTabFreq:    0,
				FSTabPassNo:  2,
			},
		},
		{
			Size: 1048576,
			Type: FilesystemDataGUID,
			UUID: FilesystemDataUUID,
			Payload: &Filesystem{
				Type:         "xfs",
				Mountpoint:   "/boot",
				FSTabOptions: "defaults",
				FSTabFreq:    0,
				FSTabPassNo:  0,
			},
		},
		{
			Type: FilesystemDataGUID,
			UUID: RootPartitionUUID,
			Payload: &Filesystem{
				Type:         "xfs",
				Label:        "root",
				Mountpoint:   "/",
				FSTabOptions: "defaults",
				FSTabFreq:    0,
				FSTabPassNo:  0,
			},
		},
	},
}

func TestDisk_ForEachFilesystem(t *testing.T) {
	rootFs := canonicalPartitionTable.Partitions[3].Payload
	bootFs := canonicalPartitionTable.Partitions[2].Payload
	efiFs := canonicalPartitionTable.Partitions[1].Payload

	// check we iterate in the correct order and throughout the whole array
	var expectedFilesystems []*Filesystem
	err := canonicalPartitionTable.ForEachFilesystem(func(fs *Filesystem) error {
		expectedFilesystems = append(expectedFilesystems, fs)
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, []*Filesystem{efiFs, bootFs, rootFs}, expectedFilesystems)

	// check we stop iterating when the callback returns false
	expectedFilesystems = make([]*Filesystem, 0)
	err = canonicalPartitionTable.ForEachFilesystem(func(fs *Filesystem) error {
		if fs.Mountpoint != "/boot" {
			return nil
		}

		// we should stop at boot, never reaching root
		assert.NotEqual(t, fs.Mountpoint, "/")

		expectedFilesystems = append(expectedFilesystems, fs)
		return StopIter
	})

	assert.NoError(t, err)
	assert.Equal(t, []*Filesystem{bootFs}, expectedFilesystems)
}

func TestDisk_ForEachEntity(t *testing.T) {

	count := 0
	err := canonicalPartitionTable.ForEachEntity(func(e Entity, path []Entity) error {
		assert.NotNil(t, e)
		assert.NotNil(t, path)

		count += 1
		return nil
	})

	assert.NoError(t, err)

	// PartitionTable, 4 partitions, 3 filesystems -> 8 entities
	assert.Equal(t, 8, count)

}
