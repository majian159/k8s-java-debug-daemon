#!/usr/bin/env bash
# 指定最忙的前N个线程并打印堆栈
max_thread_count=50
# 指定cpu占比统计的采样间隔，单位为毫秒
interval=2000
# arthas的下载地址
arthas_boot_download_url="https://arthas.gitee.io/arthas-boot.jar"

# jps命令不存在代表非java应用或没有jdk相关工具
if [ ! "$(type jps)" ]; then
  exit 0
fi

cd ~/ || exit

# arthas-boot.jar 不存在则下载
if [ ! -f "arthas-boot.jar" ]; then
  curl -sO "$arthas_boot_download_url"
fi

unset JAVA_TOOL_OPTIONS
# 处理pid 1非java应用的情况
pid=$(jps | grep -v "Jps" | awk '{print $1}')

# 获取繁忙的前20个线程
java -jar arthas-boot.jar "${pid}" -c "thread -n ${max_thread_count} -i ${interval}"
