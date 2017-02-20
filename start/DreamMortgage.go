/*
Dream Mortgage Chaincode
*/

package main

import (
	"errors"
	"fmt"
	"strings"
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
	MortgagePropertyAddress    string  `json:"MortgagePropertyAddress"`
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
	MortgageNumbers     			  []int    `json:"MortgageNumbers"`
	CustomerNames       			  []string `json:"CustomerNames"`
	MortgageStages      			  []string `json:"MortgageStages"`
	ConformedMortgages				  []bool   `json:"ConformedMortgages"`
	MortgagePropertyOwnerships  []string `json:"MortgagePropertyOwnerships"`
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

  // initialize the Mortgage number
	var mortgages mortgage_portfolio
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
	} else if function == "create_mortgage_application" {
     return t.create_mortgage_application(stub, args)
  } else if function == "modify_mortgage" {
     return t.modify_mortgage(stub, args)
  }

	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// write function
func (t *SimpleChaincode) create_mortgage_application(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    // Variable declaration
	  var mortgage Mortgage
		var mortgages mortgage_portfolio
    var err error
		var bytes []byte

		//Logging
    fmt.Println("running create_mortgage_application()")

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
		if len(mortgages.MortgageNumbers) > 0 {
			mortgage.MortgageNumber = mortgages.MortgageNumbers[len(mortgages.MortgageNumbers)-1]+1
		}else{
			mortgage.MortgageNumber =1000001
		}

    //setting default values.
		mortgage.MortgageStage="Pending-Bank:"
		mortgage.ConformedMortgage=false
		mortgage.MortgagePropertyOwnership="NOT_ACCQUIRED"

	  mortgages.MortgageNumbers             = append(mortgages.MortgageNumbers,mortgage.MortgageNumber)
	  mortgages.CustomerNames               = append(mortgages.CustomerNames,mortgage.CustomerName)
		mortgages.MortgageStages              = append(mortgages.MortgageStages,mortgage.MortgageStage)
		mortgages.ConformedMortgages          = append(mortgages.ConformedMortgages,mortgage.ConformedMortgage)
		mortgages.MortgagePropertyOwnerships  = append(mortgages.MortgagePropertyOwnerships,mortgage.MortgagePropertyOwnership)

    // Update Mortgage data into bytes.
		mortgagebytes, err := json.Marshal(mortgage)
		if err != nil {
			 return nil, errors.New("Error in Marshalling New Mortgage record")
		}

    //package updated Mortgage Portfolio data into bytes.
		bytes, err = json.Marshal(mortgages)
		if err != nil {
			 return nil, errors.New("Add to Mortgage Portfolio record")
		}

    //Store Mortgage in blockchain
		err = stub.PutState(string(mortgage.MortgageNumber),mortgagebytes)
	  if err != nil {
	      return nil, err
	  }

		//Store updated Mortgage Portfolio in blockchain
		err = stub.PutState("mortgages", bytes)

    return nil, nil
}

