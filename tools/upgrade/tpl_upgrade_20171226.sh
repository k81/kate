#!/bin/bash
regex='([^:]+):([^:]+)'
for line in $(grep -Hn 'resetting ... ' script/*.sh|cut -f1,2 -d':')
do
    [[ $line =~ $regex ]]
    F=${BASH_REMATCH[1]}
    L=${BASH_REMATCH[2]}

    echo "updating $F ..."
    awk -v n=$(expr $L - 1) '
    NR == n {
        print "echo \"checking ... \""
        print "for k in ${!confs[@]}"
        print "do"
        print "    value=\"${confs[$k]}\""
        print "    jsonerr=$(echo \"$value\" | python -mjson.tool 2>&1)"
        print "    [ $? -ne 0 ] && echo \"[$k] syntax error, \\\"$jsonerr\\\"\" && exit 1"
        print "done"
        print ""
    }
    { print }' $F >$F.tmp && mv -f $F.tmp $F && chmod a+x $F
done
