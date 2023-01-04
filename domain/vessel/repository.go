package vessel

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (d *VesselDeps) GetUserKeyRepo(ctx context.Context, in string) (*UserKey, error) {
	fmt.Printf("repo get user : %v", in)
	var out UserKey
	res := d.PQ.QueryRowContext(ctx, `select id, uuid from user_keys where uuid=$1`, in).Scan(&out.Id, &out.Uuid)
	if res != nil {
		return &UserKey{}, fmt.Errorf("user key not found: %s", res.Error())
	}
	return &out, nil
}

// func (d *VesselDeps) CreateRepo(ctx context.Context, in Vessel, userKey UserKey) (*Vessel, error) {

// 	if err := d.DB.Ping(ctx); err != nil {
// 		return &Vessel{}, fmt.Errorf("connection error: %w", err)
// 	}

// 	tx, err := d.DB.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
// 	if err != nil {
// 		return &Vessel{}, fmt.Errorf("failed to start transaction: %w", err)
// 	}
// 	query := "insert into vessel (user_key_id,name,width,length,depth,flag,call_sign,type,imo,registration,mmsi,part_of_registration,external_marking,satellite_phone,dsc_number,max_crew,hull_material,stern_type,constructor,gross_tonnage,region_of_registration,transponder,licenses,engines,fishing_capacity,owner_operators,preferred_image,created_by) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28)"

// 	_, err = tx.Exec(ctx, query, &userKey.Id, &in.Name, &in.Width, &in.Length, &in.Depth, &in.Flag, &in.CallSign, &in.Type, &in.Imo, &in.Registration, &in.Mmsi, &in.PortOfRegistration, &in.ExternalMarking, &in.SatellitePhone, &in.DscNumber, &in.MaxCrew, &in.HullMaterial, &in.SternType, &in.Constructor, &in.GrossTonnage, &in.RegionOfRegistration, &in.Transponder, &in.Licenses, &in.Engines, &in.FishingCapacity, &in.OwnerOperators, &in.PreferredImage, &userKey.Uuid)

// 	if err != nil {
// 		if errRollback := tx.Rollback(ctx); errRollback != nil {
// 			return &Vessel{}, fmt.Errorf("rollback error while trying to migrate table : %v", errRollback)
// 		}
// 		return &Vessel{}, fmt.Errorf("failed to execute insert query: %v", err)
// 	}

// 	if err := tx.Commit(ctx); err != nil {
// 		if errRollback := tx.Rollback(ctx); errRollback != nil {
// 			return &Vessel{}, fmt.Errorf("rollback error while trying to migrate table : %v", errRollback)
// 		}
// 		return &Vessel{}, fmt.Errorf("failed to commit transaction: %w", err)
// 	}

// 	return &in, nil
// }

