# order-inventory-management


## Postman Post Response Script
```js
// {baseurl}/users/login
var response = pm.response.json();


if (response.data && response.data.token) {

    pm.globals.set("access_token", response.data.token);
    console.log("Token saved to global variable access_token:", response.data.token);
} else {
    console.log("Token not found in the response");
}

```