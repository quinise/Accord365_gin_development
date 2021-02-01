// TODO: (1) Remove console.logs 
// (2) update network to ethereum main network
Web3 = require('web3')

App = {
  web3Provider: null,
  providerUrl: 'ws://127.0.0.1:7545',
  contracts: {},
  account: '0x0',
  loading: false,
  tokenPrice: 1000000000000000,
  tokensSold: 0,
  tokensAvailable: 750000,


  // Initializing function logs to the console if js is working, calls web3 configuration
  init: function() {
    console.log("App initialized...");
    return App.initWeb3();
  },

  // Sets blockchain provider and calls for contracts in the app to be inititated
  initWeb3: function() {
     if (typeof web3 !== 'undefined') {
       // If a web3 instance is already provided by Meta Mask.
       App.web3Provider = web3.currentProvider;
       web3 = new Web3(web3.currentProvider);
     } else {
      // Specify default instance if no web3 instance provided
      App.web3Provider = new Web3.providers.WebsocketProvider(App.providerUrl);
      web3 = new Web3(App.web3Provider);
     }
    return App.initContracts();
  },
 
  // Using the JSON generation of a contract js envokes a token sale (then later the token) contract using Truffle/contract
  // sets the provider, and deploys the contract
  initContracts: function() {
    $.getJSON("build/contracts/DappTokenSale.json", function(dappTokenSale) {
      App.contracts.DappTokenSale = TruffleContract(dappTokenSale);
      App.contracts.DappTokenSale.setProvider(App.web3Provider);
      App.contracts.DappTokenSale.deployed().then(function(dappTokenSale) {
      console.log("Dapp Token Sale Address", dappTokenSale.address);
      });
    }).done(function() {
        $.getJSON("build/contracts/DappToken.json", function(dappToken) {
        App.contracts.DappToken = TruffleContract(dappToken);
        App.contracts.DappToken.setProvider(App.web3Provider);
        App.contracts.DappToken.deployed().then(function(dappToken) {
          console.log("Dapp Token Address:", dappToken.address);
          });

          App.listenForEvents();
          return App.render();
        });
      })
    },

    // Listen for events emitted from the contract
    listenForEvents: function() {
      App.contracts.DappTokenSale.deployed().then(function(instance) {
        instance.contract.events.Sell( {}, {
          fromBlock: 0,
          toBlock: 'latest',
      }, function(error, event) { console.log(event); App.render();}).on("connected", function(event){
            console.log("event triggered ", event);
        }).on('data', function(event) {
            console.log("data" ,event); // same results as the optional callback above
        }).on('error', function(error, receipt) { // If the transaction was rejected by the network with a receipt, the second parameter will be the receipt.
            console.log("error", error);
        });
      })
    },

    render: function() {
      if (App.loading) {
        return;
      }
      App.loading = true;

      var loader  = $('#loader');
      var content = $('#content');
  
      loader.show();
      content.hide();

      //load account data
      web3.eth.getCoinbase(function(err, account) {
        if(err === null) {
          App.account = account;
          $('#accountAddress').html("Your Account: " + account);
        }
      });

      // Load token sale contract
      App.contracts.DappTokenSale.deployed().then(function(instance) {
        dappTokenSaleInstance = instance;
        return dappTokenSaleInstance.tokenPrice();
      }).then(function(tokenPrice) {
        App.tokenPrice = web3.utils.toBN(tokenPrice);
        App.tokenPrice = web3.utils.fromWei(App.tokenPrice, "ether");
        $('.token-price').html(App.tokenPrice);
        return dappTokenSaleInstance.tokensSold();
      }).then(function(tokensSold) {
        App.tokensSold = Number(tokensSold);
        $('.tokens-sold').html(App.tokensSold);
        $('.tokens-available').html(App.tokensAvailable);

        var progressPercent = ((Math.ceil(App.tokensSold) / App.tokensAvailable) * 100);
        $('#progress').css('width', progressPercent + '%');

          // Load Token Contract
          App.contracts.DappToken.deployed().then(function(instance) {
            dappTokenInstance = instance;
            return dappTokenInstance.balanceOf(App.account);
          }).then(function(balance){
            $('.dapp-balance').html(Number(balance));
            App.loading = false;
            loader.hide();
            content.show();
          })
      });
    },

    buyTokens: function() {
      $('#content').hide();
      $('#loader').show();
      var numberOfTokens = Number($('#numberOfTokens').val());
      App.tokenPrice = web3.utils.toWei((App.tokenPrice).toString(), 'ether');
      App.tokenPrice = web3.utils.toBN(App.tokenPrice);
      var tokenCost = (numberOfTokens * App.tokenPrice);      
      
      App.contracts.DappTokenSale.deployed().then(function(instance) {
        return instance.buyTokens(numberOfTokens, { 
          from: App.account,
          value: tokenCost,
          gas: 500000 // Gas limit
        });
      }).then(function(result) {
        console.log("tokens bought...");
        $('form').trigger('reset'); //reset number of tokens in form
        // wait for sell event to be triggered
      });
    }
  }

  $(window).on('load', function() {
    App.init();
  });