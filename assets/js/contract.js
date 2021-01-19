var Tx = require('ethereumjs-tx').Transaction;
const Web3 = require('web3');
const web3 = new Web3('wss://ropsten.infura.io/ws/v3/b522ac003fef4ed2b0d10fc0c7a3de49');
const Buffer = require('buffer/').Buffer;
const privateKeyToAddress = require('ethereum-private-key-to-address');
const stripHexPrefix = require('strip-hex-prefix');

console.log("Bundle.js loaded");

$(window).on('load', function() {  
  var validatedData = {
    transactionTitle: "",
    toAddress: "",
    fromAddress: "",
    privateKey: "",
    fullName: "",
    email: "",
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
  // Transaction Title cannot contain numbers, must only contain alphabetic characters, no special characters 
  var transactionTitleRGEX = /^[ a-zA-Z\-\’]+$/;
  var transactionTitleResult = transactionTitleRGEX.test(dataToValidate.TransactionTitle);

  if (transactionTitleResult) {
    var transactionTitleValidated = dataToValidate.TransactionTitle;
  } else {
    alert("Invalid transaction title");
    return false;
  }

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
  
  // Email validation must be a valid email address
  var emailRGEX = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/;
  var emailResult = emailRGEX.test(dataToValidate.Email);
  if (emailResult) {
    var emailValidated = dataToValidate.Email;
  } else {
    alert("Invalid email");
    return false;
  }

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
    email: emailValidated,
    value: valueValidated,
  }

  validatedData = validatedDataObj

  // Testing: console.log("validatedData ", validatedData);

  return validatedData;
}