package vessel

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gofrs/uuid"
)

// func (d *VesselDeps) CreateService(ctx context.Context, in Vessel, header string) (*Vessel, error) {

// 	userKey, err := d.GetUserKeyRepo(ctx, header)
// 	if err != nil {
// 		return &Vessel{}, fmt.Errorf("error getting user key: %v", err.Error())
// 	}

// 	res, err := d.CreateRepo(ctx, in, *userKey)
// 	if err != nil {
// 		return &Vessel{}, fmt.Errorf("error creating data: %s", err.Error())
// 	}

// 	return res, nil
// }

func (d *VesselDeps) CreateNewService(ctx context.Context, in Vessel, header string) (*Vessel, error) {

	vesselUuid, err := uuid.NewV7()
	if err != nil {
		log.Fatalf("failed to generate UUID : %v", err)
		return &Vessel{}, fmt.Errorf("failed to generate UUID: %v", err)
	}

	userKey, err := d.GetUserKeyRepo(ctx, header)
	if err != nil {
		fmt.Printf("\nservice : %v \n", in)
		return &Vessel{}, fmt.Errorf("error getting user key: %v", err)
	}

	res, err := d.CreateNewRepo(ctx, in, *userKey, fmt.Sprint(&vesselUuid))
	if err != nil {
		return &Vessel{}, fmt.Errorf("error creating data: %s", err.Error())
	}

	return res, nil
}

// func (d *VesselDeps) FindByUserKeyService(ctx context.Context, userKey UserKey) (*Vessel, error) {
// 	res, err := d.FindByUserKeyRepo(ctx, userKey)
// 	if err != nil {
// 		return &Vessel{}, err
// 	}

// 	return res, nil
// }

func (d *VesselDeps) GetVesselService(ctx context.Context, vesselUuid, userKey string) (*Vessel, error) {
	vessel, err := d.GetVesselRepo(ctx, vesselUuid, userKey)
	if err != nil {
		return &Vessel{}, err
	}

	transponder, err := d.GetTransponderRepo(ctx, vesselUuid, userKey)
	if err != nil {
		return &Vessel{}, err
	}
	if transponder.ID != "" {
		vessel.Transponder = transponder
	}

	licenses, err := d.GetLicensesRepo(ctx, vesselUuid, userKey)
	if err != nil {
		return &Vessel{}, err
	}
	if len(licenses) > 0 {
		vessel.Licenses = licenses
	} else {
		vessel.Licenses = []Licenses{}
	}

	engines, err := d.GetEnginesRepo(ctx, vesselUuid, userKey)
	if err != nil {
		return &Vessel{}, err
	}
	if len(engines) > 0 {
		vessel.Engines = engines
	} else {
		vessel.Engines = []Engines{}
	}

	fishingCapacity, id, err := d.GetFishingCapacityRepo(ctx, vesselUuid, userKey)
	if err != nil {
		return &Vessel{}, err
	}
	if id != "" {
		vessel.FishingCapacity = fishingCapacity
	}

	ownerOperators, err := d.GetOwnerOperatorsRepo(ctx, vesselUuid, userKey)
	if err != nil {
		return &Vessel{}, err
	}
	if len(ownerOperators) > 0 {
		vessel.OwnerOperators = ownerOperators
	} else {
		vessel.OwnerOperators = []OwnerOperators{}
	}

	return vessel, nil
}

func (d *VesselDeps) GetAllService(ctx context.Context) ([]VesselGet, error) {
	record, err := d.GetAllRepo(ctx)
	if err != nil {
		if errors.Is(err, ErrConnPool) {
			return nil, err
		}
		if errors.Is(err, ErrQuery) {
			return nil, err
		}
		if errors.Is(err, ErrScan) {
			return nil, err
		}
		return nil, fmt.Errorf("error on service :%v", err)
	}
	return record, nil
}

func (d *VesselDeps) GetVesselByDateService(ctx context.Context, in SelectDateIn) ([]SelectDateOut, error) {
	out, err := d.GetVesselByDateRepo(ctx, in)
	if err != nil {
		if errors.Is(err, ErrConnPool) {
			return nil, err
		}
		if errors.Is(err, ErrQuery) {
			return nil, err
		}
		if errors.Is(err, ErrScan) {
			return nil, err
		}
		return nil, fmt.Errorf("error on service :%v", err)
	}
	return out, nil
}
