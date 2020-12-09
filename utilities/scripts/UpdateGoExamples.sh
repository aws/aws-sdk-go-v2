#!/usr/bin/env bash

usage() { echo "Usage: $0 -p OUTPUT-PATH -u GITHUB-USERNAME [-l LANGUAGE] (but only go is supported) -k md | xml -t index | metadata | none" 1>&2; exit 1; }

while getopts ":d:p:u:l:k:t:" o; do
    case "${o}" in
        p)
            p=${OPTARG}
            ;;
        u)
            u=${OPTARG}
            ;;
	l)
	    l=${OPTARG}
	    ;;
	k)
	    k=${OPTARG}
	    ((k == md || k == xml)) || usage
	    ;;
	t)
	    t=${OPTARG}
	    ((t == index || t == metadata || t == none)) || usage
	    ;;
        *)
            usage
            ;;
    esac
done
shift $((OPTIND-1))

if [ -z "${p}" ]; then
    echo "You must specify a path (-p PATH)"
    usage
fi

if [ -z "${u}" ]; then
    echo "You must specify a GitHub user name (-u NAME)"
    usage
fi

if [ -z "${k}" ]; then
    echo "You must specify a file extension to keep (-k md | xml)"
    usage
fi

if [ -z "${t}" ]; then
    echo "You must specify a type of outout name (-t index | metadata | none)"
    usage
fi

echo "path = ${p}"
echo "user = ${u}"
echo "keep = ${k}"
echo "type = ${t}"

if [ "$l" != "" ]
then
    echo "lang = ${l}"
    echo
    echo "Creating empty ${p}/shared/${l}.ent"
    mkdir $p/shared
    touch $p/shared/$l.ent
fi

echo
echo "Calling: traverse-repo-file -u $u -l $l -o $p -t $t"
echo "To Create MD files from GitHub repo"
pushd ../cmd/traverse-repo-files > /dev/null 2>&1
go run . -u $u -l $l -o $p -t $t
popd > /dev/null 2>&1

if [ "$k" == "md" ]
then
    echo "Finished with MD files"
    exit
fi

echo
echo "Calling: ./cnvReadmes2Zonbook.sh $p"
echo "To convert MD files in ${p} to XML/DocBook"
./cnvReadmes2Zonbook.sh $p

echo
echo "Calling ./FixZonbookTags.sh $p $l"
echo "To convert DocBook tags -> Zonbook tags"
./FixZonbookTags.sh $p $l

echo
echo "Calling patch-zonebook-files -p $p -l $l"
echo "To create LMX files in ${p} so the title has an ID (section ID + .title)"
pushd ../cmd/patch-zonebook-file > /dev/null 2>&1
go run . -p $p -l $l
popd > /dev/null 2>&1

echo
echo "Calling ./fixXmlFiles.sh $p"
echo "To rename LMX files to XML AND deleting any MD files"
./fixXmlFiles.sh $p

echo
echo Done
