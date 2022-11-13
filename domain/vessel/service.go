package vessel

import (
	"context"
	"fmt"
	"log"

	"github.com/gofrs/uuid"
)

func (d *VesselDeps) CreateService(ctx context.Context, in Vessel, header string) (*Vessel, error) {

	userKey, err := d.GetUserKeyRepo(ctx, header)
	if err != nil {
		return &Vessel{}, fmt.Errorf("error getting user key: %v", err.Error())
	}

	res, err := d.CreateRepo(ctx, in, *userKey)
	if err != nil {
		return &Vessel{}, fmt.Errorf("error creating data: %s", err.Error())
	}

	return res, nil
}

func (d *VesselDeps) CreateNewService(ctx context.Context, in Vessel, header string) (*Vessel, error) {

	vesselUuid, err := uuid.NewV7()
	if err != nil {
		log.Fatalf("failed to generate UUID : %v", err)
		return &Vessel{}, fmt.Errorf("failed to generate UUID: %v", err)
	}

	
	userKey, err := d.GetUserKeyRepo(ctx, header)
	if err != nil {
		return &Vessel{}, fmt.Errorf("error getting user key: %s", err.Error())
	}

	res, err := d.CreateNewRepo(ctx, in, *userKey, fmt.Sprint(&vesselUuid))
	if err != nil {
		return &Vessel{}, fmt.Errorf("error creating data: %s", err.Error())
	}

	return res, nil
}

func (d *VesselDeps) FindByUserKeyService(ctx context.Context, userKey UserKey) (*Vessel, error) {
	res, err := d.FindByUserKeyRepo(ctx, userKey)
	if err != nil {
		return &Vessel{}, err
	}

	return res, nil
}

func (d *VesselDeps) GetVesselService(ctx context.Context, vesselUuid, userKey string) (*Vessel, error) {
	res, err := d.GetVesselRepo(ctx, vesselUuid, userKey)
	if err != nil {
		return &Vessel{}, err
	}
	return res, nil
}
