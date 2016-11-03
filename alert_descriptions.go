package twilio

var alertDestination = []byte(`
{
    "account_sid": "AC58f1e8f2b1c6b88ca90a012a4be0c279",
    "alert_text": "sourceComponent=14100&ErrorCode=14101&LogLevel=ERROR&Msg=The+destination+number+for+a+TwiML+message+can+not+be+the+same+as+the+originating+number+of+an+incoming+message.&EmailNotification=false",
    "api_version": "2008-08-01",
    "date_created": "2016-10-26T01:11:13Z",
    "date_generated": "2016-10-26T01:11:12Z",
    "date_updated": "2016-10-26T01:11:18Z",
    "error_code": "14101",
    "log_level": "error",
    "more_info": "https://www.twilio.com/docs/errors/14101",
    "request_headers": null,
    "request_method": "POST",
    "request_url": "https://kev.inburke.com/zombo/zombo.php",
    "request_variables": "ToCountry=US&ToState=CA&SmsMessageSid=SM1521fa559fb923c1ed64f56cd1bed8ef&NumMedia=0&ToCity=BRENTWOOD&FromZip=94514&SmsSid=SM1521fa559fb923c1ed64f56cd1bed8ef&FromState=CA&SmsStatus=received&FromCity=BRENTWOOD&Body=twilio-go+testing%21&FromCountry=US&To=%2B19253920364&ToZip=94514&NumSegments=1&MessageSid=SM1521fa559fb923c1ed64f56cd1bed8ef&AccountSid=AC58f1e8f2b1c6b88ca90a012a4be0c279&From=%2B19253920364&ApiVersion=2010-04-01",
    "resource_sid": "SM1521fa559fb923c1ed64f56cd1bed8ef",
    "response_body": "<Response><Sms>You can do anything at ZomboCom. Anything at all.</Sms></Response>",
    "response_headers": "Transfer-Encoding=chunked&X-Cache=MISS+from+ip-172-18-20-243.ec2.internal&Server=cloudflare-nginx&X-Cache-Lookup=MISS+from+ip-172-18-20-243.ec2.internal%3A3128&Content-Type=text%2Fxml%3Bcharset%3Dutf-8&Date=Wed%2C+26+Oct+2016+01%3A11%3A13+GMT&CF-RAY=2f7a0875aaca21b0-EWR&X-Powered-By=PHP%2F5.6.16",
    "service_sid": null,
    "sid": "NO7e3853acc314b52d8b6babd04ede0a39",
    "url": "https://monitor.twilio.com/v1/Alerts/NO7e3853acc314b52d8b6babd04ede0a39"
}
`)

var alert11200 = []byte(`
{
    "account_sid": "AC58f1e8f2b1c6b88ca90a012a4be0c279",
    "alert_text": "Msg&sourceComponent=12000&ErrorCode=11200&httpResponse=405&url=https%3A%2F%2Fkev.inburke.com%2Fzombo%2Fzombocom.mp3&LogLevel=ERROR",
    "api_version": "2010-04-01",
    "date_created": "2016-10-27T02:34:21Z",
    "date_generated": "2016-10-27T02:34:21Z",
    "date_updated": "2016-10-27T02:34:23Z",
    "error_code": "11200",
    "log_level": "error",
    "more_info": "https://www.twilio.com/docs/errors/11200",
    "request_headers": null,
    "request_method": "POST",
    "request_url": "https://kev.inburke.com/zombo/zombocom.mp3",
    "request_variables": "Called=%2B19252717005&ToState=CA&CallerCountry=US&Direction=outbound-api&CallerState=CA&ToZip=94596&CallSid=CA6d27370cbbfb605521fe8800bb73f2d2&To=%2B19252717005&CallerZip=94514&ToCountry=US&ApiVersion=2010-04-01&CalledZip=94596&CalledCity=PLEASANTON&CallStatus=in-progress&From=%2B19253920364&AccountSid=AC58f1e8f2b1c6b88ca90a012a4be0c279&CalledCountry=US&CallerCity=BRENTWOOD&Caller=%2B19253920364&FromCountry=US&ToCity=PLEASANTON&FromCity=BRENTWOOD&CalledState=CA&FromZip=94514&FromState=CA",
    "resource_sid": "CA6d27370cbbfb605521fe8800bb73f2d2",
    "response_body": "<html>\r\n<head><title>405 Not Allowed</title></head>\r\n<body bgcolor=\"white\">\r\n<center><h1>405 Not Allowed</h1></center>\r\n<hr><center>nginx</center>\r\n</body>\r\n</html>",
    "response_headers": "Transfer-Encoding=chunked&Server=cloudflare-nginx&CF-RAY=2f82bf9cb8102204-EWR&Set-Cookie=__cfduid%3Dd46f1cfd57d664c3038ae66f1c1de9e751477535661%3B+expires%3DFri%2C+27-Oct-17+02%3A34%3A21+GMT%3B+path%3D%2F%3B+domain%3D.inburke.com%3B+HttpOnly&Date=Thu%2C+27+Oct+2016+02%3A34%3A21+GMT&Content-Type=text%2Fhtml",
    "service_sid": null,
    "sid": "NO00ed1fb4aa449be2434d54ec8e492349",
    "url": "https://monitor.twilio.com/v1/Alerts/NO00ed1fb4aa449be2434d54ec8e492349"
}
`)

var alertUnknown = []byte(`
{
    "account_sid": "AC58f1e8f2b1c6b88ca90a012a4be0c279",
    "alert_text": "",
    "api_version": "2010-04-01",
    "date_created": "2016-10-27T02:34:21Z",
    "date_generated": "2016-10-27T02:34:21Z",
    "date_updated": "2016-10-27T02:34:23Z",
    "error_code": "235342434",
    "log_level": "error",
    "more_info": "https://www.twilio.com/docs/errors/93455",
    "request_headers": null,
    "request_method": "POST",
    "request_url": "https://kev.inburke.com/zombo/zombocom.mp3",
    "request_variables": "",
    "resource_sid": "CA6d27370cbbfb605521fe8800bb73f2d2",
    "response_body": "",
    "response_headers": "",
    "service_sid": null,
    "sid": "NO00ed1fb4aa449be2434d54ec8e492349",
    "url": "https://monitor.twilio.com/v1/Alerts/NO00ed1fb4aa449be2434d54ec8e492349"
}
`)
