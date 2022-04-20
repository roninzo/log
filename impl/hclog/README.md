[![back](https://img.shields.io/badge/navigation-back-green "Back")](../README.md) 

Extracted from [Go Microservice with Clean Architecture: Application Logging](https://medium.com/@jfeng45/go-microservice-with-clean-architecture-application-logging-b43dc5839bce)

# Comparison of log library:

Different log libraries provide different features, some are important for debugging.

Logging information that is important (the following data are required):

1. File name and line number
2. Method name and caller name
3. Message logging level
4. Timestamp
5. Error stack trace
6. Automatically logging each function call with parameters and results

Logging from [Awesome Go](https://github.com/libanglang/awesome-go#logging)