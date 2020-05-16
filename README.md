# gopher-mail
Simple serverless personal mail system. Use AWS free tier or cheap services.

- SES
- S3
- SNS
- Lambda

Moving parts:
- postmaster - handles sorting of all incoming mail
- mailman - reads mail and serves it to the web ui
- mailtruck - sends mail via SES
- web - vuejs spa for the user to interact with the systems above

---
## TODO
* [ ] Create simple front end to request mail for a mailbox
  * [x] Web skeleton
  * [ ] Add multiple user auth
  * [ ] Add sending of email
  * [ ] Add rendering of html payloads
* [ ] Create mailman service to load and serve mail
* [ ] Create mailtruck service to send mail
* [ ] Feature creep postmaster
  * [ ] Parse AWS-SES headers for virus checking and create sub-folders [spam, trash]
* [ ] Feature creep web
  * [ ] Implement Auth
  * [ ] Use markdown for email editor and MD to HTML for the html emails


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

