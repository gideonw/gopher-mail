# gopher-mail
Simple serverless personal mail system. Use AWS free tier or cheap services.

- SES
- S3
- SNS
- Lambda

---
## TODO
* [ ] Create simple front end to request mail for a mailbox
  * [ ] Add multiple user auth
  * [ ] Add sending of email
  * [ ] Add rendering of html payloads
* [ ] Create mailman service to load and serve mail
* [ ] Create mailtruck service to send mail
* [ ] Feature creep postman
  * [ ] Parse AWS-SES headers for virus checking and create sub-folders [spam, trash]


---
## Tech
### Build
```bash
./build.sh
```
or
```bash
make
```

### Deploy
Serverless is used to deploy and configure the services.
```bash
serverless deploy -v
```

