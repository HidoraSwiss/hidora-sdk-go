package hidora

// Region is a geographical location
type HidoraRegion string

const (
	// RegionChGen represents the ch-gen region
	RegionChGen = HidoraRegion("ch-gen")
	// RegionChVd represents the ch-vd region
	RegionChVd = HidoraRegion("ch-vd")
)

// Zone is an availability zone
type HidoraZone string

const (
	// ZoneChGen1 represents the ch-gen-1 zone
	ZoneChGen1 = HidoraZone("ch-gen-1")
	// ZoneChVd1 represents the ch-vd-1 zone
	ZoneChVd1 = HidoraZone("ch-vd-1")
	// ZoneChVd2 represents the ch-vd-2 zone
	ZoneChVd2 = HidoraZone("ch-vd-2")
	// ZoneChVdTrial represents the ch-vd-trial zone
	ZoneChVdTrial = HidoraZone("ch-vd-trial")
)
