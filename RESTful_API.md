Read the original article at https://hackernoon.com/restful-api-design-step-by-step-guide-2f2c9f9fcdbf

As software developers, most of us use or build REST APIs in a day to day life. APIs are the default means of communication between the systems. Amazon is the best example how of APIs can be efficiently used for communication. In this article, I am going to talk about how to design your RESTful APIs better to avoid common mistakes.

Jeff Bezos’s (Key to Success) Mandate
Some of you might have been already aware of Jeff Bezos’s mandate to the developers in Amazon. If you never got a chance to hear about it, following points are the crux of it —

All teams will henceforth expose their data and functionality through service interfaces.
Teams must communicate with each other through these interfaces.
There will be no other form of interprocess communication allowed: no direct linking, no direct reads of another team’s data store, no shared-memory model, no back-doors whatsoever. The only communication allowed is via service interface calls over the network.
It doesn’t matter what technology they use. HTTP, Corba, Pubsub, custom protocols — doesn’t matter. Bezos doesn’t care.
All service interfaces, without exception, must be designed from the ground up to be externalizable. That is to say, the team must plan and design to be able to expose the interface to developers in the outside world. No exceptions.
Anyone who doesn’t do this will be fired.
Eventually this turned to be the key to Amazon’s success. Amazon could build scalable systems and later could also offer those as services as Amazon Web Services.

Principles of Designing RESTful APIsReadt
Now let’s understand the principles we should follow while designing the RESTful APIs —

Keep it simple

Souce — Internet
We need to make sure that the base URL of the API is simple. For example, if we want to design APIs for products, it should be designed like —

/products
/products/12345
The first API is to get all products and the second one is to get specific product.

Use nouns and NOT the verbs
A lot of developers make this mistake. They generally forget that we have HTTP methods with us to describe the APIs better and end up using verbs in the API URLs. For instance, API to get all products should be

/products
and NOT as shown below

/getAllProducts
Some common URL patterns, I have seen so far

Use of right HTTP methods
RESTful APIs have various methods to indicate the type of operation we are going to perform with this API —

GET — To get a resource or collection of resources.
POST — To create a resource or collection of resources.
PUT/PATCH — To update the existing resource or collection of resources.
DELETE — To delete the existing resource or the collection of resources.
We need to make sure we use the right HTTP method for given operation.

Use Plurals
This topic is bit debatable. Some of people like to keep the resource URL with plural names while others like to keep it singular. For instance —

/products
/product
I like to keep it plural since it avoid confusion whether we are talking about getting single resource or collection. It also avoids adding additional things like attaching all to the base URL e.g. /product/all

Some people might not like this but my only suggestion is to keep is uniform across the project.

Use parameters
Sometime we need to have an API which should be telling more story than just by id. Here we should make use of query parameters to design the API.

/products?name=’ABC’ should be preffered over /getProductsByName
/products?type=’xyz’ should be preferred over /getProductsByType
This way you avoid long URLs with simplicity in design.

Use proper HTTP codes
We have plenty of HTTP codes. Most of us only end up using two — 200 and 500! This is certainly not a good practice. Following are some commonly used HTTP codes.

200 OK — This is most commonly used HTTP code to show that the operation performed is successful.
201 CREATED — This can be used when you use POST method to create a new resource.
202 ACCEPTED — This can be used to acknowledge the request sent to the server.
400 BAD REQUEST — This can be used when client side input validation fails.
401 UNAUTHORIZED / 403 FORBIDDEN— This can be used if the user or the system is not authorised to perform certain operation.
404 NOT FOUND— This can be used if you are looking for certain resource and it is not available in the system.
500 INTERNAL SERVER ERROR — This should never be thrown explicitly but might occur if the system fails.
502 BAD GATEWAY — This can be used if server received an invalid response from the upstream server.
Versioning
Versioning of APIs is very important. Many different companies use versions in different ways, some use versions as dates while some use versions as query parameters. I generally like to keep it prefixed to the resource. For instance —

/v1/products
/v2/products
I would also like to avoid using /v1.2/products as it implies the API would be frequently changing. Also dots (.) might not be easily visible in the URLs. So keep it simple.

It is always good practice to keep backward compatibility so that if you change the API version, consumers get enough time to move to the next version.

Use Pagination
Use of pagination is a must when you expose an API which might return huge data and if proper load balancing is not done, the a consumer might end up bringing down the service.

We need to always keep in mind that the API design should be full proof and fool proof.
Use of limit and offset is recommended here. For example, /products?limit=25&offset=50 It is also advised to keep a default limit and default offset.

Supported Formats
If is also important to choose how your API responds. Most of the modern day applications should return JSON responses unless you have an legacy app which still needs to get XML response.

Use Proper Error Messages
It is always a good practice to keep set of error messages application sends and respond that with proper id. For example, if you use Facebook graph APIs, in case of errors, it returns message like this —

{
  "error": {
    "message": "(#803) Some of the aliases you requested do not exist: products",
    "type": "OAuthException",
    "code": 803,
    "fbtrace_id": "FOXX2AhLh80"
  }
}
I have also seen some examples in which people return URL with error message which tells you more about the error message and how to handle it as well.

Use of Open API specifications
In order to keep all teams in your company abide to certain principles, use of OpenAPI Specification can be useful. Open API allows you to design your APIs first and share that with the consumers in easier manner.

Conclusion
It is quite evident that if you want to communicate better, APIs are the way to go. But if they are designed badly then it might increase confusion. So put best efforts to design well and rest is just the implementation.
