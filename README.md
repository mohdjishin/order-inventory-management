# order-inventory-management

**_Important:_** **Don't worry about database scripts; GORM will handle the migrations. Start by creating a user.**

## Installation and Setup

### Prerequisites

1. Docker

### Steps

1. Clone the repository
2. Run `docker compose up --build` in the root directory of the project

## Step to run locally without Docker

### Prerequisites

1. Go 
2. PostgreSQL

### Steps

1. Clone the repository
2. Set up PostgreSQL and create a database
3. Update the database configuration in `config.json`
4. Run `go run cmd/api/*.go` in the root directory of the project




## Postman Post Response Script
```js
// {baseurl}/users/login. this is a helper script to save the token in global variable. (make life simpler with scripts)
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