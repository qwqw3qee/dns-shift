#!/bin/bash

# 获取当前脚本目录
DIR="$(cd "$(dirname "$0")" && pwd)"
DNS_SHIFT="$DIR/dns-shift"

# 定义全局变量
declare -i interval=1 hour=0 minute=0


# 检查 dns-shift 程序是否存在
check_dns_shift() {
  if [ ! -f "$DNS_SHIFT" ]; then
    echo "错误：未找到 dns-shift 程序，请确保该程序存在于此脚本同目录下。"
    exit 1
  fi
}

# 配置定时任务
add_cron_job() {
  # 将天数转换为适合 cron 表达式的格式
  local cron_job="$minute $hour */$interval * * $DNS_SHIFT > $DIR/dns-shift.log 2>&1"
  (crontab -l 2>/dev/null | grep -v "$DNS_SHIFT"; echo "$cron_job") | crontab -
  printf "已将 dns-shift 程序添加到 crontab，每 %d 天的 %02d:%02d 执行一次。\n" "$interval" "$hour" "$minute"

}

# 删除定时任务
remove_cron_job() {
  (crontab -l 2>/dev/null | grep -v "$DNS_SHIFT") | crontab -
  echo "已删除 dns-shift 的定时任务。"
}

# 获取间隔天数和运行时刻
get_interval_days_and_time() {
  if [ -z "$1" ]; then
    read -p "请输入每多少天检查一次（默认1天）： " interval
    interval=${interval:-1}
    # 获取小时和分钟
    read -p "请输入运行时刻小时（0-23，默认0）： " hour_input
    hour=${hour_input:-0}
    read -p "请输入运行时刻分钟（0-59，默认0）： " minute_input
    minute=${minute_input:-0}
  else
    interval=$1
    hour=$2
    minute=$3
  fi

  # 校验输入的天数是否为有效数字，若不合法或为空则默认为1天
  if [ "$interval" -lt 1 ]; then
    interval=1
  fi

  # 校验小时是否为有效数字
  if [ "$hour" -lt 0 ] || [ "$hour" -gt 23 ]; then
    hour=0
  fi

  # 校验分钟是否为有效数字
  if [ "$minute" -lt 0 ] || [ "$minute" -gt 59 ]; then
    minute=0
  fi
}

# 主函数
main() {
  check_dns_shift

  if [ "$1" == "rm" ]; then
    remove_cron_job
    exit 0
  fi

  get_interval_days_and_time "$1"

  chmod +x "$DNS_SHIFT"
  add_cron_job
}

main "$@"
