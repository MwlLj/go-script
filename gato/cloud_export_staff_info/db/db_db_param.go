package db

type CGetAllStaffInfoOutput struct {
	InfrastructureUuid string
	InfrastructureUuidIsValid bool
	StaffName string
	StaffNameIsValid bool
	Address string
	AddressIsValid bool
	Nation string
	NationIsValid bool
	Nationality string
	NationalityIsValid bool
	NativePlace string
	NativePlaceIsValid bool
	Gender string
	GenderIsValid bool
	ContactType string
	ContactTypeIsValid bool
	ContactInfo string
	ContactInfoIsValid bool
}

type CGetInfrastructureInfoByUuidInput struct {
	InfrastructureUuid string
	InfrastructureUuidIsValid bool
}

type CGetInfrastructureInfoByUuidOutput struct {
	InfrastructureName string
	InfrastructureNameIsValid bool
	ParentUuid string
	ParentUuidIsValid bool
}

