var Tx = require('ethereumjs-tx').Transaction;
const Web3 = require('web3');
const web3 = new Web3('wss://ropsten.infura.io/ws/v3/b522ac003fef4ed2b0d10fc0c7a3de49');
const Buffer = require('buffer/').Buffer;
const privateKeyToAddress = require('ethereum-private-key-to-address');
const stripHexPrefix = require('strip-hex-prefix');

console.log("Bundle.js loaded");

$(window).on('load', function() {  
  var validatedData = {
    toAddress: "",
    fromAddress: "",
    privateKey: "",
    fullName: "",
    businessName: "",
    dateOfBirth: "",
    caseNumber: "",
    criminalCaseNumber: "",
    dlNumber: "",
    city: "",
    state: "",
    county: "",
    value: 0,
  }  
  
  var txObject = {
    nonce: '',
    from: validatedData.fromAddress,
    to: validatedData.toAddress,
    value: web3.utils.toHex(web3.utils.toWei((validatedData.value).toString(), 'ether')),
    gasLimit: web3.utils.toHex(50000),
    gasPrice: web3.utils.toHex(web3.utils.toWei('10', 'gwei')),
    data: web3.utils.toHex(""),
  };

    $('#childSupportButton').on('click', function() {
      $('#childSupport').toggle();
    });

    $('#probationButton').on('click', function() {
      $('#probation').toggle();
    });

    $('#trafficFeesButton').on('click', function() {
      $('#trafficFees').toggle();
    });

  $('#childSupportForm').on('submit', function(event) {
    event.preventDefault();
    
    function setFormData() {
      formData = Array.from(document.querySelectorAll('#childSupportForm input')).reduce((acc, input) => ({...acc, [input.id]: input.value}), {});

      if (!formData) {
        return false;
      }

      return formData;
    }
    
    dataToValidate = setFormData();
    // Test: console.log("check for dataToValidate ", dataToValidate);
    
    validatedData = validateData(dataToValidate);
    // Test: console.log("called validateData on dataToValidate ", validatedData);
        
    txObject = createTransactionObject(validatedData);
    // Test: console.log("called createTransactionObject ", txObject);
  
    provideValidatedDataTxObject(validatedData, txObject);

    document.getElementById("childSupportForm").reset();

    return false;
  });

  $('#probationForm').on('submit', function(event) {
    event.preventDefault();
    
    function setFormData() {
      formData = Array.from(document.querySelectorAll('#probationForm input')).reduce((acc, input) => ({...acc, [input.id]: input.value}), {});

      if (!formData) {
        return false;
      }

      return formData;
    }
    
    dataToValidate = setFormData();
    // Test: console.log("check for dataToValidate ", dataToValidate);
    
    validatedData = validateData(dataToValidate);
    // Test: console.log("called validateData on dataToValidate ", validatedData);
        
    txObject = createTransactionObject(validatedData);
    // Test: console.log("called createTransactionObject ", txObject);
  
    provideValidatedDataTxObject(validatedData, txObject);

    document.getElementById("probationForm").reset();

    return false;
  });

  $('#trafficFeesForm').on('submit', function(event) {
    event.preventDefault();
    
    function setFormData() {
      formData = Array.from(document.querySelectorAll('#trafficFeesForm input')).reduce((acc, input) => ({...acc, [input.id]: input.value}), {});

      if (!formData) {
        return false;
      }

      return formData;
    }
    
    dataToValidate = setFormData();
    // Test: console.log("check for dataToValidate ", dataToValidate);
    
    validatedData = validateData(dataToValidate);
    // Test: console.log("called validateData on dataToValidate ", validatedData);
        
    txObject = createTransactionObject(validatedData);
    // Test: console.log("called createTransactionObject ", txObject);
  
    provideValidatedDataTxObject(validatedData, txObject);

    document.getElementById("trafficFeesForm").reset();

    return false;
  });

});

