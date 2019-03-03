#!/bin/bash
echo "It will generate the NOMATCH file, which is the list of URLs that are in ALL, but not in DOWNLOADED."
grep -vf DOWNLOADED_*.txt ALL_*.txt > NOMATCH.txt
echo "After this, compare NOMATCH with DOWNLOADED."
