# order-inventory-management


## Postman Post Response Script
```js
// {baseurl}/users/login
var response = pm.response.json();

if (response.data && response.data.token) {
    if (response.data.user && response.data.user.role) {
       
        if (response.data.user.role === 'supplier') {
            pm.globals.set("sup_access_token", response.data.token);
            console.log("Supplier token saved to global variable 'sup_access_token':", response.data.token);
        } else if (response.data.user.role === 'customer') {
            pm.globals.set("cust_access_token", response.data.token);
            console.log("Customer token saved to global variable 'cust_access_token':", response.data.token);
        } else {
            console.log("Unknown role:", response.data.user.role);
        }
    } else {
        console.log("User role is missing or null");
    }
} else {
    console.log("Token not found in the response");
}


```