# List of sub-features and fixes needed to make the code RFC compliant

* Implement Error Estimate for Sender and Receiver based on [RFC4656](https://www.rfc-editor.org/rfc/rfc4656#section-4.1.2) and [RFC5357](https://www.rfc-editor.org/rfc/rfc5357.html#page-23)
* Implement Control messaging similar to RFC-5357
* Move the test traffic off of port 862 and use a separate connection, preferably UDP
* Set TTL to 225 and verify
* Benchmark binary serialization cost to the latency calculation in accordance to defined fields in the RFC
* Consider adding Authenticated Mode
* Interop testing with a router/switch with TWAMP enabled