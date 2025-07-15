# user
This is the service responsible for user details 

## TODO: 
Following excerpt is from [Auth0](https://auth0.com/blog/adding-salt-to-hashing-a-better-way-to-store-passwords/).
```
In production systems, we store three elements: the salt in cleartext, the resulting hash, and the username. During login, we: 1. Retrieve the stored salt for the username 2. Append it to the provided password 3. Generate a hash 4. Compare with the stored hash value
```
