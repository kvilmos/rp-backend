package service

import (
	"context"
	"room-planner/model"
	"room-planner/repository"
	"room-planner/request"

	"github.com/google/uuid"
)

type BlueprintService struct {
	BlueprintRepository repository.BlueprintRepository
	CornerRepository    repository.CornerRepository
	WallRepository      repository.WallRepository
	ItemRepository      repository.ItemRepository
}

func NewBlueprintService(bpr repository.BlueprintRepository, cr repository.CornerRepository, wr repository.WallRepository, ir repository.ItemRepository) *BlueprintService {
	return &BlueprintService{
		BlueprintRepository: bpr,
		CornerRepository:    cr,
		WallRepository:      wr,
		ItemRepository:      ir,
	}
}

func (s BlueprintService) CreateBlueprint(ctx context.Context, blueprintReq request.NewBlueprintRequest) (*model.Blueprint, error) {
	blueprint := model.Blueprint{
		UserId: blueprintReq.UserId,
	}

	err := s.BlueprintRepository.Create(ctx, &blueprint)
	if err != nil {
		return nil, err
	}

	return &blueprint, nil
}

func (s BlueprintService) GetBlueprintById(ctx context.Context, id int64) (*model.Blueprint, error) {
	return s.BlueprintRepository.GetById(ctx, id)
}

func (s BlueprintService) SaveBlueprint(ctx context.Context, blueprintReq request.NewBlueprintRequest) (*model.Blueprint, error) {
	blueprint := model.Blueprint{
		Id:     blueprintReq.Id,
		UserId: blueprintReq.UserId,
	}

	err := s.BlueprintRepository.Update(ctx, blueprint.Id)
	if err != nil {
		return nil, err
	}

	err = s.WallRepository.DeleteByBlueprintId(ctx, blueprint.Id)
	if err != nil {
		return nil, err
	}
	err = s.ItemRepository.DeleteByBlueprintId(ctx, blueprint.Id)
	if err != nil {
		return nil, err
	}
	err = s.CornerRepository.DeleteByBlueprintId(ctx, blueprint.Id)
	if err != nil {
		return nil, err
	}

	cornerMap := make(map[string]uuid.UUID)
	var corners []*model.Corner
	for _, cornerReq := range blueprintReq.Corners {
		id := cornerReq.Id

		corner := model.Corner{
			Id:          uuid.New(),
			BlueprintId: blueprintReq.Id,
			X:           cornerReq.X,
			Y:           cornerReq.Y,
		}

		corners = append(corners, &corner)
		cornerMap[id] = corner.Id
	}
	if corners != nil {
		err = s.CornerRepository.CreateMultiple(ctx, corners)
		if err != nil {
			return nil, err
		}
	}

	var walls []*model.Wall
	for _, wallReq := range blueprintReq.Walls {

		wall := model.Wall{
			StartCornerId: cornerMap[wallReq.StartCornerId],
			EndCornerId:   cornerMap[wallReq.EndCornerId],
		}

		walls = append(walls, &wall)
	}

	if walls != nil {
		err = s.WallRepository.CreateMultiple(ctx, walls)
		if err != nil {
			return nil, err
		}
	}

	var items []*model.Item
	for _, itemReq := range blueprintReq.Items {

		item := model.Item{
			BlueprintId: blueprintReq.Id,
			FurnitureId: itemReq.FurnitureId,
			PosX:        itemReq.PosX,
			PosY:        itemReq.PosY,
			PosZ:        itemReq.PosZ,
			Rot:         itemReq.Rot,
		}

		items = append(items, &item)
	}

	if items != nil {
		err = s.ItemRepository.CreateMultiple(ctx, items)
		if err != nil {
			return nil, err
		}
	}

	return &blueprint, nil
}
