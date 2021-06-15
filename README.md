# E-Commerce-Application
E - Commerce Application

* The aim is to get familiarized with microservices , coding practices and devops
* tools. Create three(or more) microservices which are capable of communicating with
* each other. The communication can be as REST API calls. For example , inventory,
* cart and payment are three different microservices. Inventory services takes care of
* maintaining the items inventory. Cart service takes care of adding item to the cart ,
* displaying items in cart etc . And payment takes care of confirming the orders
* payment. Always think of developing the application in an incremental way.
* 
* Phase 1 Microservices driven by Test cases
* Step 1: Choose your database(SQL or NoSql) and model the data . (Note: you
* should be modelling data keeping in mind all the necessary entities . For example ,
* you need two types of users - customers and Admins .Customer is the real customer
* who holds a cart. Admin is someone who loads the inventory into the system
* Step 2: Build your microservices and expose the necessary Rest APIs . Eg: v1/
* update/inventory (for admins to add inventory) , v1/add/item (for customer adding
* item to the cart) etc
* Step 3: Create test case which cover 100 customers checkout flow in parallel.
* That is - Add item to the cart - View cart- payment done- Order confirmed - inventory
* reduced(As the item is bought by someone). As there is no UI to load initial data,
* your test cases should do a setup of test data before running the tests.
* Step 4: Create docker compose - This should run the above created tests after
* bringing up docker containers of all your microservices and the database
* Notes:
* 1. Name the application . Try to understand best practices. Always think of best
* performance of application
* 2. You can ignore Authentication in Phase 1. Only thing you can check is
* availability of the user /cart/item in the DB before doing any operations.
* 3. Test your endpoints during development using the SOAPUI or any similar
* tools
* 4. Try to solve Race condition if possible which can occur in Step 3. Your
* microservices should be handling concurrency . Eg: if 100 people are trying to
* access an item which is having only 5 available quantity, all customers should be
* able to add item to the cart, but only the first 5 who completed payment will get the
* item. Remaining 95 people will get an error message on or before payment that the
* item is not available.
* 5. No need to think much about payment- Make a dummy payment when the
* service is called.
* 
* Phase 2: Application beta with basic UI
* In this phase we can focus on bringing the application to the real users.
* Step 1: Create a microservice which is for the UI of the application. Create
* basic UI fields .Eg: A login screen -> A page which shows all items with a button to
* add item to cart -> A page which shows the customers cart-> A dummy payment
* page -> A order confirmation page.
* Step3: Add authentication for customer login
* 
* Phase 3: Application in cloud
* Deploy the application in cloud.
* <<Steps to be updated>>
* Technologies/concepts
* Golang and Golang frameworks
* Docker
* Docker compose
* Microservices
* Rest API
* SOL/NoSQL DB
* Unit test
* IBM cloud
* Kubernetes
* Tools
* GIT/BIT bucket
* JIRA/any other
* SOAP UI
