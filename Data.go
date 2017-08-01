package main

type KycData struct {
	USER_NAME           string `json:"USER_NAME"`
	USER_ID             string `json:"USER_ID"`
	KYC_BANK_NAME       string `json:"KYC_BANK_NAME"`
	KYC_CREATE_DATE     string `json:"KYC_CREATE_DATE"`
	KYC_VALID_TILL_DATE string `json:"KYC_VALID_TILL_DATE"`
	KYC_STATUS          string `json:"KYC_STATUS"`
	//KYC_DOCUMENTS       KycDoc `json:"KYC_DOCUMENTS"`
}

type KycDoc struct {
	USER_ID       string `json:"USER_ID"`
	DOCUMENT_TYPE string `json:"DOCUMENT_TYPE"`
	DOCUMENT_BLOB string `json:"DOCUMENT_BLOB"`
}

/*type KycBankDetails struct {
	KYC_BANK_NAME string `json:"KYC_BANK_NAME"`
}*/

type KycCount struct {
	AllContracts      int
	ExpiringContracts int
	CreatedContracts  int
}
