# Simple-Distributed-File-in-GO
This program is a file server that splits large files into 4 chunks and stores them on 4 different storage servers. It also retrieves files from the storage servers and concatenates the chunks to return the complete file to the client.

#Getting Started
These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

#Prerequisites
You need to have Go installed on your machine.

#Installing
Clone the repository
bash
Copy code
git clone https://github.com/iabdelrhmaneyad/Simple-Distributed-File-in-GO
Change into the directory
bash
Copy code
cd Simple-Distributed-File-in-GO
Build the program
go
Copy code
go build
Start the file server
Copy code
./Simple-Distributed-File-in-GO
Running the tests
No tests are included in this repository.

Deployment
This program is not intended to be used in production environments.

Built With
Go - The programming language used
