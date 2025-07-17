Create .env at ./cmd/.env with .env.example as example
Ensure go,postgresql, and go-migrate are installed
go to cmd and run `go build .`

POST	localhost:8080/register	Creates a new user account.
POST	localhost:8080/login	Logs in an existing user.


Transactions

GET	  localhost:8080/transaction	        Retrieves user's transaction history.
POST  localhost:8080/transaction/process	Creates a new transaction.
