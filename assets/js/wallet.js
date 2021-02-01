const Portis = require('@portis/web3');
const Web3 = require('web3');
const portis = new Portis('5946cfd4-f46d-4b09-b579-020a69690dee', 'rinkeby');
const web3 = new Web3(portis.provider);


$(window).on('load', function() {
    // Leads to the ability to buy and recieve Ether/other currencies
    document.getElementById("showPortis").onclick = () => portis.provider.enable();
    document.getElementById("Logout").onclick = () => portis.logout();

    portis.onLogin((walletAddress) => {
        document.getElementById("portisAddress").innerHTML = `
        <div> Wallet Address: ${walletAddress} </div>
      `;
      });

      portis.onActiveWalletChanged(walletAddress => {
        console.log('Active wallet address:', walletAddress);
      });
      
      portis.onLogout(() => {
        alert('User logged out');
      });

      portis.showPortis();

      portis.onError(error => {
        console.log('error', error);
      });
});