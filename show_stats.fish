#!/usr/bin/env fish

## helpers
function list_coins
    psql $DB_DSN -qXAt -c "SELECT rank, name FROM coins ORDER BY rank LIMIT 40"
end

function show_graph
    argparse 'c/coin=' 'h/height=' 'w/width=' -- $argv
    if test "$_flag_height" = ""
        set _flag_height "16"
    end
    if test "$_flag_width" = ""
        set _flag_width "200"
    end

    set -l rank (psql $DB_DSN -qXAt -c "SELECT rank FROM coins WHERE name = '$_flag_coin'")
    echo "         [$rank] $_flag_coin"
    set -l query "
        SELECT value FROM (
            SELECT value, ts
            FROM rates
            WHERE coin_uuid = (SELECT uuid FROM coins WHERE name = '$_flag_coin')
            ORDER BY ts DESC
            LIMIT $_flag_width
        ) tmp
        ORDER BY ts"

    psql $DB_DSN -qXAt -c "$query" | asciigraph -w "$_flag_width" -h "$_flag_height"
    echo
end

function show_graph_by_query
    argparse 'q/query=' 'no-wait=' 'h/height=' 'w/width=' -- $argv
    for rate in (psql $DB_DSN -qXAt -c "$_flag_query")
        show_graph -c "$rate" -h "$_flag_height" -w "$_flag_width"
        if not test "$_flag_no_wait" = '--no-wait'
            sleep 2
        end
    end
end

function trending_query
    argparse 'o/order-by=' 'd/days=' 't/top=' -- $argv
    if not set -q "_flag_days"
        set -l _flag_days "30"
    end
    if not set -q "_flag_days"
        set -l _flag_days "30"
    end
    if test "$_flag_top" = ""
        set _flag_top "10"
    end

    echo 'WITH coin_stats AS (
        SELECT
            coin_uuid,
            value,
            ROW_NUMBER() OVER (PARTITION BY coin_uuid ORDER BY ts DESC) rn
        FROM rates
    )
    SELECT coins.name FROM coin_stats cs1
    INNER JOIN coin_stats cs2 USING(coin_uuid)
    INNER JOIN coins ON cs1.coin_uuid = coins.uuid'
    echo "WHERE cs1.rn = 1 AND cs2.rn = $_flag_days AND cs1.value > 0 AND cs2.value > 0 AND rank <= 40"
    echo "ORDER BY $_flag_order_by LIMIT $_flag_top"
end

## body
argparse \
    'h/help' \
    'list' \
    'show=' \
    'trending' \
    'descending' \
    'd/days=' \
    'no-wait' \
    'width=' \
    'top=' \
    'height=' -- $argv

if set -q _flag_help
    echo -n 'Usage:
    ./show_stats.fish                          # charts for all coins
    ./show_stats.fish --width 100 --height 10  # customize width and height of charts
    ./show_stats.fish --help                   # help message
    ./show_stats.fish --show Ethereum          # make Ethereum chart
    ./show_stats.fish --trending --days 7      # make charts for top 10 trending coins
    ./show_stats.fish --descending             # make charts for top 10 descending coins
'

else if set -q _flag_list
    list_coins

else if set -q _flag_show
    show_graph \
        -c "$_flag_show" \
        -w "$_flag_width" \
        -h "$_flag_height"

else if set -q _flag_trending
    set -l query (trending_query -o "cs2.value / cs1.value" -d "$_flag_days" -t "$_flag_top")
    show_graph_by_query \
        -q "$query" \
        --no-wait "$_flag_no_wait" \
        -w "$_flag_width" \
        -h "$_flag_height"

else if set -q _flag_descending
    set -l query (trending_query -o "cs1.value / cs2.value" -d "$_flag_days" -t "$_flag_top")
    show_graph_by_query \
        -q "$query" \
        --no-wait "$_flag_no_wait" \
        -w "$_flag_width" \
        -h "$_flag_height"

else
    show_graph_by_query \
        -q 'SELECT name FROM coins ORDER BY rank' \
        --no-wait "$_flag_no_wait" \
        -w "$_flag_width" \
        -h "$_flag_height"
end