func (d *VesselDeps) CreateNewRepo(ctx context.Context, in Vessel, userKey UserKey, vesselUuid string) (*Vessel, error) {
	if err := d.PQ.PingContext(ctx); err != nil {
		return &Vessel{}, fmt.Errorf("connection error: %w", err)
	}

	tx, err := d.PQ.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return &Vessel{}, fmt.Errorf("failed to start transaction: %w", err)
	}

	// fmt.Printf(in.Transponder.SerialNumber, userKey, vesselUuid)

	// create to vessel
	query := `insert into vessel
	(user_key_id, created_by, name,width,length,depth,flag,call_sign,type,imo,registration,mmsi,part_of_registration,external_marking,satellite_phone,dsc_number,max_crew,hull_material,stern_type,constructor,gross_tonnage,region_of_registration,preferred_image, uuid) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24)`

	_, err = tx.ExecContext(ctx, query, &userKey.Id, &userKey.Uuid, &in.Name, &in.Width, &in.Length, &in.Depth, &in.Flag, &in.CallSign, &in.Type, &in.Imo, &in.Registration, &in.Mmsi, &in.PortOfRegistration, &in.ExternalMarking, &in.SatellitePhone, &in.DscNumber, &in.MaxCrew, &in.HullMaterial, &in.SternType, &in.Constructor, &in.GrossTonnage, &in.RegionOfRegistration, &in.PreferredImage, &vesselUuid)

	if err != nil {
		if errRoll := tx.Rollback(); errRoll != nil {
			return &Vessel{}, fmt.Errorf("rollback error while trying to insert into vessel failed: %v", err)
		}
		return &Vessel{}, fmt.Errorf("insert into vessel failed: %v", err)
	}

	//create to transponder
	if in.Transponder != nil {
		// query = `insert into transponder (vessel_uuid, install_date, install_company,installer_name,installion_port,install_latitude,install_longtitude,install_height,vendor_type,transceiver_manufacturer,serial_number)
		// values ('01846bbc-81d8-7aa6-9ea3-61f23a580003','1651849182769','SRT Installations','Boro','d1147ab5-d76d-11ec-90d2-0242ac1c0038',51.5,101.1,5.0,'OTHER','SRT','SER123456789');`
		_, err = tx.ExecContext(ctx, "insert into transponder (vessel_uuid, install_date, install_company,installer_name,installion_port,install_latitude,install_longtitude,install_height,vendor_type,transceiver_manufacturer,serial_number) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)", &vesselUuid, in.Transponder.InstallDate, &in.Transponder.InstallCompany, &in.Transponder.InstallerName, &in.Transponder.InstallationPort, &in.Transponder.InstallLatitude, &in.Transponder.InstallLongitude, &in.Transponder.InstallHeight, &in.Transponder.VendorType, &in.Transponder.TransceiverManufacturer, &in.Transponder.SerialNumber)
		// _, err = tx.Exec(ctx, query)

		if err != nil {
			if errRoll := tx.Rollback(); errRoll != nil {
				return &Vessel{}, fmt.Errorf("rollback error while trying to insert into transponder failed: %v", err)
			}
			return &Vessel{}, fmt.Errorf("insert into transponder failed: %v", err)
		}
	}
	//create to licenses
	if len(in.Licenses) > 0 {
		for i := 0; i < len(in.Licenses); i++ {
			_, err = tx.ExecContext(ctx, "insert into licenses (vessel_uuid, type, expiry) values ($1, $2, $3)", &vesselUuid, &in.Licenses[i].Type, &in.Licenses[i].Expiry)
			if err != nil {
				if errRoll := tx.Rollback(); errRoll != nil {
					return &Vessel{}, fmt.Errorf("rollback error while trying to insert into licenses failed: %v", err)
				}
				return &Vessel{}, fmt.Errorf("insert into licenses failed: %v", err)
			}
		}
	}

	//create to engines
	if len(in.Engines) > 0 {
		for i := 0; i < len(in.Engines); i++ {
			_, err = tx.ExecContext(ctx, "insert into engines (vessel_uuid, power, type) values ($1, $2, $3)", &vesselUuid, &in.Engines[i].Power, &in.Engines[i].Type)
			if err != nil {
				if errRoll := tx.Rollback(); errRoll != nil {
					return &Vessel{}, fmt.Errorf("rollback error while trying to insert into engines failed: %v", err)
				}
				return &Vessel{}, fmt.Errorf("insert into engines failed: %v", err)
			}
		}
	}

	//create to fishing_capacities
	if in.FishingCapacity != nil {
		_, err = tx.ExecContext(ctx, "insert into fishing_capacities (vessel_uuid, sub_type, group_seine_fishing,main_gear,subsidiary_gear,freezer_snap,freezer_ice,freezer_seawater_refrigerated,freezer_seawater_chilled,freezer_blast_or_dry,freezer_other,freezer_hold_capacity) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)", &vesselUuid, &in.FishingCapacity.SubType, &in.FishingCapacity.GroupSeineFishing, &in.FishingCapacity.MainGear, &in.FishingCapacity.SubsidiaryGear, &in.FishingCapacity.FreezerSnap, &in.FishingCapacity.FreezerIce, &in.FishingCapacity.FreezerSeawaterRefrigerated, &in.FishingCapacity.FreezerSeawaterChilled, &in.FishingCapacity.FreezerBlastOrDry, &in.FishingCapacity.FreezerOther, &in.FishingCapacity.FishingHoldCapacity)
		if err != nil {
			if errRoll := tx.Rollback(); errRoll != nil {
				return &Vessel{}, fmt.Errorf("rollback error while trying to insert into transponder failed: %v", err)
			}
			return &Vessel{}, fmt.Errorf("insert into fishing_capacities failed: %v", err)
		}
	}

	//create to owner_operators
	if len(in.OwnerOperators) > 0 {
		for i := 0; i < len(in.OwnerOperators); i++ {
			_, err = tx.ExecContext(ctx, "insert into owner_operators (vessel_uuid, role, nationality, address, email,phone_number_1,phone_number_2,mobile_1,mobile_2,current,preferred_image) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)", &vesselUuid, &in.OwnerOperators[i].Role, &in.OwnerOperators[i].Nationality, &in.OwnerOperators[i].Address, &in.OwnerOperators[i].Email, &in.OwnerOperators[i].PhoneNumber1, &in.OwnerOperators[i].PhoneNumber2, &in.OwnerOperators[i].Mobile1, &in.OwnerOperators[i].Mobile2, &in.OwnerOperators[i].Current, &in.OwnerOperators[i].PreferredImage)
			if err != nil {
				if errRoll := tx.Rollback(); errRoll != nil {
					return &Vessel{}, fmt.Errorf("rollback error while trying to insert into owner_operators failed: %v", err)
				}
				return &Vessel{}, fmt.Errorf("insert into owner_operators failed: %v", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return &Vessel{}, fmt.Errorf("rollback error while trying to migrate table : %v", errRollback)
		}
		return &Vessel{}, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &in, nil
}

// func (d *VesselDeps) FindByUserKeyRepo(ctx context.Context, userKey UserKey) (*Vessel, error) {

// 	var trans Transponder
// 	re := d.DB.QueryRow(ctx, "select transponder from vessel where id = 23").Scan(&trans.ID, &trans.VendorType, &trans.InstallDate, &trans.SerialNumber, &trans.InstallHeight, &trans.InstallerName, &trans.InstallCompany, &trans.InstallLatitude, &trans.InstallLongitude, &trans.InstallationPort, &trans.TransceiverManufacturer)
// 	if re.Error() == "" {
// 		return nil, errors.New("ini error")
// 	}

// 	fmt.Println("wkwk :", re)

// 	var vessels Vessel
// 	res, err := d.DB.Query(ctx, "SELECT name,width,length,depth,flag,call_sign,type,imo,registration,mmsi,part_of_registration,external_marking,satellite_phone,dsc_number,max_crew,hull_material,stern_type,constructor,gross_tonnage,region_of_registration,transponder,licenses,engines,fishing_capacity,owner_operators,preferred_image from vessel where id = $1 and created_by = $2", userKey.Id, userKey.Uuid)

// 	if err != nil {
// 		return &Vessel{}, fmt.Errorf("vessel not found: %v", err)
// 	}

// 	for res.Next() {
// 		var vessel Vessel
// 		errScan := res.Scan(&vessel.Name, &vessel.Width, &vessel.Length, &vessel.Depth, &vessel.Flag, &vessel.CallSign, &vessel.Type, &vessel.Imo, &vessel.Registration, &vessel.Mmsi, &vessel.PortOfRegistration, &vessel.ExternalMarking, &vessel.SatellitePhone, &vessel.DscNumber, &vessel.MaxCrew, &vessel.HullMaterial, &vessel.SternType, &vessel.Constructor, &vessel.GrossTonnage, &vessel.RegionOfRegistration, &vessel.Transponder, &vessel.Licenses, &vessel.Engines, &vessel.FishingCapacity, &vessel.OwnerOperators, &vessel.PreferredImage)
// 		if errScan != nil {
// 			return &Vessel{}, fmt.Errorf("error scanning to the vessel: %v", err)
// 		}
// 		vessels = vessel
// 	}

// 	log.Print(vessels.FishingCapacity)
// 	return &vessels, nil
// }

func (d *VesselDeps) GetVesselRepo(ctx context.Context, vesselUuid, userKey string) (*Vessel, error) {

	var checkVessel bool
	r := d.PQ.QueryRowContext(ctx, `select 1 from vessel where uuid = $1 and created_by = $2`, vesselUuid, userKey).Scan(&checkVessel)
	if r != nil {
		return &Vessel{}, r
	}
	if !checkVessel {
		return &Vessel{}, errors.New("vessel not found")
	}

	query := `select name,width,length,depth,flag,call_sign,type,imo,registration,mmsi,part_of_registration,external_marking,satellite_phone,dsc_number,max_crew,hull_material,stern_type,constructor,gross_tonnage,region_of_registration,preferred_image from vessel where uuid = $1 and created_by = $2`

	var vessels Vessel
	res, err := d.PQ.QueryContext(ctx, query, vesselUuid, userKey)
	if err != nil {
		return &Vessel{}, fmt.Errorf("vessel not found: %v", err)
	}
	// Scan(&vessel.Name, &vessel.Width, &vessel.Length, &vessel.Depth, &vessel.Flag, &vessel.CallSign, &vessel.Type, &vessel.Imo, &vessel.Registration, &vessel.Mmsi, &vessel.PortOfRegistration, &vessel.ExternalMarking, &vessel.SatellitePhone, &vessel.DscNumber, &vessel.MaxCrew, &vessel.HullMaterial, &vessel.SternType, &vessel.Constructor, &vessel.GrossTonnage, &vessel.RegionOfRegistration, &vessel.PreferredImage)
	// if res.Error() == "" {
	// 	return &Vessel{}, fmt.Errorf("error querying vessel: %v", res)
	// }

	for res.Next() {
		var vessel Vessel
		errScan := res.Scan(&vessel.Name, &vessel.Width, &vessel.Length, &vessel.Depth, &vessel.Flag, &vessel.CallSign, &vessel.Type, &vessel.Imo, &vessel.Registration, &vessel.Mmsi, &vessel.PortOfRegistration, &vessel.ExternalMarking, &vessel.SatellitePhone, &vessel.DscNumber, &vessel.MaxCrew, &vessel.HullMaterial, &vessel.SternType, &vessel.Constructor, &vessel.GrossTonnage, &vessel.RegionOfRegistration, &vessel.PreferredImage)
		if errScan != nil {
			return &Vessel{}, fmt.Errorf("error scanning to the vessel: %v", err)
		}
		vessels = vessel
	}

	// // fmt.Print("wwfsg")
	// //get transponder
	// query = `select install_date, install_company,installer_name,installion_port,install_latitude,install_longtitude,install_height,vendor_type,transceiver_manufacturer,serial_number from transponder where vessel_uuid = $1`
	// var transponder Transponder
	// trans := d.PQ.QueryRowContext(ctx, query, vesselUuid).Scan(&transponder.InstallDate, &transponder.InstallCompany, &transponder.InstallerName, &transponder.InstallationPort, &transponder.InstallLatitude, &transponder.InstallLongitude, &transponder.InstallHeight, &transponder.VendorType, &transponder.TransceiverManufacturer, &transponder.SerialNumber)
	// if trans != nil {
	// 	return &Vessel{}, fmt.Errorf("data transponder not found: %v", trans)
	// }
	// //get licenses
	// var licenses []Licenses
	// query = `select id, type, expiry from licenses where vessel_uuid = $1`
	// resLicenses, err := d.PQ.QueryContext(ctx, query, vesselUuid)
	// if err != nil {
	// 	return &Vessel{}, fmt.Errorf("data licenses not found: %v", err)
	// }
	// for resLicenses.Next() {
	// 	var license Licenses
	// 	if err := resLicenses.Scan(&license.ID, &license.Type, &license.Expiry); err != nil {
	// 		return &Vessel{}, fmt.Errorf("error scanning data licenses to struct: %v", err)
	// 	}
	// 	licenses = append(licenses, license)
	// }
	// //get engines
	// var engines []Engines
	// query = `select type, power from engines where vessel_uuid = $1`
	// resEngines, err := d.PQ.QueryContext(ctx, query, vesselUuid)
	// if err != nil {
	// 	return &Vessel{}, fmt.Errorf("data engines not found: %v", err)
	// }
	// for resEngines.Next() {
	// 	var engine Engines
	// 	if err := resEngines.Scan(&engine.Type, &engine.Power); err != nil {
	// 		return &Vessel{}, fmt.Errorf("error scanning data engines to struct: %v", err)
	// 	}
	// 	engines = append(engines, engine)
	// }
	// //get fishing_capacities
	// var fishingCapacity FishingCapacity
	// query = `select sub_type, group_seine_fishing,main_gear,subsidiary_gear,freezer_snap,freezer_ice,freezer_seawater_refrigerated,freezer_seawater_chilled,freezer_blast_or_dry,freezer_other,freezer_hold_capacity from fishing_capacities where vessel_uuid = $1`
	// fishing := d.PQ.QueryRowContext(ctx, query, vesselUuid).Scan(&fishingCapacity.SubType, &fishingCapacity.GroupSeineFishing, &fishingCapacity.MainGear, &fishingCapacity.SubsidiaryGear, &fishingCapacity.FreezerSnap, &fishingCapacity.FreezerIce, &fishingCapacity.FreezerSeawaterRefrigerated, &fishingCapacity.FreezerSeawaterChilled, &fishingCapacity.FreezerBlastOrDry, &fishingCapacity.FreezerOther, &fishingCapacity.FishingHoldCapacity)
	// if fishing != nil {
	// 	return &Vessel{}, fmt.Errorf("data fishing_capacities not found: %v", fishing)
	// }
	// //get ownerOperators
	// var ownerOperators []OwnerOperators
	// query = `select role, nationality, address, email,phone_number_1,phone_number_2,mobile_1,mobile_2,current,preferred_image from owner_operators where vessel_uuid = $1`
	// owner, err := d.PQ.QueryContext(ctx, query, vesselUuid)
	// if err != nil {
	// 	return &Vessel{}, fmt.Errorf("data owner_operators not found: %v", err)
	// }
	// for owner.Next() {
	// 	var ownerOperator OwnerOperators
	// 	if err := owner.Scan(&ownerOperator.Role, &ownerOperator.Nationality, &ownerOperator.Address, &ownerOperator.Email, &ownerOperator.PhoneNumber1, &ownerOperator.PhoneNumber2, &ownerOperator.Mobile1, &ownerOperator.Mobile2, &ownerOperator.Current, &ownerOperator.PreferredImage); err != nil {
	// 		return &Vessel{}, fmt.Errorf("error scanning data owner_operators to struct: %v", err)
	// 	}
	// 	ownerOperators = append(ownerOperators, ownerOperator)
	// }
	// vessels.Transponder = &transponder
	// vessels.Licenses = licenses
	// vessels.Engines = engines
	// vessels.FishingCapacity = &fishingCapacity
	// vessels.OwnerOperators = ownerOperators

	return &vessels, nil
}

// //get transponder
func (d *VesselDeps) GetTransponderRepo(ctx context.Context, vesselUuid, userKey string) (*Transponder, error) {
	query := `select install_date, install_company,installer_name,installion_port,install_latitude,install_longtitude,install_height,vendor_type,transceiver_manufacturer,serial_number from transponder where vessel_uuid = $1`

	var transponder Transponder
	res, err := d.PQ.QueryContext(ctx, query, vesselUuid)
	if err != nil {
		return &Transponder{}, fmt.Errorf("data transponder not found: %v", err)
	}
	for res.Next() {
		if err = res.Scan(&transponder.InstallDate, &transponder.InstallCompany, &transponder.InstallerName, &transponder.InstallationPort, &transponder.InstallLatitude, &transponder.InstallLongitude, &transponder.InstallHeight, &transponder.VendorType, &transponder.TransceiverManufacturer, &transponder.SerialNumber); err != nil {
			return &Transponder{}, fmt.Errorf("scanning transponder error: %v", err)
		}
	}
	// trans := d.DB.QueryRow(ctx, query, vesselUuid).Scan(&transponder.InstallDate, &transponder.InstallCompany, &transponder.InstallerName, &transponder.InstallationPort, &transponder.InstallLatitude, &transponder.InstallLongitude, &transponder.InstallHeight, &transponder.VendorType, &transponder.TransceiverManufacturer, &transponder.SerialNumber)
	// if trans != nil {
	// 	fmt.Print("wwfsg\n")
	// 	return &Transponder{}, fmt.Errorf("data transponder not found: %v", trans)
	// }
	return &transponder, nil
}

// //get licenses
func (d *VesselDeps) GetLicensesRepo(ctx context.Context, vesselUuid, userKey string) ([]Licenses, error) {
	var licenses []Licenses
	query := `select id, type, expiry from licenses where vessel_uuid = $1`
	resLicenses, err := d.PQ.QueryContext(ctx, query, vesselUuid)
	if err != nil {
		return []Licenses{}, fmt.Errorf("data licenses not found: %v", err)
	}
	for resLicenses.Next() {
		var license Licenses
		if err := resLicenses.Scan(&license.ID, &license.Type, &license.Expiry); err != nil {
			return []Licenses{}, fmt.Errorf("error scanning data licenses to struct: %v", err)
		}
		licenses = append(licenses, license)
	}
	return licenses, nil
}

// //get engines
func (d *VesselDeps) GetEnginesRepo(ctx context.Context, vesselUuid, userKey string) ([]Engines, error) {
	var engines []Engines
	query := `select type, power from engines where vessel_uuid = $1`
	resEngines, err := d.PQ.QueryContext(ctx, query, vesselUuid)
	if err != nil {
		return []Engines{}, fmt.Errorf("data engines not found: %v", err)
	}
	for resEngines.Next() {
		var engine Engines
		if err := resEngines.Scan(&engine.Type, &engine.Power); err != nil {
			return []Engines{}, fmt.Errorf("error scanning data engines to struct: %v", err)
		}
		engines = append(engines, engine)
	}
	return engines, nil
}

// //get fishing_capacities
func (d *VesselDeps) GetFishingCapacityRepo(ctx context.Context, vesselUuid, userKey string) (*FishingCapacity, string, error) {
	var fishingCapacity FishingCapacity
	var id string
	query := `select id, sub_type, group_seine_fishing,main_gear,subsidiary_gear,freezer_snap,freezer_ice,freezer_seawater_refrigerated,freezer_seawater_chilled,freezer_blast_or_dry,freezer_other,freezer_hold_capacity from fishing_capacities where vessel_uuid = $1`
	fishing := d.PQ.QueryRowContext(ctx, query, vesselUuid).Scan(&id, &fishingCapacity.SubType, &fishingCapacity.GroupSeineFishing, &fishingCapacity.MainGear, &fishingCapacity.SubsidiaryGear, &fishingCapacity.FreezerSnap, &fishingCapacity.FreezerIce, &fishingCapacity.FreezerSeawaterRefrigerated, &fishingCapacity.FreezerSeawaterChilled, &fishingCapacity.FreezerBlastOrDry, &fishingCapacity.FreezerOther, &fishingCapacity.FishingHoldCapacity)
	if fishing != nil {
		return &FishingCapacity{}, "", nil
	}
	return &fishingCapacity, id, nil
}

// //get ownerOperators
func (d *VesselDeps) GetOwnerOperatorsRepo(ctx context.Context, vesselUuid, userKey string) ([]OwnerOperators, error) {

	var ownerOperators []OwnerOperators
	query := `select role, nationality, address, email,phone_number_1,phone_number_2,mobile_1,mobile_2,current,preferred_image from owner_operators where vessel_uuid = $1`
	owner, err := d.PQ.QueryContext(ctx, query, vesselUuid)
	if err != nil {
		return []OwnerOperators{}, fmt.Errorf("data owner_operators not found: %v", err)
	}
	for owner.Next() {
		var ownerOperator OwnerOperators
		if err := owner.Scan(&ownerOperator.Role, &ownerOperator.Nationality, &ownerOperator.Address, &ownerOperator.Email, &ownerOperator.PhoneNumber1, &ownerOperator.PhoneNumber2, &ownerOperator.Mobile1, &ownerOperator.Mobile2, &ownerOperator.Current, &ownerOperator.PreferredImage); err != nil {
			return []OwnerOperators{}, fmt.Errorf("error scanning data owner_operators to struct: %v", err)
		}
		ownerOperators = append(ownerOperators, ownerOperator)
	}
	return ownerOperators, nil
}

func (d *VesselDeps) GetAllRepo(ctx context.Context) ([]VesselGet, error) {
	err := d.PQ.Ping()
	if err != nil {
		return nil, ErrConnPool
	}
	var vesselGet []VesselGet
	query := "select uk.id, uk.uuid as user_key_id, v.uuid as vessel_secret_key, v.created_at from user_keys as uk left join vessel as v on uk.id=v.user_key_id;"
	record, err := d.PQ.QueryContext(ctx, query)

	if err != nil {
		return nil, ErrQuery
	}

	if record.Next() {
		var vessel VesselGet
		if err := record.Scan(&vessel.Id, &vessel.UserKeyId, &vessel.VesselSecretKey, &vessel.CreatedAt); err != nil {
			return nil, ErrScan
		}
		vesselGet = append(vesselGet, vessel)
	}
	return vesselGet, nil
}

type VesselGet struct {
	Id              string `json:"id"`
	UserKeyId       string `json:"user_key_id"`
	VesselSecretKey string `json:"vessel_secret_key"`
	CreatedAt       string `json:"created_at"`
}

var (
	ErrConnPool = errors.New("error getting connetion pool")
	ErrQuery    = errors.New("error getting query")
	ErrScan     = errors.New("error scanning data")
)

//"query: select uk.id, uk.uuid as user_key_id, v.uuid as vessel_secret_key, v.created_at from user_keys as uk left join vessel as v"
