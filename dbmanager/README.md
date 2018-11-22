# dbmanager

## Responsibilities

The microservice listens on port 7777 (configurable) for client imports via grpc and saves values into MySQL DB.

## Structure

Business logic is in service layer, connector to the MySQL is in repository layer. There is an initial migration. 