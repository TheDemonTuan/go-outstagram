package entity

import (
	"github.com/google/uuid"
	"time"
)

type ReportReason int

const (
	Spam ReportReason = iota + 1
	JustDontLikeIt
	SuicideSelfInjuryOrEatingDisorders
	IllegalOrRegulatedGoods
	NudityOrSexualActivity
	HateSpeechOrSymbols
	ViolenceOrDangerousOrganisations
	BullyingOrHarassment
	IntellectualPropertyViolation
	ScamOrFraud
	FalseInformation

	Drugs
	SuicideOrSelfInjury
	EatingDisorders

	PretendingMe
	PretendingSomeFriend
	PretendingCelebrity
	PretendingBusiness

	UnderThe13Age
)

type ReportType int

const (
	ReportPost ReportType = iota + 1
	ReportReel
	ReportComment
	ReportUser
)

type ReportStatus int

const (
	Pending ReportStatus = iota
	Viewed
)

type Report struct {
	ID uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`

	ByUserID uuid.UUID    `json:"by_user_id" gorm:"not null;type:uuid;index"`
	Reason   ReportReason `json:"reason" gorm:"not null"`
	Type     ReportType   `json:"type" gorm:"not null"`
	Info     string       `json:"info" gorm:"not null"`
	Status   ReportStatus `json:"status" gorm:"not null,default:0"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
