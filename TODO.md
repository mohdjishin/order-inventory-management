TODOS:


- RBAC should be implemented in. 
  - plan is super admin will approve supplier. buyer  can login and see the products.
  - currern plan platform admin will remain unauthonticated. as of now thinking to create api key or something else. not sure as of now.
  - remaining users will be wrapped with jwt token.


- Asumption settlements (payments) will be done outside the platform. 
  - so no need to implement payment gateway. 



- admin
  - blacklist supplier and all listed products. // 
  - update on all get buy products. apis like no record found. or something on blacklist products.
  - whitelist supplier

