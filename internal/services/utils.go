package services

import (
	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

func ToProtoContextType(internal sharedpb.ContextType) sharedpb.ContextType {
	switch internal {
	case sharedpb.ContextType_Colleges:
		return sharedpb.ContextType_Colleges
	case sharedpb.ContextType_Certifications:
		return sharedpb.ContextType_Certifications
	case sharedpb.ContextType_EntranceExams:
		return sharedpb.ContextType_EntranceExams
	case sharedpb.ContextType_APExams:
		return sharedpb.ContextType_APExams
	case sharedpb.ContextType_DoD:
		return sharedpb.ContextType_DoD
	case sharedpb.ContextType_UserGeneratedContent:
		return sharedpb.ContextType_UserGeneratedContent
	case sharedpb.ContextType_All:
		return sharedpb.ContextType_All
	default:
		return sharedpb.ContextType_All
	}
}

func FromProtoContextType(proto sharedpb.ContextType) sharedpb.ContextType {
	switch proto {
	case sharedpb.ContextType_Colleges:
		return sharedpb.ContextType_Colleges
	case sharedpb.ContextType_Certifications:
		return sharedpb.ContextType_Certifications
	case sharedpb.ContextType_EntranceExams:
		return sharedpb.ContextType_EntranceExams
	case sharedpb.ContextType_APExams:
		return sharedpb.ContextType_APExams
	case sharedpb.ContextType_DoD:
		return sharedpb.ContextType_DoD
	case sharedpb.ContextType_UserGeneratedContent:
		return sharedpb.ContextType_UserGeneratedContent
	case sharedpb.ContextType_All:
		return sharedpb.ContextType_All
	default:
		return sharedpb.ContextType_All
	}
}
