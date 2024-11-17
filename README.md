# order-inventory-management

**_Important:_** **Don't worry about database scripts; GORM will handle the migrations.**

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
4. Run 1`./build.sh` in the root directory of the project
5. Run `./oim` to start the server
```bash
$ ./build.sh
$ ./oim
```






## Postman Post Response Script (Helper Script)
```js
// {baseurl}/users/login. this is a helper script to save the token in global variable. (make life simpler with scripts)
const response = pm.response.json();

if (response.data && response.data.token) {
    const user = response.data.user;

    if (user && user.role) {
        switch (user.role) {
            case 'supplier':
                pm.globals.set("sup_access_token", response.data.token);
                console.log("Supplier token saved to global variable 'sup_access_token':", response.data.token);
                break;
            case 'customer':
                pm.globals.set("cust_access_token", response.data.token);
                console.log("Customer token saved to global variable 'cust_access_token':", response.data.token);
                break;
            case 'admin':
                pm.globals.set("adm_access_token", response.data.token);
                console.log("Admin token saved to global variable 'adm_access_token':", response.data.token);
                break;
            default:
                console.log("Unknown role:", user.role);
        }
    } else {
        console.log("User role is missing or null");
    }
} else {
    console.log("Token not found in the response");
}

//  use these global variables in the headers like this bearer {{sup_access_token}} or in authorization tab. 
```





## Pre-Activity Checklist

- Use the following endpoints to verify the service status and details:

  1. **Health Check**: Confirms if the service is running:
     ```bash
     curl -X GET http://localhost:8080/health
     ```
  2. **Info**: Displays server information:
     ```bash
     curl -X GET http://localhost:8080/info
     ```





<br>
<br>
<br>


## No Admin Creation Endpoint

There is currently no endpoint available for creating an Admin. However, the system is designed to automatically generate a default Super Admin user during the initial setup.

Here are the default credentials for the admin user:
- **Email**: `admin@oim.com`
- **Password**: `passwOrd@123`