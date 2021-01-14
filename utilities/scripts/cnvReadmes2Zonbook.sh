#!/usr/bin/env bash

make_list () {
    masterList=./list.txt
    # Delete existing list
    rm $masterList > /dev/nul 2>&1
    find . -type f -name *.$1 -print > $masterList
    tmp=$$.txt
    cat $masterList | sed s/' '/'\\ '/g > $tmp
    mv $tmp $masterList
}

# Searches $1 for *.md files
# and converts them to XML/DocBook files.
# e.g. myFile.md -> myFile.xml

# If we don't have an arg, bail with error message
if [ "$1" == "" ]
then
    echo "You must tell me where to find the .MD files"
    exit
fi

# Navigate to $1
pushd $1 > /dev/nul 2>&1

# Create list of markdown files
make_list md

# For each markdown file
for i in `cat $masterList`
do
#    echo "Converting $i"
    # Convert it to XML
    basename=`basename $i .md`
    dirname=`dirname $i`
    pandoc -f markdown -t docbook $1/$i > $1/$dirname/$basename.xml
done

rm $masterList
popd > /dev/nul 2>&1
echo Done
