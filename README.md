##What is Spickspan
Spickspan is an open source library that provides two main features. First, it is a client side library for interacting with Essentier Nomock's RESTful APIs. Second, it provides extensible and platform-agnostic service discovery. As we mentioned earlier, Essentier Nomock is about running integration tests against live, running services. Services often have dependencies on other services. Hence we need a service discovery mechanism to help wire up those dependencies. Spickspan is Essentier's answer to such a need.

##Why Spickspan
The challenge of service discovery is that we often want to wire up the dependencies differently depending on where our services run. For example, let's say we have a web application that depends on a database. In the production environment, we would like to connect the web application to the production database. Whereas during testing, we would like to connect the web application to the testing database. The production environment might be hosted up on a cloud such as Amazon AWS, Heroku, Google Cloud Platform, etc. The testing environment might be based on Essentier Nomock or it might be hosted on-premise. 

No matter in which environment or on which platform our services run, we want the same service discovery mechanism to work so that the code we test can be deployed to production without any change. In other words, the service discovery mechanism needs to be platform-agnostic. It also needs to be extensible so that new environments and new platforms can be supported. And that is exactly what Spickspan provides--an extensible, platform-agnostic service discovery library.

##Use Spickspan

