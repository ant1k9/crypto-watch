#!/usr/bin/env fish

for rate in (psql $DB_DSN -qXAt -c 'select name from coins order by rank')
    echo -e Rate: $rate "\n"
    psql $DB_DSN -qXAt \
        -c "select value from rates where coin_uuid = (select uuid from coins where name = '$rate') order by ts" \
        | asciigraph -w 200 -h 16
    echo; echo; sleep 2
end
