#!/usr/bin/env bash

ETCD_HOST='QA_HOST:2379'

function query_etcd() {
    [ $# -lt 2 ] && echo "no cmd specified" && exit 1
    case $1 in 
        'set')
            [ $# -ne 3 ] && echo "Usage: query_etcd set key value" && exit 1
            curl -s http://$ETCD_HOST/v2/keys/$2 -X PUT --data-urlencode "value=$3" >/dev/null
            [ $? -eq 0 ] && echo $3
            ;;
        'rm')
            [ $# -ne 2 ] && echo "Usage: query_etcd rm key" && exit 1
            curl -s http://$ETCD_HOST/v2/keys/$2?recursive=true -X DELETE >/dev/null
            [ $? -eq 0 ] && echo "key $2 deleted"
            ;;
    esac
}

CONF_ROOT="__APP_NAME__/v1/conf"

declare -A confs

#日志等级 
confs['log']='{
    "level": "DEBUG"
}'
#是否开启profiling
confs['profiling']='{
    "enabled": true,
    "port": 18000
}'
confs['redis/addrs']='{
    "addrs":[
    "QA_REDIS_HOST:6379"
	]
}'
confs['redis/pools']='{
    "connect_timeout":6,
    "read_timeout":3,
    "write_timeout":3,
    "idle_timeout":30,
    "max_idle":20, 
    "max_active":50, 
    "wait":true
}'
confs['db/basic']='{
    "debug_sql": true,
    "driver_name": "cdbpool",
    "data_source": "tcp(QA_HOST:9123)/default?timeout=5s&readTimeout=4s&writeTimeout=15s&enableCircuitBreaker=true"
}'
confs['db/pools']='{
    "max_idle_conns": 20,
    "max_open_conns": 100,
    "conn_max_lifetime": 0
}'
#confs['kafka/brokers']='{
    #"brokers": [
        #"QA_HOST:9092"
    #]
#}'
#confs['kafka/consumers']='{
    #"app": "__APP_NAME__",
    #"consumers": [
        #{
            #"topic": "EXAMPLE_TOPIC",
            #"worker_num": 1,
            #"rebalance_interval": 10,
            #"worker_report_interval": 1,
            #"file_offset_scale_factor": 0,
            #"redis_offset_scale_factor": 0
        #}
    #]
#}'
#设置http参数
confs['http']='{
    "port": 8000,
    "request_timeout_secs": 5,
    "read_timeout_secs": 3,
    "write_timeout_secs": 5,
    "max_header_bytes": 1048576,
    "max_body_bytes": 1073741824
}'

echo "checking ... "
for k in ${!confs[@]}
do
    value="${confs[$k]}"
    jsonerr=$(echo "$value" | python -mjson.tool 2>&1)
    [ $? -ne 0 ] && echo "[$k] syntax error, \"$jsonerr\"" && exit 1
done

if [ $# -eq 0 ]; then
    echo "resetting ... "
    query_etcd rm ${CONF_ROOT} 2>/dev/null

    echo

    for k in ${!confs[@]}
    do
        echo "configuring $k"
        query_etcd set "${CONF_ROOT}/${k}" "${confs[$k]}"
    done

    echo "done"
else

    until [ $# -eq 0 ]
    do
        k=$1
        echo "updating $k"
        if [ "x${confs[$k]}" = x"" ]; then
            echo "config \"$k\" not found"
            exit 1
        fi
        query_etcd set "${CONF_ROOT}/${k}" "${confs[$k]}"
        shift
    done

    echo "done"
fi

