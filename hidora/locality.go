package hidora

// Zone is an availability zone
type Zone string

const (
	// ZoneChGen1 represents the ch-gen-1 zone
	ZoneChGen1 = Zone("ch-gen-1")
	// ZoneChVd1 represents the ch-vd-1 zone
	ZoneChVd1 = Zone("ch-vd-1")
	// ZoneChVd2 represents the ch-vd-2 zone
	ZoneChVd2 = Zone("ch-vd-2")
	// ZoneChVdTrial represents the ch-vd-trial zone
	ZoneChVdTrial = Zone("ch-vd-trial")
)

// Region is a geographical location
type Region string

const (
	// RegionChGen represents the ch-gen region
	RegionChGen = Region("ch-gen")
	// RegionChVd represents the ch-vd region
	RegionChVd = Region("ch-vd")
)
