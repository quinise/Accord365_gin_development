$(window).on('load', function(){ 
    var element = document.getElementById('currencies');
    var paymentHtml = '<label for="contractPaymentValue"><i class="fa fa-envelope"></i> Payment Value</label>' +
    '<input type="number" id="contractPaymentValue" name="contractPaymentValue" placeholder="0">';
    $(element).on({
    "focus": function() {
        console.log('clicked!', this, this.value);
        this.selectedIndex = -1;
    },
    "change": function() {
        choice = $(this).val();
        console.log('changed!', this, choice);
        this.blur();
        document.getElementById("paymentInfo").innerHTML = paymentHtml;
    }
    });
})