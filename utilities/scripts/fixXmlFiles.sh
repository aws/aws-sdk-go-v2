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

# Searches $1 for *.lmx files
# and renames them to *.xml.
# e.g. myFile.lmx -> myFile.xml

# If we don't have an arg, bail with error message
if [ "$1" == "" ]
then
    echo "You must tell me where to find the .MD files"
    exit
fi

# Navigate to $1
pushd $1 > /dev/nul 2>&1

# If list of lmx or md files exists, delete them
rm lmx_masterList.txt > /dev/nul 2>&1
rm md_masterList.txt > /dev/nul 2>&1

# Create list of lmx files
make_list lmx

# For each lmx file
echo "Converting LMX files to XML"
for i in `cat $masterList`
do
#    echo "Converting $i"
    # Convert it to XML
    basename=`basename $i .lmx`
    dirname=`dirname $i`
    mv $1/$i $1/$dirname/$basename.xml
done

# Create list of md files
make_list md

# For each md file
echo "Deleting MD files"
for i in `cat $masterList`
do
#    echo "Deleting $i"
    rm $i
done

# Delete list
rm $masterList > /dev/nul 2>&1

popd > /dev/nul 2>&1
echo Done
