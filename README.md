# REST Data Archiver
This application is intended to provide a simple serverless function to dump the output of a REST API call into an S3 bucket on a regular basis. The result can be used for backup, analysis, or for access by another application. It is meant to be flexible in terms of the source and destination APIs. The primary source adapter implemented is a basic REST API, and the primary destination adapter is AWS S3. Since _destinations_ have their own unique APIs and integration methods, each _destination_ is developed individually to implement the _Destination_ interface. The runtime for this application is configured using a `config.json` file. An example is provided named 
`config.example.json`.

## Sources

### REST API
Data sources coming from simple API calls can use the `RestAPI` source. Here are some examples of how to configure it:

#### Basic Authentication
```json
{
  "Source": {
    "Type": "RestAPI",
    "AdapterConfig": {
      "Method": "GET",
      "BaseURL": "https://example.com",
      "AuthType": "basic",
      "Username": "username",
      "Password": "password",
      "UserAgent": "rest-data-archiver"
    }
  },
  "Sets": [
    {
      "Name": "Users",
      "Source": {
        "Path": "/users"
      },
      "Destination": {
      }
    }
  ]
}
```

#### Bearer Token Authentication
```json
{
  "Source": {
    "Type": "RestAPI",
    "AdapterConfig": {
      "Method": "GET",
      "BaseURL": "https://example.com",
      "AuthType": "bearer",
      "Password": "token",
      "UserAgent": "rest-data-archiver"
    }
  }
}
```
`Sets` is configured the same as for basic authentication.

#### Salesforce OAuth Authentication
```json
{
  "Source": {
    "Type": "RestAPI",
    "AdapterConfig": {
      "Method": "GET",
      "BaseURL": "https://login.salesforce.com/services/oauth2/token",
      "AuthType": "SalesforceOauth",
      "Username": "admin@example.com",
      "Password": "the-password",
      "ClientID": "put-your-client-id-here",
      "ClientSecret": "put-your-client-secret-here",
      "UserAgent": "rest-data-archiver"
    }
  },
  "Sets": [
    {
      "Name": "Contacts",
      "Source": {
        "Path": "/services/data/v20.0/query/?q=SELECT%20Email,Name%20FROM%20Contacts"
      },
      "Destination": {
      }
    }
  ]
}
```

## Destinations

### Amazon AWS S3
Using the S3 adapter, data from each Set will be saved as an object in the
configured S3 bucket.

Here is an example for how to configure it:

```json
{
  "Destination": {
    "Type": "S3",
    "AdapterConfig": {
      "BucketName": "my-archive-bucket",
      "AwsConfig": {
        "Region": "us-east-1",
        "AccessKeyId": "your-access-key-id",
        "SecretAccessKey": "secret-access-key"
      }
    }
  },
  "Sets": [
    {
      "Name": "Users",
      "Source": {
        "Path": "/users"
      },
      "Destination": {
        "ObjectNamePrefix": "users"
      }
    }
  ]
}
```

### Email Alerts

Event Log events with a level of LOG_ALERT or LOG_EMERG will result in an email 
alert sent via AWS SES. Note that the LOG_EMERG level is 0, which is the Go
zero-value. Any new event log created without a Level assigned will default to
LOG_EMERG and could result in email alerts being sent. 

The following is an example configuration:

```json
{
  "Alert": {
    "AWSRegion": "us-east-1",
    "CharSet": "UTF-8",
    "ReturnToAddr": "no-reply@example.org",
    "SubjectText": "rest-data-archiver alert",
    "RecipientEmails": [
      "admin@example.org"
    ],
    "AWSAccessKeyID": "ABCD1234",
    "AWSSecretAccessKey": "abcd1234!@#$"
  }
}
```

Alternatively, AWS credentials can be supplied by the Serverless framework by
adding the following configuration to `serverless.yml`:

```
provider:
  iamRoleStatements:
    - Effect: 'Allow'
      Action:
        - 'ses:SendEmail'
      Resource: "*"
```

Both authentication mechanisms are provided in the `lamdda-example` directory, 
but only one is needed.

### Exporting logs from CloudWatch

The log messages in CloudWatch can be viewed on the AWS Management Console. If
an exported text or json file is needed, the AWS CLI tool can be used as
follows:

```shell script
aws configure
aws logs get-log-events \
   --log-group-name "/aws/lambda/lambda-name" \
   --log-stream-name '2019/11/14/[$LATEST]0123456789abcdef0123456789abcdef' \
   --output text \
   --query 'events[*].message'
```

Replace `/aws/lambda/lambda-name` with the actual log group name and 
`2019/11/14/[$LATEST]0123456789abcdef0123456789abcdef` with the actual log
stream. Note the single quotes around the log stream name to prevent the shell
from interpreting the `$` character. `--output text` can be changed to 
`--output json` if desired. Timestamps are available if needed, but omitted
in this example by the `--query` string.

