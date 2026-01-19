# E commerce Checkout

tech stack : temporal + go + nats + psql_12

#### Why Checkout

I choose checkout system coz its more of a distributed transaction, according to me payment and inventory should be in different deployables. Temporal gives the durable execution and NATs provides us with the back pressure and fast pub/sub

Its pretty neat for me to work with. 