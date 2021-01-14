package spec

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/curltech/go-colla-biz/spec/entity"
	"github.com/curltech/go-colla-biz/spec/service"
	baseentity "github.com/curltech/go-colla-core/entity"
	"github.com/curltech/go-colla-core/logger"
	"github.com/curltech/go-colla-core/util/convert"
	"github.com/curltech/go-colla-core/util/reflect"
	"github.com/huandu/xstrings"
)

func UploadExcel(filename string) error {
	xlsx, err := excelize.OpenFile(filename)
	if err != nil {
		logger.Errorf("filename:%v can't open", filename)
		return err
	}

	rows := xlsx.GetRows("spec_role")
	roleSpec(rows)

	rows = xlsx.GetRows("spec_fixedrole")
	go fixedRoleSpec(rows)

	rows = xlsx.GetRows("spec_action")
	go actionSpec(rows)

	rows = xlsx.GetRows("spec_connection")
	go connectionSpec(rows)

	count := xlsx.SheetCount
	for c := 0; c < count; c++ {
		sheetname := xlsx.GetSheetName(c + 1)
		if sheetname != "spec_role" || sheetname != "spec_fixedrole" || sheetname != "spec_action" || sheetname != "spec_connection" {
			rows = xlsx.GetRows(sheetname)
			go attibuteSpec(sheetname, rows)
		}
	}

	return nil
}

func roleSpec(rows [][]string) {
	var head []string
	var pos = -1
	roleSpecs := make([]interface{}, 0)
	svc := service.GetRoleSpecService()
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		if i == 0 {
			head = row
			for j := 0; j < len(row); j++ {
				if "SpecId" == xstrings.FirstRuneToUpper(row[j]) {
					pos = j
				}
			}
			continue
		}
		var specId uint64
		v, err := convert.ToObject(row[pos], "uint64")
		if err == nil {
			specId = v.(uint64)
		}
		roleSpec := entity.RoleSpec{}
		roleSpec.SpecId = specId
		roleSpec.Status = baseentity.EntityStatus_Effective
		svc.Get(&roleSpec, false, "", "")
		for j := 0; j < len(row); j++ {
			fieldname := xstrings.FirstRuneToUpper(head[j])
			if fieldname != "" {
				reflect.Set(&roleSpec, fieldname, row[j])
			}
		}
		roleSpecs = append(roleSpecs, &roleSpec)
	}
	svc.Upsert(roleSpecs...)
}

func actionSpec(rows [][]string) {
	var head []string
	var pos = -1
	actionSpecs := make([]interface{}, 0)
	svc := service.GetActionSpecService()
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		if i == 0 {
			head = row
			for j := 0; j < len(row); j++ {
				if "SpecId" == xstrings.FirstRuneToUpper(row[j]) {
					pos = j
				}
			}
			continue
		}
		var specId uint64
		v, err := convert.ToObject(row[pos], "uint64")
		if err == nil {
			specId = v.(uint64)
		}
		actionSpec := entity.ActionSpec{}
		actionSpec.SpecId = specId
		actionSpec.Status = baseentity.EntityStatus_Effective
		svc.Get(&actionSpec, false, "", "")
		for j := 0; j < len(row); j++ {
			fieldname := xstrings.FirstRuneToUpper(head[j])
			if fieldname != "" {
				reflect.Set(&actionSpec, fieldname, row[j])
			}
		}
		actionSpecs = append(actionSpecs, &actionSpec)
	}
	svc.Upsert(actionSpecs...)
}

func fixedRoleSpec(rows [][]string) {
	var head []string
	var pos = -1
	fixedRoleSpecs := make([]interface{}, 0)
	svc := service.GetFixedRoleSpecService()
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		if i == 0 {
			head = row
			for j := 0; j < len(row); j++ {
				if "SpecId" == xstrings.FirstRuneToUpper(row[j]) {
					pos = j
				}
			}
			continue
		}
		var specId uint64
		v, err := convert.ToObject(row[pos], "uint64")
		if err == nil {
			specId = v.(uint64)
		}
		fixedRoleSpec := entity.FixedRoleSpec{}
		fixedRoleSpec.SpecId = specId
		fixedRoleSpec.Status = baseentity.EntityStatus_Effective
		svc.Get(&fixedRoleSpec, false, "", "")
		for j := 0; j < len(row); j++ {
			fieldname := xstrings.FirstRuneToUpper(head[j])
			if fieldname != "" {
				reflect.Set(&fixedRoleSpec, fieldname, row[j])
			}
		}
		fixedRoleSpecs = append(fixedRoleSpecs, &fixedRoleSpec)
	}
	svc.Upsert(fixedRoleSpecs...)
}

