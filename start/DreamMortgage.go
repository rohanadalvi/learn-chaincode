/*
Dream Mortgage Chaincode
*/

package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
)


//==============================================================================================================================
//	 Participating Entities
//==============================================================================================================================
/*const   FEDERAL_RESERVE   =  "federal_reserve"
const   CUSTOMER          =  "customer"
const   LENDING_BANK      =  "lendor"
const   PARTNER_BANK      =  "partner_bank"
const   AUDITOR           =  "auditor"
const   GSE               =  "gse"
const   BROKER            =  "broker"
const   CITY_COUNCIL      =   "city_council"
const   DATA_PROVIDER    =   "data_service_provider"*/


//==============================================================================================================================
//	 MORTGAGE STAGES/LIFE CYCLE.
//==============================================================================================================================
/*const   APPLICATION  			      =  0
const   LENDING_DECISION  			=  1
const   APPROVED_DENIED  			  =  2
const   RESELL                	=  3
const   SOLD            			  =  4*/

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//==============================================================================================================================
//	Mortgage - Defines the structure for a Mortgage object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON customerName -> Struct customer Name.
//==============================================================================================================================
type Mortgage struct {
	CustomerName               string  `json:"CustomerName"`
	CustomerAddress            string  `json:"CustomerAddress"`
	CustomerSSN                int     `json:"CustomerSSN"`
	CustomerDOB                string  `json:"CustomerDOB"`
	MortgageNumber             int     `json:"MortgageNumber"`
	MortgageStage              string  `json:"MortgageStage"`
	MortgagePropertyOwnership  string  `json:"MortgagePropertyOwnership"`
	MortagePropertyAddress     string  `json:"MortagePropertyAddress"`
	ReqLoanAmount              int     `json:"ReqLoanAmount"`
	GrantedLoanAmount          int     `json:"GrantedLoanAmount"`
	MortgageType               string  `json:"MortgageType"`
	RateofInterest             float64 `json:"RateofInterest"`
	MortgageStartDate          string  `json:"MortgageStartDate"`
	MortgageDuration           int     `json:"MortgageDuration"`
	LastPaymentAmount          int     `json:"LastPaymentAmount"`
  PropertyValuation          int     `json:"PropertyValuation"`
	CreditScore                int     `json:"CreditScore"`
	FinancialWorth             int     `json:"FinancialWorth"`
	RiskClassification         string  `json:"RiskClassification"`
	RiskAdjustedReturn         float64 `json:"RiskAdjustedReturn"`
	ExpectedAnnualCashflow     int     `json:"ExpectedAnnualCashflow"`
	RemainingMortgageAmount    int     `json:"RemainingMortgageAmount"`
	Ownershipcost              int     `json:"Ownershipcost"`
	ConformedMortgage          bool    `json:"ConformedMortgage"`
	ModifiedBy                 string  `json:"ModifiedBy"`
}

//==============================================================================================================================
//	Mortgage Portfolio - Defines the structure that holds all the Mortgage
//				Used as an index when querying all Mortgage.
//==============================================================================================================================

type mortgage_portfolio struct {
	MortgageNumbers []int    `json:"MortgageNumbers"`
	CustomerNames   []string `json:"CustomerNames"`
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
	mortgages.MortgageNumbers = []int{1000000}
	mortgages.CustomerNames   = []string{""}
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
  } else if function == "modify_Mortgage" {
     return t.modify_Mortgage(stub, args)
  }

	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// write function
func (t *SimpleChaincode) create_Mortgage_application(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    // Variable declaration
	  var mortgage Mortgage
		var mortgages mortgage_portfolio
    var err error
		var bytes []byte

		//Logging
    fmt.Println("running create_Mortgage_application()")

    // verify is the Json is sent.
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting one JSON object to create mortgage application")
    }
		//Assign JSON input and convert it to bytes
		mortgage_json := args[0]
    err = json.Unmarshal([]byte(mortgage_json), &mortgage)
    if err != nil {
			  return nil, errors.New("error while Unmarshalling mortgage json object")
		}

		//Get latest mortgages porfolio in blockchain and assign it to variable array
		bytes, err = stub.GetState("mortgages")
		if err != nil {
			  return nil, errors.New("error while retrieving mortgage portfolio json object")
		}
		err = json.Unmarshal(bytes,&mortgages)
		if err != nil {
 			  return nil, errors.New("error while Unmarshalling mortgages for new mortgage number")
 		}

		// Generate Unique mortgage number and append to Mortgage portfolio
		mortgage.MortgageNumber = mortgages.MortgageNumbers[len(mortgages.MortgageNumbers)-1]+1
	  mortgages.MortgageNumbers = append(mortgages.MortgageNumbers,mortgage.MortgageNumber)
	  mortgages.CustomerNames   = append(mortgages.CustomerNames,mortgage.CustomerName)

    //Store Mortgage in blockchain
		mortgagebytes, err := json.Marshal(mortgage)
		if err != nil {
			 return nil, errors.New("Error in Marshalling New Mortgage record")
		 }
		err = stub.PutState(string(mortgage.MortgageNumber),mortgagebytes)
    if err != nil {
        return nil, err
    }

    //Store current Mortgage Portfolio in blockchain.
		bytes, err = json.Marshal(mortgages)
		if err != nil {
			 return nil, errors.New("Add to Mortgage Portfolio record")
		 }
		err = stub.PutState("mortgages", bytes)

    return nil, nil
}

func (t *SimpleChaincode) modify_Mortgage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    // Variable declaration
	  var mortgage Mortgage
		var currentmortgage Mortgage
    var err error
		var mortgagebytes []byte

		//Logging
    fmt.Println("running modify_Mortgage()")

    // verify is the Json is sent.
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting one JSON object to create mortgage application")
    }
		//Assign JSON input and convert it to bytes
		mortgage_json := args[0]
    err = json.Unmarshal([]byte(mortgage_json), &mortgage)
    if err != nil {
			  return nil, errors.New("error while Unmarshalling mortgage json object")
		}

		//Get latest mortgages porfolio in blockchain and assign it to variable array
		mortgagebytes, err = stub.GetState(string(mortgage.MortgageNumber))
		if err != nil {
 			  return nil, errors.New("error while fetching mortgage number")
 		}

		err = json.Unmarshal(mortgagebytes,&currentmortgage)
		if err != nil {
 			  return nil, errors.New("error while Unmarshalling mortgages for current mortgage number")
 		}

		//Update current Mortgage Fields
		err = json.Unmarshal([]byte(mortgage_json), &currentmortgage)
    if err != nil {
			  return nil, errors.New("error while Unmarshalling mortgage json object")
		}


    //Store Mortgage in blockchain
		mortgagebytes, err = json.Marshal(currentmortgage)
		if err != nil {
			 return nil, errors.New("Error in Marshalling New Mortgage record")
		 }
		err = stub.PutState(string(currentmortgage.MortgageNumber), mortgagebytes)
    if err != nil {
        return nil, err
    }

    return nil, nil
}


// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "retrieve_mortgage_portfolio" {                            //read a variable
           return t.retrieve_mortgage_portfolio(stub, args)
        }

	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) retrieve_mortgage_portfolio(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var jsonResp string
    var err error

    //retrieve Mortgage Portfolio
    valAsbytes, err := stub.GetState("mortgages")
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to retrieve mortgage portfolio\"}"
        return nil, errors.New(jsonResp)
    }
    return valAsbytes, nil
}
