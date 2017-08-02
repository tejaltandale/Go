package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func CreateDatabase(stub shim.ChaincodeStubInterface) error {
	var err error

	//Create table "KycDetails"
	err = stub.CreateTable("KycDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "USER_ID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "KYC_BANK_NAME", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "USER_NAME", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "KYC_CREATE_DATE", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "KYC_VALID_TILL_DATE", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "KYC_STATUS", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return errors.New("Failed creating KycDetails table.")
	}

	//Create table "KycDocDetails"
	err = stub.CreateTable("KycDocDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "USER_ID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "DOCUMENT_TYPE", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "DOCUMENT_BLOB", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return errors.New("Failed creating KycDocDetails table.")
	}

	//Create table "BankDetails"
	err = stub.CreateTable("BankDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "BankName", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "UserList", Type: shim.ColumnDefinition_BYTES, Key: false},
	})
	if err != nil {
		return errors.New("Failed creating BankDetails table.")
	}

	//Create Bank List
	var BankList []string
	jsonAsBytes, _ := json.Marshal(BankList)
	err = stub.PutState("BankList", jsonAsBytes)
	if err != nil {
		return errors.New("Failed to put Bank List")
	}
	return nil
}

func InsertKYCDetails(stub shim.ChaincodeStubInterface, Kycdetails KycData) (bool, error) {
	return stub.InsertRow("KycDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: Kycdetails.USER_ID}},
			&shim.Column{Value: &shim.Column_String_{String_: Kycdetails.KYC_BANK_NAME}},
			&shim.Column{Value: &shim.Column_String_{String_: Kycdetails.USER_NAME}},
			&shim.Column{Value: &shim.Column_String_{String_: Kycdetails.KYC_CREATE_DATE}},
			&shim.Column{Value: &shim.Column_String_{String_: Kycdetails.KYC_VALID_TILL_DATE}},
			&shim.Column{Value: &shim.Column_String_{String_: Kycdetails.KYC_STATUS}},
		},
	})
}

func InsertKYCDocumentDetails(stub shim.ChaincodeStubInterface, KycDocDetails KycDoc) (bool, error) {
	return stub.InsertRow("KycDocDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: KycDocDetails.USER_ID}},
			&shim.Column{Value: &shim.Column_String_{String_: KycDocDetails.DOCUMENT_TYPE}},
			&shim.Column{Value: &shim.Column_String_{String_: KycDocDetails.DOCUMENT_BLOB}},
		},
	})
}

func InsertBankDetails(stub shim.ChaincodeStubInterface, BankName string, UserList []string) (bool, error) {
	var ok bool
	var err error
	var BankList []string

	JsonAsBytes, _ := json.Marshal(UserList)

	ok, err = stub.InsertRow("BankDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: BankName}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: JsonAsBytes}},
		},
	})
	if !ok && err == nil {
		return false, errors.New("Error in adding BankDetails record.")
	}

	BankList, err = GetBankList(stub)
	if err != nil {
		return false, err
	}

	//Update Bank List
	BankList = append(BankList, BankName)

	ok, err = PutBankList(stub, BankList)
	if !ok {
		return false, err
	}

	return true, nil

}

func GetBankList(stub shim.ChaincodeStubInterface) ([]string, error) {
	// Get bank List
	var BankList []string
	jsonAsBytes, err := stub.GetState("BankList")
	if err != nil {
		return nil, errors.New("Failed to get Bank List")
	}
	json.Unmarshal(jsonAsBytes, &BankList)
	return BankList, nil
}

func PutBankList(stub shim.ChaincodeStubInterface, BankList []string) (bool, error) {
	//Put Bank List
	jsonAsBytes, _ := json.Marshal(BankList)
	err := stub.PutState("BankList", jsonAsBytes)
	if err != nil {
		return false, errors.New("Failed to put Bank List")
	}
	return true, nil
}

