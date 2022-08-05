package main

import (
	"context"
	pb "github.com/Jimmy01010/protocol/vessel-service"
)

type handler struct {
	repository
}

func (h *handler) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {

	// Find the next available vessel
	Vessel, err := h.repository.FindAvailable(ctx, MarshalSpecification(req))
	if err != nil {
		return err
	}

	res.Vessel = UnmarshalVessel(Vessel)
	return nil
}

func (h *handler) Create(ctx context.Context, req *pb.Vessel, res *pb.Response) error {
	if err := h.repository.Create(ctx, MarshalVessel(req)); err != nil {
		return err
	}

	res.Vessel = req
	return nil
}
