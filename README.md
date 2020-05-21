# gopher-mail
Simple serverless personal mail system. Use AWS free tier or cheap services.

- SES
- S3
- SNS
- Lambda
- API Gateway

Moving parts:
- postmaster - handles sorting of all incoming mail
- mailman - reads mail and serves it to the web ui
- mailtruck - sends mail via SES
- web - vuejs spa for the user to interact with the systems above

## TODO
* [x] Move to terraform for infrastructure and lambda deployment
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


## Tech
### Build
The build script will build each of the services into the bin folder.

```bash
./build.sh
```

### Deploy

The deploy script will take the built binaries in the bin folder and create an archive for deployment to a lambda function.

Terraform is used to deploy and configure the services required to glue everything together. Create a terraform.tfvars from the example in the terrafrom folder and run the command below.
```bash
./deploy.sh

cd terraform/
terraform apply
```
_**Note:** The `terraform apply` command will take around 10 minutes._

If you would like to use an S3 bucket for your terraform state, there is an example override in the terraform folder. You will need to create a private S3 bucket in your account and enter the bucket name into the override file.

## Restrictions
Using SES to S3 email delivery caps email size at 30MB. At a later time this can be updated to use lambdas exclusively for a payload size only limited by the HTTP protocol and the lambda memory.

## Price
_**Note:** The following figures are preliminary and do not include data transfer rates._

The following table describes the price per month.
| Price | Service |
| --- | --- |
| $0.40 | AWS Secrets |

Total per year: $4.80

The following table describes worst case price per email, beyond the free tier.
| Price | Service | Free Tier Limit |
| --- | --- | --- |
| ~ $0.00000001 | AWS Lambda | 1,000,000 |
| ~ $0.0001 | SES | 1,000 |
| ~ $0.000005 | S3 | 20,000 |

Web UI costs for viewing one email.
| Price | Service | Free Tier Limit |
| --- | --- | --- |
| ~ $0.00000004 | AWS Lambda | 1,000,000 |
| ~ $0.000005 | S3 | 20,000 |