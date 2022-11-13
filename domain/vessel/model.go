package vessel

import "github.com/jackc/pgx/v5"

type VesselDeps struct {
	DB *pgx.Conn
}

type VesselIn struct {
	Uuid string `json:"uuid"`
	Vessel
}

type UserKey struct {
	Id   string `json:"id" form:"id"`
	Uuid string `json:"uuid" form:"uuid"`
}

// author generated I name it as vessel table
type Vessel struct {
	Name                 string           `json:"name" form:"name"`
	Width                int              `json:"width" form:"width"`
	Length               int              `json:"length" form:"depth"`
	Depth                int              `json:"depth" form:""`
	Flag                 string           `json:"flag" form:"flag"`
	CallSign             string           `json:"callSign" form:"callSign"`
	Type                 int              `json:"type" form:"type"`
	Imo                  string           `json:"imo" form:"imo"`
	Registration         string           `json:"registration" form:"registration"`
	Mmsi                 string           `json:"mmsi" form:"mmsi"`
	PortOfRegistration   string           `json:"portOfRegistration" form:"portOfRegistration"`
	ExternalMarking      string           `json:"externalMarking" form:"externalMarking"`
	SatellitePhone       string           `json:"satellitePhone" form:"satellitePhone"`
	DscNumber            string           `json:"dscNumber" form:"dscNumber"`
	MaxCrew              int              `json:"maxCrew" form:"maxCrew"`
	HullMaterial         string           `json:"hullMaterial" form:"hullMaterial"`
	SternType            string           `json:"sternType" form:"sternType"`
	Constructor          string           `json:"constructor" form:"constructor"`
	GrossTonnage         float64          `json:"grossTonnage" form:"grossTonnage"`
	NetTonnage           float64          `json:"netTonnage" form:"netTonnage"`
	RegionOfRegistration string           `json:"regionOfRegistration" form:"regionOfRegistration"`
	Transponder          *Transponder      `jsonb:"transponder" form:"transponder"`
	Licenses             []Licenses       `jsonb:"licenses" form:"licenses"`
	Engines              []Engines        `jsonb:"engines" form:"engines"`
	FishingCapacity      *FishingCapacity  `jsonb:"fishingCapacity" form:"fishingCapacity"`
	OwnerOperators       []OwnerOperators `jsonb:"ownerOperators" form:"ownerOperators"`
	PreferredImage       string           `json:"preferredImage" form:"preferredImage"`
}

type Transponder struct {
	InstallDate             int64   `json:"installDate" form:"installDate"`
	InstallCompany          string  `json:"installCompany" form:"installCompany"`
	InstallerName           string  `json:"installerName" form:"installerName"`
	InstallationPort        string  `json:"installationPort" form:"installationPort"`
	InstallLatitude         float64 `json:"installLatitude" form:"installLatitude"`
	InstallLongitude        float64 `json:"installLongitude" form:"installLongitude"`
	InstallHeight           float64 `json:"installHeight" form:"installHeight"`
	VendorType              string  `json:"vendorType" form:"vendorType"`
	TransceiverManufacturer string  `json:"transceiverManufacturer" form:"transceiverManufacturer"`
	ID                      string  `json:"id" form:"id"`
	SerialNumber            string  `json:"serialNumber" form:"serialNumber"`
}
type Licenses struct {
	Type   string `json:"type" form:"type"`
	ID     string `json:"id" form:"id"`
	Expiry int64  `json:"expiry" form:"expiry"`
}
type Engines struct {
	Power float64 `json:"power" form:"power"`
	Type  string  `json:"type" form:"type"`
}

type FishingCapacity struct {
	SubType                     int     `json:"subType" form:"subType"`
	GroupSeineFishing           bool    `json:"groupSeineFishing" form:"groupSeineFishing"`
	MainGear                    string  `json:"mainGear" form:"mainGear"`
	SubsidiaryGear              string  `json:"subsidiaryGear" form:"subsidiaryGear"`
	FreezerSnap                 bool    `json:"freezerSnap" form:"freezerSnap"`
	FreezerIce                  bool    `json:"freezerIce" form:"freezerIce"`
	FreezerSeawaterRefrigerated bool    `json:"freezerSeawaterRefrigerated" form:"freezerSeawaterRefrigerated"`
	FreezerSeawaterChilled      bool    `json:"freezerSeawaterChilled" form:"freezerSeawaterChilled"`
	FreezerBlastOrDry           bool    `json:"freezerBlastOrDry" form:"freezerBlastOrDry"`
	FreezerOther                bool    `json:"freezerOther" form:"freezerOther"`
	FishingHoldCapacity         float64 `json:"fishingHoldCapacity" form:"fishingHoldCapacity"`
}

type OwnerOperators struct {
	Role           string `json:"role" form:"role"`
	Name           string `json:"name" form:"name"`
	Nationality    string `json:"nationality" form:"nationality"`
	Address        string `json:"address" form:"address"`
	Email          string `json:"email" form:"email"`
	PhoneNumber1   string `json:"phoneNumber1" form:"phoneNumber1"`
	PhoneNumber2   string `json:"phoneNumber2" form:"phoneNumber2"`
	Mobile1        string `json:"mobile1" form:"mobile1"`
	Mobile2        string `json:"mobile2" form:"mobile2"`
	Current        bool   `json:"current" form:"current"`
	PreferredImage string `json:"preferredImage" form:"preferredImage"`
}

type PreferredImage struct {
	VesselImage         string `json:"vesselImage"`
	OwnerOperatorsImage string `json:"ownerOperatorsImage"`
}