function provideValidatedDataTxObject(validatedData, txObject) {
  var noncePromise = new Promise(function(resolve, reject) {
    const fromAccount = validatedData.fromAddress;
    // Testing: console.log("validatedData and txObject in noncePromise = new Promise", validatedData, txObject);
    if (fromAccount) {
      resolve(web3.eth.getTransactionCount(fromAccount));
    }
    
    else {
      alert("Sorry, unable to comlete transaction.");
      reject(Error("Sorry, unable to complete transaction."));
    }
  })

  noncePromise.then(function(result) {
      // Testing completed promimse
      // console.log(result);
      // console.log("validatedData and txObject in noncePromise.then(function ", validatedData, txObject);

      if (!result) {
        console.log("no nonce to return");
        return false;
      } else if(result) {
        txObject.nonce = result;
      }

      // Remove leading 0x from privateKey and create buffer
      var privateKeyBuffer = stripHexPrefix(validatedData.privateKey);
      privateKeyBuffer = Buffer.from(privateKeyBuffer, 'hex');
      
      delete validatedData.privateKey;
      txObject.data = web3.utils.toHex(validatedData)
      // Testing: console.log("txObject in noncepromise.then(**2** ", txObject);

      // Sign the transaction
      const tx = new Tx(txObject, {chain: 'ropsten', hardfork: 'petersburg'});
      tx.sign(privateKeyBuffer);

      const serializedTx = tx.serialize();
      const raw = '0x' + serializedTx.toString('hex');

      // Broadcast the transaction
      web3.eth.sendSignedTransaction(raw, (err, txHash) => {
        // var oldHashes = [];

        if (txHash) {
          // Testing: console.log('txHash:', txHash)
          alert("Check transaction hash on Ethereum blockchain: " +  txHash);
        } else {
          alert("Transaction failed!");
          return false;
        }
      })
  }, function(err) {
      console.log("Error ", err);
      return false;
    });
}
  
function createTransactionObject(validatedData) {
  txObjectCreate = {
      nonce: '',
      from: validatedData.fromAddress,
      to: validatedData.toAddress,
      value: web3.utils.toHex(web3.utils.toWei((validatedData.value).toString(), 'ether')),
      gasLimit: web3.utils.toHex(50000),
      gasPrice: web3.utils.toHex(web3.utils.toWei('10', 'gwei')),
      data: web3.utils.toHex("test")
    };

    txObject = txObjectCreate;
    // Testing: console.log("TxObjectCreate ", txObject);

    return txObject;

}