func (t *SimpleChaincode) modify_mortgage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    // Variable declaration
	  var mortgage Mortgage
		var currentmortgage Mortgage
		var mortgages mortgage_portfolio
    var err error
		var mortgagebytes, bytes []byte
		var amountDisbursed bool
		var counter, value, Ratio_1, Ratio_2, Ratio_3, Rating_Ratio int

		//Logging
    fmt.Println("running modify_mortgage()")

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

		//Get latest mortgage in blockchain and assign it to variable array
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

    // smart contract fields
		// Update Mortgage Stage.
		if strings.ToUpper(currentmortgage.MortgageStage)== "APPROVED:" && currentmortgage.Ownershipcost > 0 {
			 currentmortgage.GrantedLoanAmount = currentmortgage.Ownershipcost
			 currentmortgage.MortgageStage="Disbursed:"
			 amountDisbursed=true
		}else{
			 amountDisbursed=false
			 if (strings.ToUpper(currentmortgage.MortgageStage)== "DISBURSED:READY TO PURCHASE" || strings.ToUpper(currentmortgage.MortgageStage)== "DISBURSED:READY TO SELL") && mortgage.Ownershipcost > 0 {
				  currentmortgage.MortgageStage="Disbursed:Sold"
			 }
		}

    // Update Mortgage Property Ownership
     if amountDisbursed  {
			  currentmortgage.MortgagePropertyOwnership="LENDING_BANK"
		 }else if strings.ToUpper(currentmortgage.MortgageStage)== "DISBURSED:READY TO PURCHASE" && mortgage.Ownershipcost > 0 {
			  currentmortgage.MortgagePropertyOwnership="GSE"
		 }else if strings.ToUpper(currentmortgage.MortgageStage)== "DISBURSED:READY TO SELL" && mortgage.Ownershipcost > 0 {
			  currentmortgage.MortgagePropertyOwnership="PARTNER_BANK"
		 } else if strings.Index(strings.ToUpper(currentmortgage.MortgageStage),"DISBURSED:") < 0 {
			 currentmortgage.MortgagePropertyOwnership="NOT_ACCQUIRED"
		 }

		//Calculate RemainingMortgageAmount
		if strings.Index(strings.ToUpper(currentmortgage.MortgageStage),"DISBURSED:") < 0 {
			  if currentmortgage.GrantedLoanAmount > 0 {
			     currentmortgage.RemainingMortgageAmount = currentmortgage.GrantedLoanAmount
				 }else{
					 currentmortgage.RemainingMortgageAmount  = currentmortgage.ReqLoanAmount
				 }
		} else if amountDisbursed {
			  currentmortgage.RemainingMortgageAmount = currentmortgage.GrantedLoanAmount
		} else if (currentmortgage.RemainingMortgageAmount - currentmortgage.LastPaymentAmount) > 0 {
		    currentmortgage.RemainingMortgageAmount = currentmortgage.RemainingMortgageAmount - currentmortgage.LastPaymentAmount
		} else {
			  currentmortgage.RemainingMortgageAmount=0
		}

    // if customer pays out property is moved back to customer.
		if currentmortgage.RemainingMortgageAmount <=0 {
			 currentmortgage.MortgagePropertyOwnership="CUSTOMER"
		}

		// Calculate Risk Classification.
			if currentmortgage.RemainingMortgageAmount > 0 && currentmortgage.FinancialWorth > 0 && currentmortgage.CreditScore > 0 && currentmortgage.PropertyValuation > 0 {
			   switch {
			   case  currentmortgage.PropertyValuation*100/currentmortgage.RemainingMortgageAmount > 75:
				Ratio_1 = 100
			   case  currentmortgage.PropertyValuation*100/currentmortgage.RemainingMortgageAmount > 50 :
				Ratio_1 = 75
			   case  currentmortgage.PropertyValuation*100/currentmortgage.RemainingMortgageAmount > 25 :
				Ratio_1 = 50
			   default :
				Ratio_1 = 25
			   }
		           switch {
			   case  currentmortgage.FinancialWorth*100/currentmortgage.RemainingMortgageAmount > 75:
				Ratio_2 = 100
			   case  currentmortgage.FinancialWorth*100/currentmortgage.RemainingMortgageAmount > 50 :
				Ratio_2 = 75
			   case  currentmortgage.FinancialWorth*100/currentmortgage.RemainingMortgageAmount > 25 :
				Ratio_2 = 50
			   default :
				Ratio_2 = 25
			   }
			   switch {
			   case  currentmortgage.CreditScore > 700:
				Ratio_3 = 100
			   case  currentmortgage.CreditScore > 500 :
				Ratio_3 = 75
			   case  currentmortgage.CreditScore > 250 :
				Ratio_3 = 50
			   default :
				Ratio_3 = 25
			   }
			   Rating_Ratio = (Ratio_1 + Ratio_2 + Ratio_3) / 3
			   switch {
			   case Rating_Ratio > 75:
				currentmortgage.RiskClassification = "A"
			   case Rating_Ratio > 50 :
				currentmortgage.RiskClassification = "B"
			   case Rating_Ratio > 25 :
				currentmortgage.RiskClassification = "C"
			   default :
				currentmortgage.RiskClassification = "D"
			   }
			}else {
			     currentmortgage.RiskClassification=""
			}
			switch currentmortgage.RiskClassification{
			    case "A":
			         currentmortgage.RiskAdjustedReturn=currentmortgage.RateofInterest
			    case "B":
			         currentmortgage.RiskAdjustedReturn=currentmortgage.RateofInterest*3/4
			    case "C":
			         currentmortgage.RiskAdjustedReturn=currentmortgage.RateofInterest*2/4
			    case "D":
			         currentmortgage.RiskAdjustedReturn=currentmortgage.RateofInterest*1/4
			    default :
			         currentmortgage.RiskAdjustedReturn=0
			}
			// Calculate Expected Annual CashFlow.
			if currentmortgage.MortgageDuration > 365 {
			   currentmortgage.ExpectedAnnualCashflow=currentmortgage.RemainingMortgageAmount/currentmortgage.MortgageDuration*365
			}else {
			   currentmortgage.ExpectedAnnualCashflow=currentmortgage.RemainingMortgageAmount
			}
			// Calculate if conformed currentmortgage.
			if (currentmortgage.RiskClassification=="A" || currentmortgage.RiskClassification=="B" || currentmortgage.RiskClassification=="C") && currentmortgage.RemainingMortgageAmount <= 424100  && strings.Contains(strings.ToUpper(currentmortgage.MortgageStage),strings.ToUpper("Disbursed:"))  {
			   currentmortgage.ConformedMortgage=true
			}else{
			   currentmortgage.ConformedMortgage=false
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


    // identify place to update mortgage portfolio
		 for counter, value = range mortgages.MortgageNumbers {
				 if value == currentmortgage.MortgageNumber {
						break
			 }
		 }

     // Update mortgage portfolio with new details
		 mortgages.CustomerNames[counter]               = currentmortgage.CustomerName
		 mortgages.MortgageStages[counter]              = currentmortgage.MortgageStage
     mortgages.ConformedMortgages[counter]          = currentmortgage.ConformedMortgage
		 mortgages.MortgagePropertyOwnerships[counter]  = currentmortgage.MortgagePropertyOwnership

		 //package updated Mortgage Portfolio data into bytes.
	 		bytes, err = json.Marshal(mortgages)
	 		if err != nil {
	 			 return nil, errors.New("Add to Mortgage Portfolio record")
	 		}

		//Store updated Mortgage data in blockchain
		mortgagebytes, err = json.Marshal(currentmortgage)
		if err != nil {
			 return nil, errors.New("Error in Marshalling New Mortgage record")
		 }

		err = stub.PutState(string(currentmortgage.MortgageNumber), mortgagebytes)
    if err != nil {
        return nil, err
    }

		//Store updated Mortgage Portfolio in blockchain
		err = stub.PutState("mortgages", bytes)

    return nil, nil
}


// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "retrieve_mortgage_portfolio" {                            //read a variable
           return t.retrieve_mortgage_portfolio(stub, args)
  }else if function == "retrieve_mortgage" {
			     return t.retrieve_mortgage(stub, args)
	}else if function == "retrieve_mortgages" {
			     return t.retrieve_mortgages(stub, args)
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

func (t *SimpleChaincode) retrieve_mortgage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	// Variable declaration
	var mortgage Mortgage
	var err error
	var mortgagebytes []byte


	//Logging
	fmt.Println("running retrieve_mortgage()")

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

	//Get latest mortgages porfolio in blockchain and assign it to struct
	mortgagebytes, err = stub.GetState(string(mortgage.MortgageNumber))
	if err != nil {
			return nil, errors.New("error while fetching mortgage number")
	}
    return mortgagebytes, nil
}

func (t *SimpleChaincode) retrieve_mortgages(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var jsonResp string
    var err error
		var mortgages mortgage_portfolio
		var mortgage_list []Mortgage
		var mortgage Mortgage
		var value int
		var mortgagebytes []byte
    //retrieve Mortgage Portfolio
    valAsbytes, err := stub.GetState("mortgages")
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to retrieve mortgage portfolio\"}"
        return nil, errors.New(jsonResp)
    }
		err = json.Unmarshal(valAsbytes,&mortgages)
		if err != nil {
				return nil, errors.New("error while Unmarshalling mortgages")
		}

	 // identify place to update mortgage portfolio
		for _ , value = range mortgages.MortgageNumbers {
				//Get latest mortgage in blockchain and assign it to variable array
				mortgagebytes, err = stub.GetState(string(value))
				if err != nil {
		 			  return nil, errors.New("error while fetching mortgage number")
		 		}

				err = json.Unmarshal(mortgagebytes,&mortgage)
				if err != nil {
		 			  return nil, errors.New("error while Unmarshalling mortgages for mortgage number")
		 		}else {
					  mortgage_list = append(mortgage_list,mortgage)
			}
		}
		mortgagelist_bytes, err := json.Marshal(mortgage_list)
		if err != nil {
				return nil, errors.New("error while marshalling the mortage list")
		}
    return mortgagelist_bytes, nil
}
