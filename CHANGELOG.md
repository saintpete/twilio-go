# Changes

## 0.55

Handle new HTTPS-friendly media URLs.

Support deleting/releasing phone numbers via IncomingNumbers.Release(ctx, sid).

Initial support for the Pricing API (https://pricing.twilio.com).

Add an AUTHORS.txt.

## 0.54

Add Recordings.GetTranscriptions() to get Transcriptions for a recording. The
Transcriptions resource doesn't support filtering by Recording Sid.

## 0.53

Add Alert.StatusCode() function for retrieving a HTTP status code (if one
exists) from an alert.

## 0.52

Copy url.Values for GetXInRange() functions before modifying them.

## 0.51

Implement GetNextConferencesInRange

## 0.50

Implement GetConferencesInRange. Fix paging error in
GetCallsInRange/GetMessagesInRange.

## 0.47

Implement GetNextXInRange - if you have a next page URI and want to get an
Iterator (instead of starting with a url.Values).

## 0.45

Fix go 1.6 (messages_example_test) relied on the stdlib Context package by
accident.

## 0.44

Support filtering Calls/Messages down to the nanosecond in a TZ-aware way, with
Calls.GetCallsInRange / Messages/GetMessagesInRange.

## 0.42

Add more Description fields based on errors I've received in the past. There
are probably more to be found, but this is a good start.

## 0.41

Use the same JWT library instead of using two different ones.

Add Description() for Alert bodies.

## 0.40

Fix next page URL's for Twilio Monitor

## 0.39

The data in Update() requests was silently being ignored. They are not ignored
any more.

Support the Accounts resource.

Add RequestOnBehalfOf function to make requests on behalf of a subaccount.

Fixes short tests that were broken in 0.38

## 0.37

Support Outgoing Caller ID's

## 0.36

Support Keys

## 0.35

Added Ended(), EndedUnsuccessfully() helpers to a Call, and FriendlyPrice() to
a Transcription.