function validateData(dataToValidate) {
  // To address must be a hex-string
  var toAddressRGEX = /^0x[a-fA-F0-9]{40}$/;
  var toAddressResult = toAddressRGEX.test(dataToValidate.ToAddress);
  if (toAddressResult) {
    var toAddressValidated = dataToValidate.ToAddress;
  } else {
    alert("Invalid recipient address");
    return false;
  }
  
  // validate private key and obtain fromAccount
  try {
    var privateKeyToAddressResult = privateKeyToAddress(dataToValidate.PrivateKey);
    var fromAdressValidated = privateKeyToAddressResult;
    // Testing: console.log("fromAddress", fromAdressValidated);
    var privateKeyValidated = dataToValidate.PrivateKey;
    
  } catch(err) {
    alert("Invalid private key " + err);
    return false;
    
  }
   
  // Full name validation must be only characters and at least two units in length
  var fullNameRGEX = /^[ a-zA-Z\-\’]+$/;
  var fullNameResult = fullNameRGEX.test(dataToValidate.FullName);
  if (fullNameResult) {
    var fullNameValidated = dataToValidate.FullName;
  } else {
    alert("Invalid full name");
    return false;
  }

  // Business name validation must be only characters and at least two units in length, "/" included
  var businessNameRGEX = /^[\w]+([-_\s]{1}[a-z0-9]+)*$/;
  var businessNameResult = businessNameRGEX.test(dataToValidate.businessName);
  if (businessNameResult) {
    var businessNameValidated = dataToValidate.businessName;
  } else {
    alert("Invalid business name");
    return false;
  }

  // Date of Birth validation, must be a date
  var d = new Date((dataToValidate.dateOfBirth));
  var now = new Date();
  if ((d.getTime()-now.getTime()) > 0) {
    alert("Invalid date of birth");
    return false;
  } else {
    var dateOfBirthValidated = dataToValidate.dateOfBirth;
  }

  // Date of citation, must be a date in the past
  var d = new Date(dataToValidate.dateOfCitation);
  var now = new Date();
  if ((d.getTime()-now.getTime()) > 0) {
    alert("Invalid date of citation");
    return false;
  } else {
    var dateOfCitationValidated = dataToValidate.dateOfCitation;
  }

  // Case number must be an alphanumeric string of length 15 (not case sensitive)
  var caseNumberRGEX = /\w{15}/;
  var caseNumberResult = caseNumberRGEX.test(dataToValidate.caseNumber);
  if (caseNumberResult) {
    var caseNumberValidated = dataToValidate.caseNumber;
  } else {
    alert("Invalid case number");
    return false;
  }

 // Criminal Case Number must be an alphanumeric string of length 15 (not case sensitive)
 var criminalCaseNumberRGEX = /\w{15}/;
 var criminalCaseNumberResult = criminalCaseNumberRGEX.test(dataToValidate.criminalCaseNumber);
 if (criminalCaseNumberResult) {
   var criminalCaseNumberValidated = dataToValidate.criminalCaseNumber;
  } else {
    alert("Invalid criminal case number");
    return false;
  }

  // Drivers License Number must be an alphanumeric string of length 15 (not case sensitive)
  var dlNumberRGEX = /\w{15}/;
  var dlNumberResult = dlNumberRGEX.test(dataToValidate.dlNumber);
  if (dlNumberResult) {
    var dlNumberValidated = dataToValidate.dlNumber;
  } else {
    alert("Invalid drivers license number");
    return false;
  }

  // City (of citation) validation must be only characters and at least two units in length
  var cityRGEX = /^[ a-zA-Z\-\’]+$/;
  var cityResult = cityRGEX.test(dataToValidate.city);
  if (cityResult) {
    var cityValidated = dataToValidate.city;
  } else {
    alert("Invalid city name");
    return false;
  }

  // State must be a string of length 2
  var stateRGEX = /[A-Z]{2}/;
  var stateResult = stateRGEX.test(dataToValidate.state);
  if (stateResult) {
    var stateValidated = dataToValidate.state;
  } else {
    alert("Invalid state name");
    return false;
  }

  // County must be a string
  var countyRGEX = /^[A-Za-z]+$/;
  var countyResult = countyRGEX.test(dataToValidate.county);
  if (countyResult) {
    var countyValidated = dataToValidate.county;
  } else {
    alert("Invalid county name");
    return false;
  }

  // Ethereum transaction validation must be a number (decimal accepted)
  var valueRGEX = /^\d{0,9}(?:[.]\d{0,9})?$/;
  var valueResult = valueRGEX.test(dataToValidate.value);
  if (valueResult) {
    var valueValidated = dataToValidate.value;
  } else {
    alert("Invalid Ether value");
    return false;
  }

  validatedDataObj = {
    transactionTitle: transactionTitleValidated,
    toAddress: toAddressValidated,
    fromAddress: fromAdressValidated,
    privateKey: privateKeyValidated,
    fullName: fullNameValidated,
    businessName: businessNameValidated,
    dateOfBirth: dateOfBirthValidated,
    dateOfCitation: dateOfCitationValidated,
    caseNumber: caseNumberValidated,
    criminalCaseNumber: criminalCaseNumberValidated,
    dlNumber: dlNumberValidated,
    city: cityValidated,
    state, stateValidated,
    value: valueValidated,
  }

  validatedData = validatedDataObj

  // Testing: console.log("validatedData ", validatedData);

  return validatedData;
}