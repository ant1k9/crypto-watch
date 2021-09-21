#!/usr/bin/env fish

function show_graph
    psql $DB_DSN -qXAt \
        -c "select value from rates where coin_uuid = (select uuid from coins where name = '$argv[1]') order by ts" \
        | asciigraph -w 200 -h 16
    echo -e "\n"
end

if test (count $argv) -eq 1
    show_graph "$argv[1]"
    exit
end

for rate in (psql $DB_DSN -qXAt -c 'select name from coins order by rank')
    echo -e "Rate: $rate \n"
    show_graph "$rate"
    sleep 2
end
