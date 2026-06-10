package models

import (
	"time"

	"github.com/google/uuid"
)

type MovableOwnerType string

type StorageLocation string

type MovableCondition string

type MovableRecord struct {
	ID          uuid.UUID    `json:"id"`
	RecordID    string       `json:"record_id"`
	RegistrarID uuid.UUID    `json:"registrar_id"`
	Status      RecordStatus `json:"status"`

	NameAmharic string  `json:"name_amharic"`
	NameLocal   *string `json:"name_local,omitempty"`
	Category    *string `json:"category,omitempty"`

	LocationName *string `json:"location_name,omitempty"`
	Woreda       *string `json:"woreda,omitempty"`
	Kebele       *string `json:"kebele,omitempty"`
	HouseNumber  *string `json:"house_number,omitempty"`
	CurrentUse   *string `json:"current_use,omitempty"`
	PreviousID   *string `json:"previous_id,omitempty"`

	OwnerType            *MovableOwnerType `json:"owner_type,omitempty"`
	OwnerName            *string           `json:"owner_name,omitempty"`
	StorageLocation      *StorageLocation  `json:"storage_location,omitempty"`
	StorageLocationOther *string           `json:"storage_location_other,omitempty"`

	MadeBy             *string    `json:"made_by,omitempty"`
	PeriodMade         *string    `json:"period_made,omitempty"`
	AgeMethod          *AgeMethod `json:"age_method,omitempty"`
	AcquisitionMethods []string   `json:"acquisition_methods,omitempty"`

	HeightCM         *float64 `json:"height_cm,omitempty"`
	WidthCM          *float64 `json:"width_cm,omitempty"`
	LengthCM         *float64 `json:"length_cm,omitempty"`
	DiameterCM       *float64 `json:"diameter_cm,omitempty"`
	ThicknessCM      *float64 `json:"thickness_cm,omitempty"`
	WeightKG         *float64 `json:"weight_kg,omitempty"`
	NumPages         *int     `json:"num_pages,omitempty"`
	NumChapters      *int     `json:"num_chapters,omitempty"`
	NumIllustrations *int     `json:"num_illustrations,omitempty"`

	ColorType       *string  `json:"color_type,omitempty"`
	HasDecoration   *bool    `json:"has_decoration,omitempty"`
	Materials       []string `json:"materials,omitempty"`
	MaterialOther   *string  `json:"material_other,omitempty"`
	Description     *string  `json:"description,omitempty"`
	NotableBecause  []string `json:"notable_because,omitempty"`
	NotableOther    *string  `json:"notable_other,omitempty"`
	Significance    *string  `json:"significance,omitempty"`

	Condition          *MovableCondition   `json:"condition,omitempty"`
	HasThreat          *bool               `json:"has_threat,omitempty"`
	ThreatDescription  *string             `json:"threat_description,omitempty"`
	MaintenanceDone    *bool               `json:"maintenance_done,omitempty"`
	MaintenanceBy      *string             `json:"maintenance_by,omitempty"`
	MaintenanceDate    *time.Time          `json:"maintenance_date,omitempty"`
	MaintenanceCount   *int                `json:"maintenance_count,omitempty"`
	PreventiveLevel    *QualityLevel       `json:"preventive_level,omitempty"`
	Accessibility      *AccessibilityLevel `json:"accessibility,omitempty"`
	Notes              *string             `json:"notes,omitempty"`

	RelatedDocs []string `json:"related_docs,omitempty"`

	InformantName       *string  `json:"informant_name,omitempty"`
	InformantSex        *SexType `json:"informant_sex,omitempty"`
	InformantAge        *int     `json:"informant_age,omitempty"`
	InformantOccupation *string  `json:"informant_occupation,omitempty"`
	CaretakerName       *string  `json:"caretaker_name,omitempty"`
	CaretakerRole       *string  `json:"caretaker_role,omitempty"`
	RegistrarDate       *time.Time `json:"registrar_date,omitempty"`

	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	ApprovedBy *uuid.UUID `json:"approved_by,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type MovableRecordInput struct {
	NameAmharic *string `json:"name_amharic"`
	NameLocal   *string `json:"name_local"`
	Category    *string `json:"category"`

	LocationName *string `json:"location_name"`
	Woreda       *string `json:"woreda"`
	Kebele       *string `json:"kebele"`
	HouseNumber  *string `json:"house_number"`
	CurrentUse   *string `json:"current_use"`
	PreviousID   *string `json:"previous_id"`

	OwnerType            *MovableOwnerType `json:"owner_type"`
	OwnerName            *string           `json:"owner_name"`
	StorageLocation      *StorageLocation  `json:"storage_location"`
	StorageLocationOther *string           `json:"storage_location_other"`

	MadeBy             *string    `json:"made_by"`
	PeriodMade         *string    `json:"period_made"`
	AgeMethod          *AgeMethod `json:"age_method"`
	AcquisitionMethods []string   `json:"acquisition_methods"`

	HeightCM         *float64 `json:"height_cm"`
	WidthCM          *float64 `json:"width_cm"`
	LengthCM         *float64 `json:"length_cm"`
	DiameterCM       *float64 `json:"diameter_cm"`
	ThicknessCM      *float64 `json:"thickness_cm"`
	WeightKG         *float64 `json:"weight_kg"`
	NumPages         *int     `json:"num_pages"`
	NumChapters      *int     `json:"num_chapters"`
	NumIllustrations *int     `json:"num_illustrations"`

	ColorType      *string  `json:"color_type"`
	HasDecoration  *bool    `json:"has_decoration"`
	Materials      []string `json:"materials"`
	MaterialOther  *string  `json:"material_other"`
	Description    *string  `json:"description"`
	NotableBecause []string `json:"notable_because"`
	NotableOther   *string  `json:"notable_other"`
	Significance   *string  `json:"significance"`

	Condition         *MovableCondition   `json:"condition"`
	HasThreat         *bool               `json:"has_threat"`
	ThreatDescription *string             `json:"threat_description"`
	MaintenanceDone   *bool               `json:"maintenance_done"`
	MaintenanceBy     *string             `json:"maintenance_by"`
	MaintenanceDate   *time.Time          `json:"maintenance_date"`
	MaintenanceCount  *int                `json:"maintenance_count"`
	PreventiveLevel   *QualityLevel       `json:"preventive_level"`
	Accessibility     *AccessibilityLevel `json:"accessibility"`
	Notes             *string             `json:"notes"`

	RelatedDocs []string `json:"related_docs"`

	InformantName       *string  `json:"informant_name"`
	InformantSex        *SexType `json:"informant_sex"`
	InformantAge        *int     `json:"informant_age"`
	InformantOccupation *string  `json:"informant_occupation"`
	CaretakerName       *string  `json:"caretaker_name"`
	CaretakerRole       *string  `json:"caretaker_role"`
	RegistrarDate       *time.Time `json:"registrar_date"`
}

func ApplyMovableInput(record *MovableRecord, input MovableRecordInput) {
	if input.NameAmharic != nil {
		record.NameAmharic = *input.NameAmharic
	}
	if input.NameLocal != nil {
		record.NameLocal = input.NameLocal
	}
	if input.Category != nil {
		record.Category = input.Category
	}
	if input.LocationName != nil {
		record.LocationName = input.LocationName
	}
	if input.Woreda != nil {
		record.Woreda = input.Woreda
	}
	if input.Kebele != nil {
		record.Kebele = input.Kebele
	}
	if input.HouseNumber != nil {
		record.HouseNumber = input.HouseNumber
	}
	if input.CurrentUse != nil {
		record.CurrentUse = input.CurrentUse
	}
	if input.PreviousID != nil {
		record.PreviousID = input.PreviousID
	}
	if input.OwnerType != nil {
		record.OwnerType = input.OwnerType
	}
	if input.OwnerName != nil {
		record.OwnerName = input.OwnerName
	}
	if input.StorageLocation != nil {
		record.StorageLocation = input.StorageLocation
	}
	if input.StorageLocationOther != nil {
		record.StorageLocationOther = input.StorageLocationOther
	}
	if input.MadeBy != nil {
		record.MadeBy = input.MadeBy
	}
	if input.PeriodMade != nil {
		record.PeriodMade = input.PeriodMade
	}
	if input.AgeMethod != nil {
		record.AgeMethod = input.AgeMethod
	}
	if input.AcquisitionMethods != nil {
		record.AcquisitionMethods = input.AcquisitionMethods
	}
	if input.HeightCM != nil {
		record.HeightCM = input.HeightCM
	}
	if input.WidthCM != nil {
		record.WidthCM = input.WidthCM
	}
	if input.LengthCM != nil {
		record.LengthCM = input.LengthCM
	}
	if input.DiameterCM != nil {
		record.DiameterCM = input.DiameterCM
	}
	if input.ThicknessCM != nil {
		record.ThicknessCM = input.ThicknessCM
	}
	if input.WeightKG != nil {
		record.WeightKG = input.WeightKG
	}
	if input.NumPages != nil {
		record.NumPages = input.NumPages
	}
	if input.NumChapters != nil {
		record.NumChapters = input.NumChapters
	}
	if input.NumIllustrations != nil {
		record.NumIllustrations = input.NumIllustrations
	}
	if input.ColorType != nil {
		record.ColorType = input.ColorType
	}
	if input.HasDecoration != nil {
		record.HasDecoration = input.HasDecoration
	}
	if input.Materials != nil {
		record.Materials = input.Materials
	}
	if input.MaterialOther != nil {
		record.MaterialOther = input.MaterialOther
	}
	if input.Description != nil {
		record.Description = input.Description
	}
	if input.NotableBecause != nil {
		record.NotableBecause = input.NotableBecause
	}
	if input.NotableOther != nil {
		record.NotableOther = input.NotableOther
	}
	if input.Significance != nil {
		record.Significance = input.Significance
	}
	if input.Condition != nil {
		record.Condition = input.Condition
	}
	if input.HasThreat != nil {
		record.HasThreat = input.HasThreat
	}
	if input.ThreatDescription != nil {
		record.ThreatDescription = input.ThreatDescription
	}
	if input.MaintenanceDone != nil {
		record.MaintenanceDone = input.MaintenanceDone
	}
	if input.MaintenanceBy != nil {
		record.MaintenanceBy = input.MaintenanceBy
	}
	if input.MaintenanceDate != nil {
		record.MaintenanceDate = input.MaintenanceDate
	}
	if input.MaintenanceCount != nil {
		record.MaintenanceCount = input.MaintenanceCount
	}
	if input.PreventiveLevel != nil {
		record.PreventiveLevel = input.PreventiveLevel
	}
	if input.Accessibility != nil {
		record.Accessibility = input.Accessibility
	}
	if input.Notes != nil {
		record.Notes = input.Notes
	}
	if input.RelatedDocs != nil {
		record.RelatedDocs = input.RelatedDocs
	}
	if input.InformantName != nil {
		record.InformantName = input.InformantName
	}
	if input.InformantSex != nil {
		record.InformantSex = input.InformantSex
	}
	if input.InformantAge != nil {
		record.InformantAge = input.InformantAge
	}
	if input.InformantOccupation != nil {
		record.InformantOccupation = input.InformantOccupation
	}
	if input.CaretakerName != nil {
		record.CaretakerName = input.CaretakerName
	}
	if input.CaretakerRole != nil {
		record.CaretakerRole = input.CaretakerRole
	}
	if input.RegistrarDate != nil {
		record.RegistrarDate = input.RegistrarDate
	}
}

func NewDraftMovableRecord(registrarID uuid.UUID) MovableRecord {
	return MovableRecord{
		RegistrarID: registrarID,
		Status:      StatusDraft,
		NameAmharic: "",
	}
}
