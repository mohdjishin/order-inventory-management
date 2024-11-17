TODOS:


<!-- - RBAC should be implemented in.  -->
  <!-- - plan is super admin will approve supplier. buyer  can login and see the products.// -->
  <!-- - all users will be wrapped with jwt token. -->


- Asumption settlements (payments) will be done outside the platform. 
  - so no need to implement payment gateway. 



- admin
  <!-- - blacklist supplier and all listed products. //  -->
  <!-- - update on all get buy products. apis like no record found. or something on blacklist products. -->
  - whitelist supplier
  - statistics of products.


- supplier
   <!-- 1. add product -->
   2. update product
   <!-- 3. delete product -->
   <!-- 4. Shipmemt once order is approved by by supplier.  -->
   <!-- 4. view all products -->
    <!-- 5. add stocks to existing products -->
    <!-- 6. approve or reject orders by customers. -->


- buyer
  - 1. return product

- bugs
 - 1. when blaclisting by admin. admins email id should be stored in the database. to track who blacklisted the supplier.
 - 2. when supplier is not approved. he should not be able to login. (currently he can login but can not do other operations - kept like this as functitonality as of now.)
 - 3. once approved supplier is able to do approve again. fix this. 
 - 4.  when product delivery is cancelled and stock is increased. price of the product should be reduced write a db trigger or on cancel api write a logic to reduce the price.
 - 4. delivery address should be stored in the database. as of now not need not very important. for task completion.
 - 5. password length and complexity should be checked. as of now not very important. for task completion.
