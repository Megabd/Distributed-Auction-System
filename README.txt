# Distributed-Auction-System

How to run the Distributed Auction System:

1. Go to the Distributed-Auction-System folder in your favorite shell
2. Go to the server sub folder
3. Write: "go run server.go 0"
3. Repeat previous 3 steps on two other terminals, with the numbers 1 or 2 replacing the 0 after server.go. (all terminals should have a different number between 0, 1 and 2)
4. Go to the client sub folder.
5. Write: "go run client.go (id)" with id being replaced with an int of your choosing
6. Repeat the previous two steps at least twice, with a different id
7. You can now enjoy the Auction
8. Call your mother and tell her you love her.
9. Find your personal purpose with life.

Nice to know: 
There are two commands, "info" and "bid (amount)" where amount is the amount you wish to bid.
Auctions last 1 minute and start over when someone bids while there is no Auction active.
Winners can be seen in log file.

