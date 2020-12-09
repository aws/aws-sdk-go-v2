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

# Convert tags in XML files in current folder to XML/Zonbook tags
# It takes two args, $1 is the folder where to start looking for XML files
# and $2 is the default programming language for programlisting tags
# If we don't have an arg, we just leave the language tag out of programlisting tags

pushd $1
make_list xml

# <title>           -> <info><title>
# </title>          -> </title></info>
# <literal>         -> <code>
# </literal>        -> </code>
# <link xlink:href  -> href/uling url
# </link>           -> </ulink>

for i in `cat $masterList`
do
    tmpfile=$$.tmp

    echo "<!DOCTYPE section PUBLIC \"-//OASIS//DTD DocBook XML V4.5//EN\" \"file://zonbook/docbookx.dtd\"" >> $tmpfile
    echo "[" >> $tmpfile
    echo "    <!ENTITY % xinclude SYSTEM \"file://AWSShared/common/xinclude.mod\">" >> $tmpfile
    echo "    %xinclude;" >> $tmpfile
    echo "    <!ENTITY % phrases-shared SYSTEM \"file://AWSShared/common/phrases-shared.ent\">" >> $tmpfile
    echo "    %phrases-shared;" >> $tmpfile

    if [ "$2" != "" ]
    then
	echo "    <!ENTITY % phrases-${2} SYSTEM \"../../shared/${2}.ent\">" >> $tmpfile
	echo "    %phrases-${2};" >> $tmpfile
    fi
    
    echo "]>" >> $tmpfile

    cat $i >> $tmpfile
    rm $i

    cat $tmpfile |
	sed "s/xml:id/id/g" | \
	    sed "s/<title>/\<info\>\<title\>/g" | \
	    sed "s/<\/title>/\<\/title\>\<\/info>/g" | \
	    sed "s/<literal>/<code>/g" | \
	    sed "s/<\/literal>/<\/code>/g" | \
	    sed "s/link xlink:href/ulink url/g" | \
            sed "s/<\/link>/<\/ulink>/g"	    > $i

    if [ "$2" != "" ]
    then
	rm $tmpfile
	cat $i | sed "s/<programlisting>/<programlisting language=\"$2\">/g" > $tmpfile
	mv $tmpfile $i
    else
	rm $tmpfile
    fi
done

rm $masterList > /dev/nul 2>&1
