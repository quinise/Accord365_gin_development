const Web3 = require('web3');
const web3 = new Web3('wss://ropsten.infura.io/ws/v3/b522ac003fef4ed2b0d10fc0c7a3de49');

console.log("transaction history bundle loaded");

$(window).on('load', function() {
  function getTenBlocks () {
    var blockArray = [];
    // get latest 10 blocks
    // TODO: Does client really want this feature?
    web3.eth.getBlockNumber().then((latest) => {
      for (let i = 0; i < 10; i++) {
        web3.eth.getBlock(latest - i).then(function(blockNumber) {
          console.log("blockNumber", blockNumber);
          blockArray.push(blockNumber);
          console.log("blockArray ", blockArray);

          return blockArray;
        }).then(function(blockArray){
          console.log("blockArray in promise chain ", blockArray);

          function generateTableHead(table, data) {
            let thead = table.createTHead();
            let row = thead.insertRow();
            for (let key of data) {
              let th = document.createElement("th");
              let text = document.createTextNode(key);
              th.appendChild(text);
              row.appendChild(th);
            }
          }
          
          function generateTable(table, data) {
            for (let element of data) {
              let row = table.insertRow();
              for (key in element) {
                let cell = row.insertCell();
                let text = document.createTextNode(element[key]);
                cell.appendChild(text);
              }
            }
          }

          let table = document.querySelector("#get-blocks-table");
          let data = Object.keys(blockArray[0]);
          generateTableHead(table, data);
          generateTable(table, blockArray);

        }, function(err) {
          alert("Error: unable to load blocks." + err);
          return false;
        });
      }
    })
  }

  getTenBlocks();
  

  $('#getTransactionForm').on('submit', function(event) {
      event.preventDefault();

      function setFormData() {
        formData = Array.from(document.querySelectorAll('#getTransactionForm input')).reduce((acc, input) => ({...acc, [input.id]: input.value}), {});
  
        if (!formData) {
          console.log("form data failed");
          return false;
        }
        console.log("form data ", formData);

        return formData;
      }
  
      var dataToValidate = setFormData();
      console.log("dataToValidate in .on( ", dataToValidate);

     var validatedData = validateData(dataToValidate);
     console.log("validated data .on( ", validatedData);

    web3.eth.getTransaction(validatedData).then(function(transaction) {      
      // post transaction data to page: 
      $('.hash').html(transaction.hash);
      $('.nonce').html(transaction.nonce);
      $('.transaction-from-block').html(transaction.blockHash);
      $('.block-number').html(transaction.blockNumber);
      $('.transaction-index').html(transaction.transactionIndex);
      $('.from').html(transaction.from);
      $('.to').html(transaction.to);
      $('.value').html(transaction.value);
      $('.gas').html(transaction.gas);
      $('.gas-price').html(transaction.gasPrice);
      $('.input').html(transaction.input);
   
    }, function(err) {
        alert("Error: Unable to display transaction" + err);
        return false;
    });
    
    return false;

  });
});

function validateData(dataToValidate) {
  var transactionHashRGEX = /^0x([A-Fa-f0-9]{64})$/;
  var transactionHashResult = transactionHashRGEX.test(dataToValidate.BlockHash);
  if (transactionHashResult) {
    var transactionHashValidated = dataToValidate.BlockHash;
    return transactionHashValidated;
  } else {
    alert("Incorrect transaction hash, please try again. ");
    return false;
  }
}