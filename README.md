# tcpproxy
tcpproxy is simple tool for communicating with inner network tcp server.     

```  
                                +----------+      +----------------+    
            +------------+      | gateway  |      |  inner network |    
  YOU <-->  | SSH tunnel | <--> | tcpproxy | <--> |  tcp server    |    
            +------------+      +----------+      +----------------+     
```

SSH tunnel is an optional component. You may make it with:
`ssh -p 36000 -fNL 3306:localhost:3306 toby@mnet.inner-network.com`

### What is tcpproxy use for
mnet.inner-network.com is the gateway of inner network(192.168.0.1) where deplyed a MySQL server(192.168.1.2:3306). You can only access mnet.inner-network.com in your mac but you want to access the MySQL service. Just execute 

	./tcpproxy 3306 192.168.1.2 3306
	 
at gateway then you can access the MySQL service with 

	mysql -h mnet.inner-network.com -u root -p

In some case, mnet.inner-network.com is a jump server and you have to ssh login it to access the inner network resource. Just make a ssh tunnel between your mac and the jump server, then you can access the MySQL service with

	mysql -h localhost -u root -p
