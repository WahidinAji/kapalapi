package vessel

// author generated I name it as vessel table
type Vessel struct {
	Name                 string           `json:"name"`
	Width                int              `json:"width"`
	Length               int              `json:"length"`
	Depth                int              `json:"depth"`
	Flag                 string           `json:"flag"`
	CallSign             string           `json:"callSign"`
	Type                 int              `json:"type"`
	Imo                  string           `json:"imo"`
	Registration         string           `json:"registration"`
	Mmsi                 string           `json:"mmsi"`
	PortOfRegistration   string           `json:"portOfRegistration"`
	ExternalMarking      string           `json:"externalMarking"`
	SatellitePhone       string           `json:"satellitePhone"`
	DscNumber            string           `json:"dscNumber"`
	MaxCrew              int              `json:"maxCrew"`
	HullMaterial         string           `json:"hullMaterial"`
	SternType            string           `json:"sternType"`
	Constructor          string           `json:"constructor"`
	GrossTonnage         float64          `json:"grossTonnage"`
	NetTonnage           float64          `json:"netTonnage"`
	RegionOfRegistration string           `json:"regionOfRegistration"`
	Transponder          Transponder      `json:"transponder"`
	Licenses             []Licenses       `json:"licenses"`
	Engines              []Engines        `json:"engines"`
	FishingCapacity      FishingCapacity  `json:"fishingCapacity"`
	OwnerOperators       []OwnerOperators `json:"ownerOperators"`
	PreferredImage       string           `json:"preferredImage"`
}

type Transponder struct {
	InstallDate             int64   `json:"installDate"`
	InstallCompany          string  `json:"installCompany"`
	InstallerName           string  `json:"installerName"`
	InstallationPort        string  `json:"installationPort"`
	InstallLatitude         float64 `json:"installLatitude"`
	InstallLongitude        float64 `json:"installLongitude"`
	InstallHeight           float64 `json:"installHeight"`
	VendorType              string  `json:"vendorType"`
	TransceiverManufacturer string  `json:"transceiverManufacturer"`
	ID                      string  `json:"id"`
	SerialNumber            string  `json:"serialNumber"`
}
type Licenses struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Expiry int64  `json:"expiry"`
}
type Engines struct {
	Power float64 `json:"power"`
	Type  string  `json:"type"`
}

type FishingCapacity struct {
	SubType                     int     `json:"subType"`
	GroupSeineFishing           bool    `json:"groupSeineFishing"`
	MainGear                    string  `json:"mainGear"`
	SubsidiaryGear              string  `json:"subsidiaryGear"`
	FreezerSnap                 bool    `json:"freezerSnap"`
	FreezerIce                  bool    `json:"freezerIce"`
	FreezerSeawaterRefrigerated bool    `json:"freezerSeawaterRefrigerated"`
	FreezerSeawaterChilled      bool    `json:"freezerSeawaterChilled"`
	FreezerBlastOrDry           bool    `json:"freezerBlastOrDry"`
	FreezerOther                bool    `json:"freezerOther"`
	FishingHoldCapacity         float64 `json:"fishingHoldCapacity"`
}

type OwnerOperators struct {
	Role           string `json:"role"`
	Name           string `json:"name"`
	Nationality    string `json:"nationality"`
	Address        string `json:"address"`
	Email          string `json:"email"`
	PhoneNumber1   string `json:"phoneNumber1"`
	PhoneNumber2   string `json:"phoneNumber2"`
	Mobile1        string `json:"mobile1"`
	Mobile2        string `json:"mobile2"`
	Current        bool   `json:"current"`
	PreferredImage string `json:"preferredImage"`
}
