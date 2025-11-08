package service

import (
	"context"
	"room-planner/model"
	"room-planner/repository"
	"room-planner/request"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BlueprintService struct {
	Db                  *gorm.DB
	BlueprintRepository repository.BlueprintRepository
	CornerRepository    repository.CornerRepository
	WallRepository      repository.WallRepository
	ItemRepository      repository.ItemRepository
}

func NewBlueprintService(db *gorm.DB, bpr repository.BlueprintRepository, cr repository.CornerRepository, wr repository.WallRepository, ir repository.ItemRepository) *BlueprintService {
	return &BlueprintService{
		Db:                  db,
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

func (s BlueprintService) PageForUser(ctx context.Context, filter request.BlueprintFilter) ([]*model.Blueprint, int, error) {
	return s.BlueprintRepository.Page(ctx, filter)
}

func (s BlueprintService) GetBlueprintById(ctx context.Context, id int64) (*model.Blueprint, error) {
	return s.BlueprintRepository.GetById(ctx, id)
}

func (s BlueprintService) SaveBlueprint(ctx context.Context, blueprintReq request.NewBlueprintRequest) (*model.Blueprint, error) {
	blueprint := &model.Blueprint{
		Id:     blueprintReq.Id,
		Name:   blueprintReq.Name,
		UserId: blueprintReq.UserId,
	}

	tx := s.Db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	blueprintTx := s.BlueprintRepository.WithTx(tx)
	cornerTx := s.CornerRepository.WithTx(tx)
	wallTx := s.WallRepository.WithTx(tx)
	itemTx := s.ItemRepository.WithTx(tx)

	err := blueprintTx.Update(ctx, blueprint, blueprint.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = wallTx.DeleteByBlueprintId(ctx, blueprint.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = itemTx.DeleteByBlueprintId(ctx, blueprint.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = cornerTx.DeleteByBlueprintId(ctx, blueprint.Id)
	if err != nil {
		tx.Rollback()
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
		err = cornerTx.CreateMultiple(ctx, corners)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	var walls []*model.Wall
	for _, wallReq := range blueprintReq.Walls {

		wall := model.Wall{
			BlueprintId:   blueprint.Id,
			StartCornerId: cornerMap[wallReq.StartCornerId],
			EndCornerId:   cornerMap[wallReq.EndCornerId],
		}

		walls = append(walls, &wall)
	}

	if walls != nil {
		err = wallTx.CreateMultiple(ctx, walls)
		if err != nil {
			tx.Rollback()
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
		err = itemTx.CreateMultiple(ctx, items)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return blueprint, nil
}

func (s BlueprintService) GetCompleteBlueprintById(ctx context.Context, userId int64, id int64) (*model.Blueprint, error) {
	blueprint, err := s.BlueprintRepository.GetCompleteById(ctx, id)
	if err != nil {
		return nil, err
	}

	return blueprint, err
}

func (s BlueprintService) DeleteUserBlueprint(ctx context.Context, userId int64, blueprintId int64) error {
	return s.BlueprintRepository.DeleteForUser(ctx, userId, blueprintId)
}