func GetKYCDetails(stub shim.ChaincodeStubInterface, UserId string) (KycData, error) {
	var KycDataObj KycData
	var columns []shim.Column
	var err error

	col1 := shim.Column{Value: &shim.Column_String_{String_: UserId}}
	columns = append(columns, col1)

	row, err := stub.GetRow("KycDetails", columns)
	if err != nil {
		return KycDataObj, errors.New("Failed to query")
	}

	KycDataObj.USER_ID = row.Columns[0].GetString_()
	KycDataObj.KYC_BANK_NAME = row.Columns[1].GetString_()
	KycDataObj.USER_NAME = row.Columns[2].GetString_()
	KycDataObj.KYC_CREATE_DATE = row.Columns[3].GetString_()
	KycDataObj.KYC_VALID_TILL_DATE = row.Columns[4].GetString_()

	return KycDataObj, nil
}

func GetBankSpecificKYCDetails(stub shim.ChaincodeStubInterface, UserId string, BankName string) (KycData, error) {
	var KycDataObj KycData
	var columns []shim.Column
	var err error

	col1 := shim.Column{Value: &shim.Column_String_{String_: UserId}}
	columns = append(columns, col1)

	row, err := stub.GetRow("KycDetails", columns)
	if err != nil {
		return KycDataObj, errors.New("Failed to query")
	}

	if row.Columns == nil {
		return KycDataObj, nil
	}

	KycDataObj.USER_ID = row.Columns[0].GetString_()
	KycDataObj.KYC_BANK_NAME = row.Columns[1].GetString_()
	KycDataObj.USER_NAME = row.Columns[2].GetString_()
	KycDataObj.KYC_CREATE_DATE = row.Columns[3].GetString_()
	KycDataObj.KYC_VALID_TILL_DATE = row.Columns[4].GetString_()

	return KycDataObj, nil
}

func GetUserList(stub shim.ChaincodeStubInterface, BankName string) ([]string, error) {
	var UserList []string
	var columns []shim.Column

	col1 := shim.Column{Value: &shim.Column_String_{String_: BankName}}
	columns = append(columns, col1)

	row, err := stub.GetRow("BankDetails", columns)
	if err != nil {
		return UserList, errors.New("Failed to query table BankDetails")
	}

	UsersAsBytes := row.Columns[1].GetBytes()
	json.Unmarshal(UsersAsBytes, &UserList)

	return UserList, nil
}

func GetDocument(stub shim.ChaincodeStubInterface, UserId string, DocumentType string) (string, error) {
	var columns []shim.Column
	var err error

	col1 := shim.Column{Value: &shim.Column_String_{String_: UserId}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: DocumentType}}
	columns = append(columns, col1)
	columns = append(columns, col2)

	row, err := stub.GetRow("KycDocDetails", columns)
	if err != nil {
		return "", errors.New("Failed to query")
	}
	return row.Columns[2].GetString_(), nil
}

func UpdateBankDetails(stub shim.ChaincodeStubInterface, BankName string, Userlist []string) (bool, error) {

	JsonAsBytes, _ := json.Marshal(Userlist)

	return stub.ReplaceRow("BankDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: BankName}},
			&shim.Column{Value: &shim.Column_Bytes{Bytes: JsonAsBytes}},
		},
	})
}

func UpdateKycDetails(stub shim.ChaincodeStubInterface, KycDetails KycData) (bool, error) {

	return stub.ReplaceRow("KycDetails", shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: KycDetails.USER_ID}},
			&shim.Column{Value: &shim.Column_String_{String_: KycDetails.KYC_BANK_NAME}},
			&shim.Column{Value: &shim.Column_String_{String_: KycDetails.USER_NAME}},
			&shim.Column{Value: &shim.Column_String_{String_: KycDetails.KYC_CREATE_DATE}},
			&shim.Column{Value: &shim.Column_String_{String_: KycDetails.KYC_VALID_TILL_DATE}},
		},
	})
}

