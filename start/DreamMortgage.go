/*
Dream Mortgage Chaincode
*/

package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"regexp"
)


//==============================================================================================================================
//	 Participating Entities
//==============================================================================================================================
const   FEDERAL_RESERVE   =  "federal_reserve"
const   CUSTOMER          =  "customer"
const   LENDING_BANK      =  "lendor"
const   PARTNER_BANK      =  "partner_bank"
const   AUDITOR           =  "auditor"
const   GSE               =  "gse"
const   BROKER            =  "broker"
const   CITY_COUNCIL      =   "city_council"
const   DATA_PROVIDER    =   "data_service_provider"


//==============================================================================================================================
//	 MORTGAGE STAGES/LIFE CYCLE.
//==============================================================================================================================
const   APPLICATION  			      =  0
const   LENDING_DECISION  			=  1
const   APPROVED_DENIED  			  =  2
const   RESELL                	=  3
const   SOLD            			  =  4

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//==============================================================================================================================
//	Mortgage - Defines the structure for a Mortgage object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON customerName -> Struct customer Name.
//==============================================================================================================================
type Mortgage struct {
	customerName               string  `json:"customerName"`
	customerAddress            string  `json:"customerAddress"`
	customerSSN                string  `json:"customerSSN"`
	customerDOB                string  `json:"customerDOB"`
	mortgageNumber             int     `json:"mortgageNumber"`
	mortgageStage              string  `json:"mortgageStage"`
	mortgagePropertyOwnership  string  `json:"MortgagePropertyOwnership"`
	mortagePropertyAddress     string  `json:"mortagePropertyAddress"`
	reqLoanAmount              int     `json:"reqLoanAmount"`
	grantedLoanAmount          int     `json:"grantedLoanAmount"`
	mortgageType               string  `json:"mortgageType"`
	rateofInterest             float64 `json:"mortgageType"`
	mortgageStartDate          string  `json:"mortgageStartDate"`
	mortgageDuration           int     `json:"mortgageDuration"`
	lastPaymentAmount          int     `json:"lastPaymentAmount"`
  propertyValuation          int     `json:"propertyValuation"`
	creditScore                int     `json:"propertyValuation"`
	financialWorth             int     `json:"propertyValuation"`
	riskClassification         string  `json:"riskClassification"`
	riskAdjustedReturn         float64  `json:"riskAdjustedReturn"`
	expectedAnnualCashflow     int     `json:"expectedAnnualCashflow"`
	remainingMortgageAmount    int     `json:"remainingMortgageAmount"`
	ownershipcost              int     `json:"ownershipcost"`
	conformedMortgage          bool    `json:"conformedMortgage"`
	modifiedBy                 string  `json:"modifiedBy"`
}

//==============================================================================================================================
//	Mortgage Portfolio - Defines the structure that holds all the Mortgage
//				Used as an index when querying all Mortgage.
//==============================================================================================================================

type mortgage_portfolio struct {
	mortgageNumbers []int    `json:"mortgageNumbers"`
	customerNames   []string `json:"customerNames"`
}
// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

  // initialize the Mortgage number to 1000000
	var mortgages mortgage_portfolio
	mortgages.mortgageNumbers = []int{1000000}
	mortgages.customerNames   = []string{""}
	bytes, err := json.Marshal(mortgages)
	if err != nil {
		 return nil, errors.New("Error creating Mortgage Portfolio record")
	 }

  // initialize the Mortgage Portfolio
	err = stub.PutState("mortgages", bytes)
	if err != nil {
		return nil, errors.New("Error storing Mortgage Portfolio record in blockchain")
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		 return t.Init(stub, "init", args)
	} else if function == "create_Mortgage_application" {
                return t.create_Mortgage_application(stub, args)
        }

	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// write function
func (t *SimpleChaincode) create_Mortgage_application(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var mortgage Mortgage
		var mortgages mortgage_portfolio
    var err error
		var bytes []byte
    fmt.Println("running create_Mortgage_application()")

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting one JSON object to create mortgage application")
    }
		mortgage_json := args[0]
    err = json.Unmarshal([]byte(mortgage_json), &mortgage)
    if err != nil {
			  return nil, error.New("error while Unmarshalling mortgage json object")
		}
		stub.GetState("mortgages", bytes)
		err = json.Unmarshal(bytes,&mortgages)
		if err != nil {
 			  return nil, error.New("error while Unmarshalling mortgages for new mortgage number")
 		}
		mortgage.mortgageNumber = mortgages.mortgageNumbers[len(mortgages.mortgageNumbers)-1]+1
	  mortgages.mortgageNumbers = append(mortgages.mortgageNumbers,mortgage.mortgageNumber)
	  mortgages.customerNames   = append(mortgages.customerNames,mortgage.customerName)

    err = stub.PutState(mortgage.mortgageNumber, []byte(mortgage))  //write the variable into the chaincode state
    if err != nil {
        return nil, err
    }
    return nil, nil
}


// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {                            //read a variable
           return t.read(stub, args)
        }

	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var name, jsonResp string
    var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
    }

    name = args[0]
    valAsbytes, err := stub.GetState(name)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
        return nil, errors.New(jsonResp)
    }

    return valAsbytes, nil
}
