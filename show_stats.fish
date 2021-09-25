#!/usr/bin/env fish

## helpers
function show_graph
    argparse 'c/coin=' -- $argv
    set -l rank (psql $DB_DSN -qXAt -c "SELECT rank FROM coins WHERE name = '$_flag_coin'")
    echo -e "[$rank] $_flag_coin\n"
    psql $DB_DSN -qXAt \
        -c "SELECT value FROM rates WHERE coin_uuid = (SELECT uuid FROM coins WHERE name = '$_flag_coin') ORDER BY ts" \
        | asciigraph -w 200 -h 16
    echo -e "\n"
end

function show_graph_by_query
    argparse 'q/query=' -- $argv
    for rate in (psql $DB_DSN -qXAt -c "$_flag_query")
        show_graph -c "$rate"
        sleep 2
    end
end

function trending_query
    argparse 'o/order-by=' -- $argv
    echo 'WITH coin_stats AS (
        SELECT
            coin_uuid,
            value,
            ROW_NUMBER() OVER (PARTITION BY coin_uuid ORDER BY ts DESC) rn
        FROM rates
    )
    SELECT coins.name FROM coin_stats cs1
    INNER JOIN coin_stats cs2 USING(coin_uuid)
    INNER JOIN coins ON cs1.coin_uuid = coins.uuid
    WHERE
        cs1.rn = 1 AND cs2.rn = 45 AND cs1.value > 0 AND cs2.value > 0'
    echo "ORDER BY $_flag_order_by LIMIT 10"
end

## body
argparse 'h/help' 'show=' 't/trending' 'd/descending' -- $argv

if set -q _flag_help
    echo -n 'Usage:
    ./show_stats.fish                  # help charts for all coins
    ./show_stats.fish --help           # help message
    ./show_stats.fish --show Ethereum  # make Ethereum chart
    ./show_stats.fish --trending       # make charts for top 10 trending coins
    ./show_stats.fish --descending     # make charts for top 10 descending coins
'
else if set -q _flag_show
    show_graph -c "$_flag_show"
else if set -q _flag_trending
    set -l query (trending_query -o "cs2.value / cs1.value")
    show_graph_by_query -q "$query"
else if set -q _flag_descending
    set -l query (trending_query -o "cs1.value / cs2.value")
    show_graph_by_query -q "$query"
else
    show_graph_by_query -q 'SELECT name FROM coins ORDER BY rank'
end
