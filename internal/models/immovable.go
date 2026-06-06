package models

import (
	"time"

	"github.com/google/uuid"
)

type ImmovableRecord struct {
	ID          uuid.UUID    `json:"id"`
	RecordID    string       `json:"record_id"`
	RegistrarID uuid.UUID    `json:"registrar_id"`
	Status      RecordStatus `json:"status"`

	NameAmharic     string   `json:"name_amharic"`
	NameLocal       *string  `json:"name_local,omitempty"`
	Category        []string `json:"category,omitempty"`
	CurrentUse      []string `json:"current_use,omitempty"`
	CurrentUseOther *string  `json:"current_use_other,omitempty"`
	PreviousID      *string  `json:"previous_id,omitempty"`

	Woreda        string  `json:"woreda"`
	Kebele        string  `json:"kebele"`
	HouseNumber   *string `json:"house_number,omitempty"`
	StreetNumber  *string `json:"street_number,omitempty"`
	Gate          *string `json:"gate,omitempty"`

	OwnerType    *ImmovableOwnerType `json:"owner_type,omitempty"`
	OwnerName    *string             `json:"owner_name,omitempty"`
	MapReference *string             `json:"map_reference,omitempty"`
	GPSEast      *float64            `json:"gps_east,omitempty"`
	GPSNorth     *float64            `json:"gps_north,omitempty"`
	ElevationM   *float64            `json:"elevation_m,omitempty"`

	BuiltBy              *string    `json:"built_by,omitempty"`
	ConstructionPeriod   *string    `json:"construction_period,omitempty"`
	AgeMethod            *AgeMethod `json:"age_method,omitempty"`
	HeightM              *float64   `json:"height_m,omitempty"`
	LengthM              *float64   `json:"length_m,omitempty"`
	WidthM               *float64   `json:"width_m,omitempty"`
	NumDoors             *int       `json:"num_doors,omitempty"`
	NumWindows           *int       `json:"num_windows,omitempty"`
	NumRooms             *int       `json:"num_rooms,omitempty"`
	Material             *string    `json:"material,omitempty"`
	Description          *string    `json:"description,omitempty"`
	HarariHouseGrades    []string   `json:"harari_house_grades,omitempty"`
	NeighborhoodType     *string    `json:"neighborhood_type,omitempty"`

	OverallCondition *OverallCondition `json:"overall_condition,omitempty"`
	DamageRoof       *DamageLevel      `json:"damage_roof,omitempty"`
	DamageCornice    *DamageLevel      `json:"damage_cornice,omitempty"`
	DamageWall       *DamageLevel      `json:"damage_wall,omitempty"`
	DamageFloor      *DamageLevel      `json:"damage_floor,omitempty"`
	DamageDoor       *DamageLevel      `json:"damage_door,omitempty"`
	DamageCupboard   *DamageLevel      `json:"damage_cupboard,omitempty"`
	DamageUpperFloor *DamageLevel      `json:"damage_upper_floor,omitempty"`
	DamageDera       *DamageLevel      `json:"damage_dera,omitempty"`
	DamagePillar     *DamageLevel      `json:"damage_pillar,omitempty"`

	ValueHistorical    *string `json:"value_historical,omitempty"`
	ValueCraftsmanship *string `json:"value_craftsmanship,omitempty"`
	ValueArtistic      *string `json:"value_artistic,omitempty"`
	ValueScientific    *string `json:"value_scientific,omitempty"`
	ValueCultural      *string `json:"value_cultural,omitempty"`

	HasThreat          *bool               `json:"has_threat,omitempty"`
	MaintenanceDone    *bool               `json:"maintenance_done,omitempty"`
	MaintenanceReason  *string             `json:"maintenance_reason,omitempty"`
	MaintenanceBy      *string             `json:"maintenance_by,omitempty"`
	MaintenanceDate    *time.Time          `json:"maintenance_date,omitempty"`
	MaintenanceCount   *int                `json:"maintenance_count,omitempty"`
	PreventiveLevel    *QualityLevel       `json:"preventive_level,omitempty"`
	Accessibility      *AccessibilityLevel `json:"accessibility,omitempty"`
	Notes              *string             `json:"notes,omitempty"`

	RelatedDocs     []string `json:"related_docs,omitempty"`
	HasOralHistory  *bool    `json:"has_oral_history,omitempty"`

	CaretakerName *string  `json:"caretaker_name,omitempty"`
	CaretakerRole *string  `json:"caretaker_role,omitempty"`
	InformantName *string  `json:"informant_name,omitempty"`
	InformantSex  *SexType `json:"informant_sex,omitempty"`
	InformantAge  *int     `json:"informant_age,omitempty"`
	RegistrarDate *time.Time `json:"registrar_date,omitempty"`

	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	ApprovedBy *uuid.UUID `json:"approved_by,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type ImmovableRecordInput struct {
	NameAmharic     *string  `json:"name_amharic"`
	NameLocal       *string  `json:"name_local"`
	Category        []string `json:"category"`
	CurrentUse      []string `json:"current_use"`
	CurrentUseOther *string  `json:"current_use_other"`
	PreviousID      *string  `json:"previous_id"`

	Woreda       *string `json:"woreda"`
	Kebele       *string `json:"kebele"`
	HouseNumber  *string `json:"house_number"`
	StreetNumber *string `json:"street_number"`
	Gate         *string `json:"gate"`

	OwnerType    *ImmovableOwnerType `json:"owner_type"`
	OwnerName    *string             `json:"owner_name"`
	MapReference *string             `json:"map_reference"`
	GPSEast      *float64            `json:"gps_east"`
	GPSNorth     *float64            `json:"gps_north"`
	ElevationM   *float64            `json:"elevation_m"`

	BuiltBy            *string    `json:"built_by"`
	ConstructionPeriod *string    `json:"construction_period"`
	AgeMethod          *AgeMethod `json:"age_method"`
	HeightM            *float64   `json:"height_m"`
	LengthM            *float64   `json:"length_m"`
	WidthM             *float64   `json:"width_m"`
	NumDoors           *int       `json:"num_doors"`
	NumWindows         *int       `json:"num_windows"`
	NumRooms           *int       `json:"num_rooms"`
	Material           *string    `json:"material"`
	Description        *string    `json:"description"`
	HarariHouseGrades  []string   `json:"harari_house_grades"`
	NeighborhoodType   *string    `json:"neighborhood_type"`

	OverallCondition *OverallCondition `json:"overall_condition"`
	DamageRoof       *DamageLevel      `json:"damage_roof"`
	DamageCornice    *DamageLevel      `json:"damage_cornice"`
	DamageWall       *DamageLevel      `json:"damage_wall"`
	DamageFloor      *DamageLevel      `json:"damage_floor"`
	DamageDoor       *DamageLevel      `json:"damage_door"`
	DamageCupboard   *DamageLevel      `json:"damage_cupboard"`
	DamageUpperFloor *DamageLevel      `json:"damage_upper_floor"`
	DamageDera       *DamageLevel      `json:"damage_dera"`
	DamagePillar     *DamageLevel      `json:"damage_pillar"`

	ValueHistorical    *string `json:"value_historical"`
	ValueCraftsmanship *string `json:"value_craftsmanship"`
	ValueArtistic      *string `json:"value_artistic"`
	ValueScientific    *string `json:"value_scientific"`
	ValueCultural      *string `json:"value_cultural"`

	HasThreat         *bool               `json:"has_threat"`
	MaintenanceDone   *bool               `json:"maintenance_done"`
	MaintenanceReason *string             `json:"maintenance_reason"`
	MaintenanceBy     *string             `json:"maintenance_by"`
	MaintenanceDate   *time.Time          `json:"maintenance_date"`
	MaintenanceCount  *int                `json:"maintenance_count"`
	PreventiveLevel   *QualityLevel       `json:"preventive_level"`
	Accessibility     *AccessibilityLevel `json:"accessibility"`
	Notes             *string             `json:"notes"`

	RelatedDocs    []string `json:"related_docs"`
	HasOralHistory *bool    `json:"has_oral_history"`

	CaretakerName *string  `json:"caretaker_name"`
	CaretakerRole *string  `json:"caretaker_role"`
	InformantName *string  `json:"informant_name"`
	InformantSex  *SexType `json:"informant_sex"`
	InformantAge  *int     `json:"informant_age"`
	RegistrarDate *time.Time `json:"registrar_date"`
}

func ApplyImmovableInput(record *ImmovableRecord, input ImmovableRecordInput) {
	if input.NameAmharic != nil {
		record.NameAmharic = *input.NameAmharic
	}
	if input.NameLocal != nil {
		record.NameLocal = input.NameLocal
	}
	if input.Category != nil {
		record.Category = input.Category
	}
	if input.CurrentUse != nil {
		record.CurrentUse = input.CurrentUse
	}
	if input.CurrentUseOther != nil {
		record.CurrentUseOther = input.CurrentUseOther
	}
	if input.PreviousID != nil {
		record.PreviousID = input.PreviousID
	}
	if input.Woreda != nil {
		record.Woreda = *input.Woreda
	}
	if input.Kebele != nil {
		record.Kebele = *input.Kebele
	}
	if input.HouseNumber != nil {
		record.HouseNumber = input.HouseNumber
	}
	if input.StreetNumber != nil {
		record.StreetNumber = input.StreetNumber
	}
	if input.Gate != nil {
		record.Gate = input.Gate
	}
	if input.OwnerType != nil {
		record.OwnerType = input.OwnerType
	}
	if input.OwnerName != nil {
		record.OwnerName = input.OwnerName
	}
	if input.MapReference != nil {
		record.MapReference = input.MapReference
	}
	if input.GPSEast != nil {
		record.GPSEast = input.GPSEast
	}
	if input.GPSNorth != nil {
		record.GPSNorth = input.GPSNorth
	}
	if input.ElevationM != nil {
		record.ElevationM = input.ElevationM
	}
	if input.BuiltBy != nil {
		record.BuiltBy = input.BuiltBy
	}
	if input.ConstructionPeriod != nil {
		record.ConstructionPeriod = input.ConstructionPeriod
	}
	if input.AgeMethod != nil {
		record.AgeMethod = input.AgeMethod
	}
	if input.HeightM != nil {
		record.HeightM = input.HeightM
	}
	if input.LengthM != nil {
		record.LengthM = input.LengthM
	}
	if input.WidthM != nil {
		record.WidthM = input.WidthM
	}
	if input.NumDoors != nil {
		record.NumDoors = input.NumDoors
	}
	if input.NumWindows != nil {
		record.NumWindows = input.NumWindows
	}
	if input.NumRooms != nil {
		record.NumRooms = input.NumRooms
	}
	if input.Material != nil {
		record.Material = input.Material
	}
	if input.Description != nil {
		record.Description = input.Description
	}
	if input.HarariHouseGrades != nil {
		record.HarariHouseGrades = input.HarariHouseGrades
	}
	if input.NeighborhoodType != nil {
		record.NeighborhoodType = input.NeighborhoodType
	}
	if input.OverallCondition != nil {
		record.OverallCondition = input.OverallCondition
	}
	if input.DamageRoof != nil {
		record.DamageRoof = input.DamageRoof
	}
	if input.DamageCornice != nil {
		record.DamageCornice = input.DamageCornice
	}
	if input.DamageWall != nil {
		record.DamageWall = input.DamageWall
	}
	if input.DamageFloor != nil {
		record.DamageFloor = input.DamageFloor
	}
	if input.DamageDoor != nil {
		record.DamageDoor = input.DamageDoor
	}
	if input.DamageCupboard != nil {
		record.DamageCupboard = input.DamageCupboard
	}
	if input.DamageUpperFloor != nil {
		record.DamageUpperFloor = input.DamageUpperFloor
	}
	if input.DamageDera != nil {
		record.DamageDera = input.DamageDera
	}
	if input.DamagePillar != nil {
		record.DamagePillar = input.DamagePillar
	}
	if input.ValueHistorical != nil {
		record.ValueHistorical = input.ValueHistorical
	}
	if input.ValueCraftsmanship != nil {
		record.ValueCraftsmanship = input.ValueCraftsmanship
	}
	if input.ValueArtistic != nil {
		record.ValueArtistic = input.ValueArtistic
	}
	if input.ValueScientific != nil {
		record.ValueScientific = input.ValueScientific
	}
	if input.ValueCultural != nil {
		record.ValueCultural = input.ValueCultural
	}
	if input.HasThreat != nil {
		record.HasThreat = input.HasThreat
	}
	if input.MaintenanceDone != nil {
		record.MaintenanceDone = input.MaintenanceDone
	}
	if input.MaintenanceReason != nil {
		record.MaintenanceReason = input.MaintenanceReason
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
	if input.HasOralHistory != nil {
		record.HasOralHistory = input.HasOralHistory
	}
	if input.CaretakerName != nil {
		record.CaretakerName = input.CaretakerName
	}
	if input.CaretakerRole != nil {
		record.CaretakerRole = input.CaretakerRole
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
	if input.RegistrarDate != nil {
		record.RegistrarDate = input.RegistrarDate
	}
}

func NewDraftImmovableRecord(registrarID uuid.UUID) ImmovableRecord {
	return ImmovableRecord{
		RegistrarID: registrarID,
		Status:      StatusDraft,
		NameAmharic: "",
		Woreda:      "",
		Kebele:      "",
	}
}
