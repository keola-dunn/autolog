# Schema
This dir contains the DB migrations for the postgres db powering the autlog application.

## Example 
```.sh
goose postgres "user=postgres password=postgres host=127.0.0.1 port=5432  dbname=postgres" up
```
