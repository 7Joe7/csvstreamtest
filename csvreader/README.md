# csvreader

## Responsibilities

The microservice reads a .tar file and searches for .csv files within, if there are such files the microservice parses 
the files and sends parsed contents to dbmanager. So far the only data accepted is of type client (based on headers in the .csv file).

## Structure

Business logic is in service layer, accepting the data is in controller layer. 