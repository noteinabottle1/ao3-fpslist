# ao3-fpslist
A collection of Golang scripts that scrape Ao3 for data

## How to use this
Go to the Releases page and download the latest binary. To run, go to a terminal and
run:
`./bplist-generator`

To see the options, run:
```
./bplist-generator -help
Usage of ./bplist-generator:
      --authors strings   Authors with BP statements (e.g. a,b,c)
      --id int            AO3 Fandom ID
      --maxWords int      Max words in a fic (default 5000)
      --minWords int      Min words in a fic (default 50)
```

## Program Inputs

There are two required inputs to this script:
1) The fandom ID

To find this, you'll need to go to a user, go to their Dashboard page, and click on the
fandom tag you're interested in filtering by. For example:
`https://archiveofourown.org/users/orphan_account/works?fandom_id=911149`

2) The list of authors with a BP statement.

To find this, go to fpslist.org, search for the fandom you're interested in, and see the
list of authors with BP statements in that fandom. For example:
https://www.fpslist.org/fandom/13862/

Then, copy the column of authors and convert them into a comma separated list.
(I use https://convert.town/column-to-comma-separated-list).

Run this program, plug both values in when prompted, and you should have an output csv file at
`fandom_id_bplist.csv`!
