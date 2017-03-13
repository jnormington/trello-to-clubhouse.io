# trello-to-clubhouse.io

## What is it ?
This is a script written with Go which migrates the cards from a specific list on a board in Trello and into
[Clubhouse.io](https://clubhouse.io).

This program takes an interactive approach by asking questions and querying the api for the things we need
and takes action from there instead of you trying hard to find whats needs and it taking too long.

It currently migrates the following from a Trello Card what we felt was valuable;

- Name
- Description
- Labels
- Due Date
- Creator
- Created At
- Comments
- Members ([Requested by @meganchinburg](https://github.com/jnormington/trello-to-clubhouse.io/issues/3))
- Checklists (and also whether the checklist item is completed)
- ShortURL (optional comment added with Trello link)
- Attachments (optional uploads attachments to dropbox)

If you are also making the move from Trello to Clubhouse.io and want some extra attributes copied from a Trello Card
feel free to create an issue or submit a pull request.

## Why the use of dropbox?

I know what you are thinking, why, why not just use the current Trello url for the attachment ? Well the attachments are on the Trello S3 bucket meaning that once the card is deleted so are the images. Even if you don't delete the card maybe
Trello might see some traffic from Clubhouse and ban it so for safety we just download and upload to dropbox and share the resource.

All the attachments are uploaded under trello

I understand its not perfect and maybe using the direct link is your preferred route if this is the case please fork and modify.

## Setup

Before we can run the program we need to get all the keys and tokens from the services

#### Trello (Key and Token)
[You can get the key here](https://trello.com/app-key)

To get your token use the following url, but remember to replace the key at the key

```
https://trello.com/1/authorize?expiration=1day&name=MigrationFromTrelloToClubhouse&response_type=token&key=REPLACEWITHYOURKEY
```

#### Clubhouse (Token)

[You can create a token here](https://app.clubhouse.io/tester1234/settings/account/api-tokens)

#### Dropbox (Token, only if you plan to migrate attachments)
[You can create a token here](https://www.dropbox.com/developers/apps/create)


## Usage

Now you have all your tokens you need to set them up as environment variables.

Download the binary for your platform from the list below

- [Windows (x86)](https://github.com/jnormington/trello-to-clubhouse.io/releases/download/v0.2.1/trello_to_clubhouse_windows_x86.exe)
- [Windows (x64)](https://github.com/jnormington/trello-to-clubhouse.io/releases/download/v0.2.1/trello_to_clubhouse_windows_x64.exe)
- [Linux (x64)](https://github.com/jnormington/trello-to-clubhouse.io/releases/download/v0.2.1/trello_to_clubhouse_linux_x64)
- [OSX (x64)](https://github.com/jnormington/trello-to-clubhouse.io/releases/download/v0.2.1/trello_to_clubhouse_osx_x64)


Then move to the next section for settings the environment variables for your platform

#### Windows

Open command prompt and type the following one line at a time and press return key on each line

```
set CLUBHOUSE_TOKEN=YOURTOKEN
set TRELLO_KEY=YOURKEY
set TRELLO_TOKEN=YOURTOKEN
set DROPBOX_TOKEN=YOURTOKEN
```

Now with the same command window open drag the downloaded binary into it and press return key. You
should now be asked several questions and before long cards will be in clubhouse.io.

You can see an example below if similar output and the questions you will be asked along the journey.

#### OSX/Linux

Open terminal app for your platform and type the following one line at a time and press return key on each line

```
export CLUBHOUSE_TOKEN=YOURTOKEN
export TRELLO_KEY=YOURKEY
export TRELLO_TOKEN=YOURTOKEN
export DROPBOX_TOKEN=YOURTOKEN
```

Now you might need to make the binary executable to do so type the following `chmod +x ` and then drag
the downloaded binary press return key. This should now be executable by any user.

With the same terminal window open drag the downloaded binary into it and press return key. You
should now be asked several questions and before long the cards will be in clubhouse.io.

You can see an example below if similar output and the questions you will be asked along the journey.



## Example program questions/output (specific to my accounts)

```
$ ./trello-to-clubhouse.io

Would you like to migrate all attachments from trello cards?
This will entail downloading the attachments and uploading to dropbox
A dropbox account will be required for the token
[0] Yes
[1] No

Please select a board by its number
[0] Bugs
[1] Scorpian

Please select the list to import by number
[0] New
[1] High
[2] Medium
[3] Low

Please wait while we retrieve your cards... This might take a few minutes.
Please select a project by it number to import the cards into
[0] Project Two
[1] Project Zero
[2] Bugs

Please select a workflow state to import the cards into
[0] Ready for Development
[1] In Development
[2] Completed

Please select the user account to import the cards as
[0] Test Account
[1] Jon

Please select the story type all cards should be imported as
[0] feature
[1] chore
[2] bug

Would you like a comment added with the original trello ticket link?
[0] Yes
[1] No

To correctly map ticket owners to Clubhouse we need a user mapping CSV.
If this is the first time running this program you need to generate one.
We generate a csv of a best guess user mapping which you can edit to be correct
If you already have one that is correct please select option 1
Please select your option based on the above information:
[0] Yes
[1] No

*********************
 CSV generated: /home/jon/Documents/userMappingTtoC.csv
*********************
Is your CSV user mapping correct ?
CSV file: /home/jon/Documents/userMappingTtoC.csv
Are you ready to continue ?
[1] Yes

****** WARNING ******
Please review carefully before you continue

Export cards from Trello
        Board: Bugs
        List: New


Import cards into clubhouse
        Project: Bugs
        Workflow State: Ready for Development
        Story Type: bug
        Add Comment with Trello Link: true

Is the above correct select the number representing your answer ?
[0] Yes
[1] No

Importing trello cards into Clubhouse...
Trello Card Link                         Import Status     Error/Story ID

https://trello.com/c/TUQxmMFB            Success           Story ID: 649
https://trello.com/c/GBWkKhW0            Success           Story ID: 652
https://trello.com/c/YolLMisX            Success           Story ID: 655
https://trello.com/c/VuaXkO2X            Success           Story ID: 660
https://trello.com/c/9eIaDF7n            Success           Story ID: 666
*** Looks like we finished go and have fun & joy with Clubhouse ***
```
