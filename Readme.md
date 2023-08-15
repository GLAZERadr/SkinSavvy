# Doculex AI 

Doculex AI: your ultimate solution for efficient legal document management

## Introduction 

Doculex AI is a web application that allows users to upload, store, and search for legal documents. The application is built using the Go Fiber framework and is hosted on AWS EC2 Virtual Private Server.

## Installation Guides

* Install Go 

[Go Installation Link Here!](https://go.dev/doc/install)

* Clone This Repository

    ```
        https://github.com/DataaAlchemist/Doculex-API.git
    ```

* Add Value To ".env.example" File >> Change ".env.example" To ".env"
    You Must Add:
    1. MongoDB URL
    2. Database Name
    3. Host
    4. Port 
    5. JWT Secret Key

* Run The Application With Your Terminal

    ```
        go run main.go
    ``` 

## API URL

[Link For API Integration](http://18.143.77.67:8000/)

## Project Support Features

* Users can signup and login to their accounts

* Public (non-authenticated) users can access all public encpoint on the platform

* Authenticated users can access all protected endpoint.

## API Documentation
* User: [Click here for API documentation!](https://www.postman.com/api-documentation-tool/)

## Technologies Used
* Golang - Backend Language
* Fiber - Backend Framework
* MongoDB - Database
* AWS EC2 - Hosting Platform
* AWS S3 - Storage
* Alibaba Cloud - LLM llama2
* Github - Version Control
* Postman - API Documentation

## License

The MIT License (MIT)

Copyright (c) 2023 Adrian Glazer

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.