func connectionSpec(rows [][]string) {
	var head []string
	var parentPos = -1
	var subPos = -1
	connectionSpecs := make([]interface{}, 0)
	svc := service.GetConnectionSpecService()
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		if i == 0 {
			head = row
			for j := 0; j < len(row); j++ {
				if "ParentSpecId" == xstrings.FirstRuneToUpper(row[j]) {
					parentPos = j
				}
				if "SubSpecId" == xstrings.FirstRuneToUpper(row[j]) {
					subPos = j
				}
			}
			continue
		}
		var parentSpecId uint64
		var subSpecId uint64
		v, err := convert.ToObject(row[parentPos], "uint64")
		if err == nil {
			parentSpecId = v.(uint64)
		}
		v, err = convert.ToObject(row[subPos], "uint64")
		if err == nil {
			subSpecId = v.(uint64)
		}
		connectionSpec := entity.ConnectionSpec{}
		connectionSpec.ParentSpecId = parentSpecId
		connectionSpec.SubSpecId = subSpecId
		connectionSpec.Status = baseentity.EntityStatus_Effective
		svc.Get(&connectionSpec, false, "", "")
		for j := 0; j < len(row); j++ {
			fieldname := xstrings.FirstRuneToUpper(head[j])
			if fieldname != "" {
				reflect.Set(&connectionSpec, fieldname, row[j])
			}
		}
		connectionSpec.SpecType = entity.SpecType_Role
		if connectionSpec.RelationType == "" {
			connectionSpec.RelationType = entity.RelationType_Dependency
		}
		connectionSpecs = append(connectionSpecs, &connectionSpec)
	}
	svc.Upsert(connectionSpecs...)
}

func attibuteSpec(sheetname string, rows [][]string) {
	roleSpec := entity.RoleSpec{}
	roleSpec.Kind = sheetname
	roleSpec.Status = baseentity.EntityStatus_Effective
	roleSpecSvc := service.GetRoleSpecService()
	ok := roleSpecSvc.Get(&roleSpec, false, "", "")
	var parentSpecId uint64
	if ok {
		parentSpecId = roleSpec.SpecId
	}
	if parentSpecId <= 0 {
		return
	}

	var head []string
	var pos = -1
	attibuteSpecs := make([]interface{}, 0)
	connectionSpecs := make([]interface{}, 0)
	attributeSpecSvc := service.GetAttributeSpecService()
	connectionSpecSvc := service.GetConnectionSpecService()
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		if i == 0 {
			head = row
			for j := 0; j < len(row); j++ {
				if "SpecId" == xstrings.FirstRuneToUpper(row[j]) {
					pos = j
				}
			}
			continue
		}
		var specId uint64
		v, err := convert.ToObject(row[pos], "uint64")
		if err == nil {
			specId = v.(uint64)
		}
		attibuteSpec := entity.AttributeSpec{}
		attibuteSpec.SpecId = specId
		attibuteSpec.Status = baseentity.EntityStatus_Effective
		attributeSpecSvc.Get(&attibuteSpec, false, "", "")
		for j := 0; j < len(row); j++ {
			fieldname := xstrings.FirstRuneToUpper(head[j])
			if fieldname != "" {
				reflect.Set(&attibuteSpec, fieldname, row[j])
			}
		}
		if attibuteSpec.DataType == "BigDecimal" {
			attibuteSpec.DataType = "float64"
		} else if attibuteSpec.DataType == "String" {
			attibuteSpec.DataType = "string"
		} else if attibuteSpec.DataType == "Date" {
			attibuteSpec.DataType = "time.Time"
		} else if attibuteSpec.DataType == "Boolean" {
			attibuteSpec.DataType = "bool"
		} else if attibuteSpec.DataType == "Integer" {
			attibuteSpec.DataType = "int"
		} else if attibuteSpec.DataType == "Long" {
			attibuteSpec.DataType = "int64"
		}
		attibuteSpecs = append(attibuteSpecs, &attibuteSpec)

		connectionSpec := entity.ConnectionSpec{}
		connectionSpec.ParentSpecId = parentSpecId
		connectionSpec.SubSpecId = specId
		connectionSpec.Status = baseentity.EntityStatus_Effective
		connectionSpecSvc.Get(&connectionSpec, false, "", "")
		for j := 0; j < len(row); j++ {
			fieldname := xstrings.FirstRuneToUpper(head[j])
			if fieldname != "" {
				reflect.Set(&connectionSpec, fieldname, row[j])
			}
		}
		connectionSpec.Maxmium = 1
		connectionSpec.Minmium = 1
		connectionSpec.BuildNum = 1
		connectionSpec.SpecType = entity.SpecType_Property
		connectionSpec.RelationType = entity.RelationType_Dependency
		connectionSpecs = append(connectionSpecs, &connectionSpec)
	}
	attributeSpecSvc.Upsert(attibuteSpecs...)
	connectionSpecSvc.Upsert(connectionSpecs...)
}

func init() {

}
