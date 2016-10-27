// Package twilio simplifies interaction with the Twilio API.
//
// The twilio-go library should be your first choice for interacting with the
// Twilio API; it offers forward compatibility, very fine-grained control of
// API access, best-in-class control over how long to wait for requests to
// complete, and great debuggability when things go wrong. Get started by
// creating a Client:
//
//     client := twilio.NewClient("AC123", "123", nil)
//
// All of the Twilio resources are available as properties on the Client. Let's
// walk through some of the example use cases.
//
// Creating a Resource
//
// Resources that can create new methods take a url.Values as an argument,
// and pass all arguments to the Twilio API. This method ensures forward
// compatibility; any new arguments that get invented can be added in
// client-side code.
//
//     data := url.Values{"To": []string{"+1foo"}, "From": []string{"+1foo"},
//         "Body": []string{"+1foo"}}
//     msg, err := client.Messages.Create(context.TODO(), data)
//
// Getting an Instance Resource
//
// Call Get() with a particular sid.
//
//     number, err := client.IncomingNumbers.Get(context.TODO(), "PN123")
//     fmt.Println(number.PhoneNumber)
//
// Getting a List Resource
//
// There are two flavors of interaction. First, if all you want is a single
// Page of resources, optionally with filters:
//
//      page, err := client.Recordings.GetPage(context.TODO(), url.Values{})
//
// To control the page size, set "PageSize": "N" in the url.Values{} field.
// Twilio defaults to returning 50 results per page if this is not set.
//
// Alternatively you can get a PageIterator and call Next() to repeatedly
// retrieve pages.
//
//     iterator := client.Calls.GetPageIterator(url.Values{"From": []string{"+1foo"}})
//     for {
//         page, err := iterator.Next(context.TODO())
//         // NoMoreResults means you've reached the end.
//         if err == twilio.NoMoreResults {
//             break
//         }
//         fmt.Println("start", page.Start)
//     }
package twilio
