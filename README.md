cfg配置文件解析
===============

考试信息
--------

* 出题人: 王永振
* 考试时间: 20180510-20180516
* 审阅时间: 下周四五两天看哪天有时间
* 代码提交地址: https://git.cloud.top/microservices/go-challenge/20180509/002/姓名


必选功能：
--------

1. 可以从指定配置目录读取cfg文件（默认从tos.cfg读取，可指定cfg文件）
1. 文件可以读取cfg文件的命令段配置（每个文件可能存在多个）
1. 每个文件可能存在多个可以跳过空行，注释（#开头的字符串）
1. 开头的字符串如果存在引用分级cfg文件的，需要支持级联解析（最多四层）
1. 最多四层能够拆分文件字符串，获取节点信息
1. 获取节点信息能够判断出给定的字符串是否命中语法树中的节点（显示出命中的节点信息）
1. 显示出命中的节点信息能够打印出所有的命令行（手动输入的命令行）


增分功能：
--------

1. 增分功能可以根据查询条件，返回所有符合条件的命令行
1. 返回所有符合条件的命令行实现两组cfg文件比较，显示差异
1. 显示差异提供web服务，可通过api方式查看结果


评分标准：
--------

1. 评分标准满足基本功能
1. 评分标准满足基本功能在满足1.基础上，性能高者加分
1. 性能高者加分有实现增分项者酌情加分


技能要点：
--------

1. 技能要点golang文件读取
1. 文件读取Golang字符串拆分，拼接
1. 拼接Golang结构体使用
1. 结构体使用有限状态机（配置文件解析）
1. 配置文件解析Golang树结构


附录
----

**附1：示例**

示例文件1：tos.cfg

```
cmd_begin
_obj_type_keyword = system
_file = system.cfg
_obj_meaning_c = 系统设置
_obj_meaning_e = system configure
#_conf = system
cmd_end
# error code (-1000 ~ -2000)
error_begin
-1000= ivalid input parameter
-1000 = 无效的输入参数
error_end
```


示例文件2：system.cfg

```
cmd_begin
_obj_type_keyword = version
_cmd_keyword =
_cmd_function = /tos/so/libtos_system_config.so:tos_system_version
_conf = local-service
_popedom = read
_obj_meaning_e = system version
_act_meaning_e = show system version
_obj_meaning_c = 系统版本号
_act_meaning_c = 查看系统版本号
_xml_layout_begin
_xml_key = tos_version
_xml_layout_end
cmd_end

# error code (-101000 ~ -102999)
error_begin
-101101 = ivalid input parameter
-101101 = 无效的输入参数
error_end
```

**附2：cfg文件语法**

cfg文件格式语法（每一个命令描述段中不一定存在全部的字段属性）

```
cmd_begin                             #段的开始
# 段中不支持嵌套，如果出现段的开始符号，未出现段的结束符号，编译过程将失败。
_obj_type_keyword =    # 对象类型关键字。
_cmd_keyword =         # 动作,与命令类型编码相对应，是该命令的关键字。
                              对象可以有几种操作：增、移、删、改、清、查看。
_cmd_function =        # 执行函数名，包括动态库的文件名称。动态库的输出函数使用统一的接口。
_file =                # 分级文件名，全路径。执行函数名和分级文件名只能出现一个。
                              对于分级命令，该段描述分级命令的语法规则文件。其他段将失效。
_ha=                   #HA同步,分enable、disable
_license=              #命令许可证
_conf=                 #功能配置，与用户权限有关
_popedom =             #命令的执行权限，分super、write、read、gui、cli、noha、root、inter。
                              在执行命令之前，由认证权限管理模块对比该执行权限与用户的许可权限对比
                              所的结果决定执行与否。
_para = [ 关键字keyword= ] [ 必选可选required= ] [ 单选多选multi-value= ] [ 参数值类型type=]
                              [许可证license=][系统配置对象obj=][无序unsequent=]
                              [ 参数英文描述meaning_e= ][参数中文描述meaning_c=]
                       #关键字：属性名。
                       #必选可选：说明该属性是否必须填写。必选属性将排列在可选属性前面。
                       #参数描述：对该参数的文字描述，便于理解该属性的含义。
                       #参数值类型：可以定义几种固定的类型，比如IP地址、掩码、字符串、整形值等等，
                              选择型参数提供多个可选项，命中可选项即可通过检查。
                       #单选多选：说明该属性可以填写多个值。
_obj_meaning_e =       # 命令含义的英文解释语句。用户键入TAB时将输出该段。
_obj_meaning_c =       # 命令含义的中文解释语句。用户键入TAB时将输出该段。
_act_meaning_e =       # 动作含义的英文解释语句。用户键入TAB时将输出该段。
_act_meaning_c =       # 动作含义的中文解释语句。用户键入TAB时将输出该段。
cmd_end                        #段的结束

```

**错误信息段语法：**

```
error_begin=-N:-M
-N=example
-N=示例
error_end
```
