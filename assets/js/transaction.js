var Tx = require('ethereumjs-tx').Transaction;
const Web3 = require('web3');
const web3 = new Web3('wss://ropsten.infura.io/ws/v3/b522ac003fef4ed2b0d10fc0c7a3de49');
const Buffer = require('buffer/').Buffer;
const privateKeyToAddress = require('ethereum-private-key-to-address');
const stripHexPrefix = require('strip-hex-prefix');

console.log("Bundle.js loaded");

$(window).on('load', function() { 
  // ToDo: add value input, wait for wallet to be implemented so we can use tokens
  var validatedData = {
    transactionTitle: "",
    toAddress: "",
    fromAddress: "",
    privateKey: "",
    fullName: "",
    email: "",
  }  
  
  var txObject = {
    nonce: '',
    from: validatedData.fromAddress,
    to: validatedData.toAddress,
    value: web3.utils.toHex(web3.utils.toWei('.01', 'ether')),
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

        if (txHash){
          // Testing: console.log('txHash:', txHash)
          alert("Check transaction hash on Ethereum blockchain: " +  txHash);
          // TODO: prevent form from submitting on pageload
          function SubmitForm(event) {
            event.preventDefault();

            document.getElementById("txValueHidden").value = txHash;
            formData = Array.from(document.querySelectorAll('#txHashForm input')).reduce((acc, input) => ({...acc, [input.id]: input.value}), {});            document.getElementById("txHashForm").submit();
            if (!formData) {
              return false;
            }
            console.log("formData ", formData)
            return formData;
          };
          SubmitForm()
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
      value: web3.utils.toHex(web3.utils.toWei('.01', 'ether')),
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
  
  // //Phone number must be only numbers
  // dataToValidate.phoneNumber = Number(dataToValidate.phoneNumber)
  // var phoneRGEX = /^\d{10}$/;
  // var phoneResult = phoneRGEX.test(dataToValidate.phoneNumber);
  // if (phoneResult) {
  //   var phoneNumberValidated = dataToValidate.phoneNumber;
  // } else {
  //    alert("phone:"+phoneResult );
  // }

  validatedDataObj = {
    transactionTitle: transactionTitleValidated,
    toAddress: toAddressValidated,
    fromAddress: fromAdressValidated,
    privateKey: privateKeyValidated,
    fullName: fullNameValidated,
    email: emailValidated,
  }

  validatedData = validatedDataObj

  // Testing: console.log("validatedData ", validatedData);

  return validatedData;
}