package service

import (
	"fmt"
	"github.com/curltech/go-colla-biz/actual/entity"
	"github.com/curltech/go-colla-core/container"
	baseentity "github.com/curltech/go-colla-core/entity"
	"github.com/curltech/go-colla-core/service"
	"github.com/curltech/go-colla-core/util/message"
)

/**
同步表结构，服务继承基本服务的方法
*/
type RoleService struct {
	service.OrmBaseService
}

var roleService = &RoleService{}

func GetRoleService() *RoleService {
	return roleService
}

var seqname = "seq_actual"

func (this *RoleService) GetSeqName() string {
	return seqname
}

func (this *RoleService) NewEntity(data []byte) (interface{}, error) {
	entity := &entity.Role{}
	if data == nil {
		return entity, nil
	}
	err := message.Unmarshal(data, entity)
	if err != nil {
		return nil, err
	}

	return entity, err
}

func (this *RoleService) NewEntities(data []byte) (interface{}, error) {
	entities := make([]*entity.Role, 0)
	if data == nil {
		return &entities, nil
	}
	err := message.Unmarshal(data, &entities)
	if err != nil {
		return nil, err
	}

	return &entities, err
}

/**
当LoadNum=0，只加载LoadNum=0的数据；当LoadNum=-1，加载所有数据
*/
func (this *RoleService) FindOffspring(schemaName string, parentIds []uint64, LoadNum int) ([]*entity.Role, []uint64) {
	if parentIds != nil && len(parentIds) > 0 {
		roles := make([]*entity.Role, 0)
		var conds string
		if LoadNum == -1 {
			conds = fmt.Sprintf("Id in (%s)", service.PlaceQuestionMark(len(parentIds)))
		} else {
			conds = fmt.Sprintf("Id in (%s) and LoadNum=0", service.PlaceQuestionMark(len(parentIds)))
		}

		params := make([]interface{}, len(parentIds))
		for i, id := range parentIds {
			params[i] = id
		}
		err := this.Find(&roles, nil, baseentity.FieldName_Id, 0, 0, conds, params...)
		if err == nil && len(roles) > 0 {
			var pids = make([]uint64, len(roles))
			var ids = make([]uint64, len(roles))
			for i, r := range roles {
				pids[i] = r.Id
				ids[i] = r.Id
			}
		loop:
			rs := make([]*entity.Role, 0)
			idss := Split(ids)
			ids = make([]uint64, 0)
			for _, is := range idss {
				conds = fmt.Sprintf("ParentId in (%s) and LoadNum=0", service.PlaceQuestionMark(len(is)))
				params = make([]interface{}, len(is))
				for i, id := range is {
					params[i] = id
				}
				err = this.Find(&rs, nil, baseentity.FieldName_Id, 0, 0, conds, params...)
				if err == nil && len(rs) > 0 {
					for _, r := range rs {
						roles = append(roles, r)
						pids = append(pids, r.Id)
						ids = append(ids, r.Id)
					}
				}
			}
			if len(ids) > 0 {
				goto loop
			}
			return roles, pids
		}
	}

	return nil, nil
}

func (this *RoleService) FindSpecOffspring(schemaName string, parentId uint64, specId uint64) ([]*entity.Role, []uint64) {
	if parentId > 0 && specId > 0 {
		roles := make([]*entity.Role, 0)
		condiBean := &entity.Role{}
		condiBean.ParentId = parentId
		condiBean.SpecId = specId
		err := this.Find(&roles, condiBean, baseentity.FieldName_Id, 0, 0, "")
		if err == nil && len(roles) > 0 {
			var pids = make([]uint64, len(roles))
			var ids = make([]uint64, len(roles))
			for i, r := range roles {
				pids[i] = r.Id
				ids[i] = r.Id
			}
		loop:
			rs := make([]*entity.Role, 0)
			idss := Split(ids)
			ids = make([]uint64, 0)
			for _, is := range idss {
				conds := fmt.Sprintf("ParentId in (%s)", service.PlaceQuestionMark(len(is)))
				params := make([]interface{}, len(is))
				for i, id := range is {
					params[i] = id
				}
				err = this.Find(&rs, nil, baseentity.FieldName_Id, 0, 0, conds, params...)
				if err == nil && len(rs) > 0 {
					for _, r := range rs {
						roles = append(roles, r)
						pids = append(pids, r.Id)
						ids = append(ids, r.Id)
					}
				}
			}
			if len(ids) > 0 {
				goto loop
			}
			return roles, pids
		}
	}

	return nil, nil
}

func Split(parentIds []uint64) [][]uint64 {
	var is = make([][]uint64, 0)
	for i := 0; i < len(parentIds); i = i + 999 {
		start := i * 999
		end := 999
		if end > len(parentIds)-start {
			end = len(parentIds) - start
		}
		ids := parentIds[start:end]
		is = append(is, ids)
	}

	return is
}

func init() {
	service.GetSession().Sync(new(entity.Role))
	roleService.OrmBaseService.GetSeqName = roleService.GetSeqName
	roleService.OrmBaseService.FactNewEntity = roleService.NewEntity
	roleService.OrmBaseService.FactNewEntities = roleService.NewEntities
	service.RegistSeq(seqname, 0)
	container.RegistService("atlRole", roleService)
}
