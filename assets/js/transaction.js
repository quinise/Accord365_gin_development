var Tx = require('ethereumjs-tx').Transaction;
const Web3 = require('web3');
const web3 = new Web3('wss://ropsten.infura.io/ws/v3/b522ac003fef4ed2b0d10fc0c7a3de49');
const Buffer = require('buffer/').Buffer;
const privateKeyToAddress = require('ethereum-private-key-to-address');
const stripHexPrefix = require('strip-hex-prefix');

console.log("transactionBundle.js loaded");

$(window).on('load', function() { 
  var validatedData = {
    transactionTitle: "",
    toAddress: "",
    fromAddress: "",
    privateKey: "",
    fullName: "",
    pocFullName: "",
    pocPhoneNumber: "",
    pocBusinessEmail: "",
    businessName: "",
    hqPhoneNumber: "",
    beWebsite: "",
    businessEmail: "",
    feinNumber: "",
    seinNumber: "",
    taxIDNumber: "",
    beAdr: "",
    beCity: "",
    beState: "",
    beZip: "",
    beCountryCode: "",
    doingAdr: "",
    dbaCity: "",
    dbastate: "",
    dbaZip: "",
    value: 0,
  }  
  
  var txObject = {
    nonce: '',
    from: validatedData.fromAddress,
    to: validatedData.toAddress,
    value: web3.utils.toHex(web3.utils.toWei((validatedData.value).toString(), 'ether')),
    gasLimit: web3.utils.toHex(50000),
    gasPrice: web3.utils.toHex(web3.utils.toWei('10', 'gwei')),
    data: web3.utils.toHex("test"),
  };
  
  $('#newPaymentForm').on('submit', function(event) {
    
    event.preventDefault();
    
    function setFormData() {
      formData = Array.from(document.querySelectorAll('#newPaymentForm input')).reduce((acc, input) => ({...acc, [input.id]: input.value}), {});

      if (!formData) {
        return false;
      } else {
        return formData;
      }
    }
    
    dataToValidate = setFormData();
    // Test: console.log("check for dataToValidate ", dataToValidate);
    
    validatedData = validateData(dataToValidate);
    // Test: console.log("called validateData on dataToValidate ", validatedData);
        
    txObject = createTransactionObject(validatedData);
     //Test: console.log("called createTransactionObject ", txObject);
  
    provideValidatedDataTxObject(validatedData, txObject);

    document.getElementById("newPaymentForm").reset();

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
        if (txHash) {
          // Testing: console.log('txHash:', txHash)
          alert("Check transaction hash on Ethereum blockchain: " +  txHash);
          $(function() {
            $('#hash-container').load('../templates/new_payment.html #transaction-hashes', function() {
              formData = Array.from(document.querySelectorAll('#txHashForm input')).reduce((acc, input) => ({...acc, [input.id]: input.value}), {});

              if (!formData) {
              alert("no form data");
              return false;
            } else {
              document.getElementById("txValueHidden").value = txHash;
              formData = Array.from(document.querySelectorAll('#txHashForm input')).reduce((acc, input) => ({...acc, [input.id]: input.value}), {});            
              document.getElementById("txHashForm").submit();
              document.getElementById("txValueHidden").value = "";
              console.log("formData ", formData);
              return formData;
            }
          });
        }); 
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
  // Transaction Title cannot contain numbers, must only contain alphabetic characters, no special characters 
  var transactionTitleRGEX = /^[ a-zA-Z\-\’]+$/;
  var transactionTitleResult = transactionTitleRGEX.test(dataToValidate.TransactionTitle);

  if (transactionTitleResult) {
    var transactionTitleValidated = dataToValidate.TransactionTitle;
  } else {
    alert("Invalid transaction title");
    return false;
  }

  // To address must be a hex-string of length 40
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
   
  // Full name (of the person doing business) validation must be only characters
  var fullNameRGEX = /^[ a-zA-Z\-\’]+$/;
  var fullNameResult = fullNameRGEX.test(dataToValidate.FullName);
  if (fullNameResult) {
    var fullNameValidated = dataToValidate.FullName;
  } else {
    alert("Invalid full name os the person doing business");
    return false;
  }

    // Point of contact full name validation must be only characters
    var pocFullNameRGEX = /^[ a-zA-Z\-\’]+$/;
    var pocFullNameResult = pocFullNameRGEX.test(dataToValidate.pocFullName);
    if (pocFullNameResult) {
      var pocFullNameValidated = dataToValidate.pocFullName;
    } else {
      alert("Invalid point of contact full name");
      return false;
    }
  
    // Point of contact phone number must be only numbers
    dataToValidate.pocPhoneNumber = Number(dataToValidate.pocPhoneNumber)
    var pocPhoneNumberRGEX = /^\d{10}$/;
    var pocPhoneNumberResult = pocPhoneNumberRGEX.test(dataToValidate.pocPhoneNumber);
    if (pocPhoneNumberResult) {
      var pocPhoneNumberValidated = dataToValidate.pocPhoneNumber;
    } else {
      alert("Invalid point of contact phone number");
    }

  // Point of contact email validation must be a valid email address
  var pocBusinessEmailRGEX = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/;
  var pocBusinessEmailResult = pocBusinessEmailRGEX.test(dataToValidate.pocBusinessEmail);
  if (pocBusinessEmailResult) {
    var pocBusinessEmailValidated = dataToValidate.pocBusinessEmail;
  } else {
    alert("Invalid point of contact email");
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
  
  // Headquarters phone number must be only numbers
  dataToValidate.hqPhoneNumber = Number(dataToValidate.hqPhoneNumber)
  var hqPhoneNumberRGEX = /^\d{10}$/;
  var hqPhoneNumberResult = hqPhoneNumberRGEX.test(dataToValidate.hqPhoneNumber);
  if (hqPhoneNumberResult) {
    var hqPhoneNumberValidated = dataToValidate.hqPhoneNumber;
  } else {
    alert("Invalid business headquarters phone number");
  }

  // Business website validation must be a valid url
  var beWebsiteRGEX = /^(?:(?:https?|ftp):\/\/)?(?:(?!(?:10|127)(?:\.\d{1,3}){3})(?!(?:169\.254|192\.168)(?:\.\d{1,3}){2})(?!172\.(?:1[6-9]|2\d|3[0-1])(?:\.\d{1,3}){2})(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-z\u00a1-\uffff0-9]-*)*[a-z\u00a1-\uffff0-9]+)(?:\.(?:[a-z\u00a1-\uffff0-9]-*)*[a-z\u00a1-\uffff0-9]+)*(?:\.(?:[a-z\u00a1-\uffff]{2,})))(?::\d{2,5})?(?:\/\S*)?$/;
  var beWebsiteResult = beWebsiteRGEX.test(dataToValidate.beWebsite);
  if (beWebsiteResult) {
    var beWebsiteValidated = dataToValidate.beWebsite;
  } else {
    alert("Invalid business website url");
    return false;
  }

   // Business headquarters email validation must be a valid email address
   var businessEmailRGEX = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/;
   var businessEmailResult = businessEmailRGEX.test(dataToValidate.businessEmail);
   if (businessEmailResult) {
     var businessEmailValidated = dataToValidate.businessEmail;
   } else {
     alert("Invalid business headquarters email");
     return false;
   }

  // Business FEIN number validation must be a 9 digit number
  dataToValidate.feinNumber = Number(dataToValidate.feinNumber)
  var feinNumberRGEX = /^\d{9}$/;
  var feinNumberResult = feinNumberRGEX.test(dataToValidate.feinNumber);
  if (feinNumberResult) {
    var feinNumberValidated = dataToValidate.feinNumber;
  } else {
    alert("Invalid FEIN number");
  }

  // Business SEIN number validation must be a 8 digit number
  dataToValidate.seinNumber = Number(dataToValidate.seinNumber)
  var seinNumberRGEX = /^\d{8}$/;
  var seinNumberResult = seinNumberRGEX.test(dataToValidate.seinNumber);
  if (seinNumberResult) {
    var seinNumberValidated = dataToValidate.seinNumber;
  } else {
    alert("Invalid SEIN number");
  }

  // Business TaxID number validation must be a 9 digit number that starts with the number 9
  dataToValidate.taxIDNumber = Number(dataToValidate.taxIDNumber)
  var taxIDNumberRGEX = /^[9]\d{8}$/;
  var taxIDNumberResult = taxIDNumberRGEX.test(dataToValidate.taxIDNumber);
  if (taxIDNumberResult) {
    var taxIDNumberValidated = dataToValidate.taxIDNumber;
  } else {
    alert("Invalid tax ID number");
  }

  // Business headquarters address validation must be a valid street address
  var beAdrRGEX = /^[a-zA-Z0-9\s,'-]*$/;
  var beAdrResult = beAdrRGEX.test(dataToValidate.beAdr);
  if (beAdrResult) {
    var beAdrValidated = dataToValidate.beAdr;
  } else {
    alert("Invalid business address");
    return false;
  }

  // Business headquarters city (of citation) validation must be only characters and at least two units in length
  var beCityRGEX = /^[ a-zA-Z\-\’]+$/;
  var beCityResult = beCityRGEX.test(dataToValidate.beCity);
  if (beCityResult) {
    var beCityValidated = dataToValidate.beCity;
  } else {
    alert("Invalid business city name");
    return false;
  }

  // Business headquarters state must be a string of length 2
  var beStateRGEX = /[A-Z]{2}/;
  var beStateResult = beStateRGEX.test(dataToValidate.beState);
  if (beStateResult) {
    var beStateValidated = dataToValidate.beState;
  } else {
    alert("Invalid state name");
    return false;
  }

  // Business zipcode validation must be a number of lenght 5
  var beZipRGEX = /[0-9]{5}/;
  var beZipResult = beZipRGEX.test(dataToValidate.beZip);
  if (beZipResult) {
    var beZipValidated = dataToValidate.beZip;
  } else {
    alert("Invalid zip-code for business");
    return false;
  }

  // Business country code validation must be 2 numbers
  var beCountryCodeRGEX = /^(\+?\d{1,3}|\d{1,4})$/;
  var beCountryCodeResult = beCountryCodeRGEX.test(dataToValidate.beCountryCode);
  if (beCountryCodeResult) {
    var beCountryCodeValidated = dataToValidate.beCountryCode;
  } else {
    alert("Invalid country code for business");
    return false;
  }

  // Doing Business As address validation must be a valid street address
  var doingAdrRGEX = /^[a-zA-Z0-9\s,'-]*$/;
  var doingAdrResult = doingAdrRGEX.test(dataToValidate.doingAdr);
  if (doingAdrResult) {
    var doingAdrValidated = dataToValidate.doingAdr;
  } else {
    alert("Invalid doing business as address");
    return false;
  }

  // Doing Business As City (of citation) validation must be only characters and at least two units in length
  var dbaCityRGEX = /^[ a-zA-Z\-\’]+$/;
  var dbaCityResult = dbaCityRGEX.test(dataToValidate.dbaCity);
  if (dbaCityResult) {
    var dbaCityValidated = dataToValidate.dbaCity;
  } else {
    alert("Invalid doing business as city name");
    return false;
  }

  // Doing Business As State must be a string of length 2
  var dbaStateRGEX = /[A-Z]{2}/;
  var dbaStateResult = dbaStateRGEX.test(dataToValidate.dbaState);
  if (dbaStateResult) {
    var dbaStateValidated = dataToValidate.dbaState;
  } else {
    alert("Invalid doing business as state");
    return false;
  }

  // Doing Business As Business zipcode validation must be a number of lenght 5
  var dbaZipRGEX = /[0-9]{5}/;
  var dbaZipResult = dbaZipRGEX.test(dataToValidate.dbaZip);
  if (dbaZipResult) {
    var dbaZipValidated = dataToValidate.dbaZip;
  } else {
    alert("Invalid zip-code for doing business as");
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
    pocFullName: pocFullNameValidated,
    pocPhoneNumber: pocPhoneNumberValidated,
    pocBusinessEmail: pocBusinessEmailValidated,
    businessName: businessNameValidated,
    hqPhoneNumber: hqPhoneNumberValidated,
    beWebsite: beWebsiteValidated,
    businessEmail: businessEmailValidated,
    feinNumber: feinNumberValidated,
    seinNumber: seinNumberValidated,
    taxIDNumber: taxIDNumberValidated,
    beAdr: beAdrValidated,
    beCity: beCityValidated,
    beState: beStateValidated,
    beZip: beZipValidated,
    beCountryCode: beCountryCodeValidated,
    doingAdr: doingAdrValidated,
    dbaCity: dbaCityValidated,
    dbaState: dbaStateValidated,
    dbaZip: dbaZipValidated,
    value: valueValidated,
  }

  validatedData = validatedDataObj

  // Testing: console.log("validatedData ", validatedData);

  return validatedData;
}